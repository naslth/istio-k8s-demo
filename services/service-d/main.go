package main

import (
	"encoding/json"
	"net/http"
	"os"
	"time"

	runtime "github.com/banzaicloud/logrus-runtime-formatter"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	amqp "github.com/rabbitmq/amqp091-go"
	log "github.com/sirupsen/logrus"
)

var (
	logLevel     = getEnv("LOG_LEVEL", "debug")
	port         = getEnv("PORT", ":8080")
	serviceName  = getEnv("SERVICE_NAME", "Service D")
	message      = getEnv("GREETING", "Shalom (שָׁלוֹם), from Service D!")
	queueName    = getEnv("QUEUE_NAME", "service-d.greeting")
	rabbitMQConn = getEnv("RABBITMQ_CONN", "amqp://guest:guest@rabbitmq:5672")
)

type Greeting struct {
	ID          string    `json:"id,omitempty"`
	ServiceName string    `json:"service,omitempty"`
	Message     string    `json:"message,omitempty"`
	CreatedAt   time.Time `json:"created,omitempty"`
	Hostname    string    `json:"hostname,omitempty"`
}

var greetings []Greeting

func GreetingHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)

	// Create a new greeting message
	greetings = nil
	tmpGreeting := Greeting{
		ID:          uuid.New().String(),
		ServiceName: serviceName,
		Message:     message,
		CreatedAt:   time.Now().Local(),
		Hostname:    getHostname(),
	}

	// Store it in the greetings slice
	greetings = append(greetings, tmpGreeting)

	// Respond with the greeting message
	err := json.NewEncoder(w).Encode(tmpGreeting)
	if err != nil {
		log.Error("Error encoding greeting to JSON:", err)
	}

	// Prepare headers for Jaeger tracing
	incomingHeaders := []string{
		"x-b3-flags",
		"x-b3-parentspanid",
		"x-b3-sampled",
		"x-b3-spanid",
		"x-b3-traceid",
		"x-ot-span-context",
		"x-request-id",
	}

	rabbitHeaders := amqp.Table{}
	for _, header := range incomingHeaders {
		if r.Header.Get(header) != "" {
			rabbitHeaders[header] = r.Header.Get(header)
		}
	}

	log.Debug("Sending message with headers:", rabbitHeaders)

	// Marshal the greeting for RabbitMQ
	body, err := json.Marshal(tmpGreeting)
	if err != nil {
		log.Error("Error marshalling greeting to JSON:", err)
		return
	}

	// Send the message to RabbitMQ
	sendMessage(rabbitHeaders, body, rabbitMQConn)
}

func HealthCheckHandler(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	_, err := w.Write([]byte("{\"alive\": true}"))
	if err != nil {
		log.Error("Error sending health check response:", err)
	}
}

func sendMessage(headers amqp.Table, body []byte, rabbitMQConn string) {
	conn, err := amqp.Dial(rabbitMQConn)
	if err != nil {
		log.Error("Failed to connect to RabbitMQ:", err)
		return
	}
	defer func() {
		if err := conn.Close(); err != nil {
			log.Error("Error closing RabbitMQ connection:", err)
		}
	}()

	ch, err := conn.Channel()
	if err != nil {
		log.Error("Failed to open a channel:", err)
		return
	}
	defer func() {
		if err := ch.Close(); err != nil {
			log.Error("Error closing RabbitMQ channel:", err)
		}
	}()

	q, err := ch.QueueDeclare(
		queueName,
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Error("Failed to declare queue:", err)
		return
	}

	err = ch.Publish(
		"",
		q.Name,
		false,
		false,
		amqp.Publishing{
			Headers:     headers,
			ContentType: "application/json",
			Body:        body,
		})
	if err != nil {
		log.Error("Failed to send message to RabbitMQ:", err)
	}
}

func getHostname() string {
	hostname, err := os.Hostname()
	if err != nil {
		log.Error("Error getting hostname:", err)
		return ""
	}
	return hostname
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

func run() error {
	router := mux.NewRouter()
	api := router.PathPrefix("/api").Subrouter()
	api.HandleFunc("/greeting", GreetingHandler).Methods("GET")
	api.HandleFunc("/health", HealthCheckHandler).Methods("GET")
	api.Handle("/metrics", promhttp.Handler())
	return http.ListenAndServe(port, router)
}

func init() {
	formatter := runtime.Formatter{ChildFormatter: &log.JSONFormatter{}}
	formatter.Line = true
	log.SetFormatter(&formatter)
	log.SetOutput(os.Stdout)
	level, err := log.ParseLevel(logLevel)
	if err != nil {
		log.Error("Error parsing log level:", err)
	}
	log.SetLevel(level)
}

func main() {
	if err := run(); err != nil {
		log.Fatal("Error running the server:", err)
		os.Exit(1)
	}
}
