package main

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"os"
	"time"
	"sync"

	runtime "github.com/banzaicloud/logrus-runtime-formatter"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	amqp "github.com/rabbitmq/amqp091-go"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	logLevel     = getEnv("LOG_LEVEL", "debug")
	port         = getEnv("PORT", ":8080")
	serviceName  = getEnv("SERVICE_NAME", "Service F")
	message      = getEnv("GREETING", "Hola, from Service F!")
	queueName    = getEnv("QUEUE_NAME", "service-d.greeting")
	mongoConn    = getEnv("MONGO_CONN", "mongodb+srv://naslth:9015@k8s-istio-mongo.fpfctzc.mongodb.net/?retryWrites=true&w=majority&appName=k8s-istio-mongo")
	rabbitMQConn = getEnv("RABBITMQ_CONN", "amqps://pplaujxh:pN4wi0uhlSFAU2enS7EZTKWdIXNhZh2H@albatross.rmq.cloudamqp.com/pplaujxh")
)

type Greeting struct {
	ID          string    `json:"id,omitempty"`
	ServiceName string    `json:"service,omitempty"`
	Message     string    `json:"message,omitempty"`
	CreatedAt   time.Time `json:"created,omitempty"`
	Hostname    string    `json:"hostname,omitempty"`
}

var greetings []Greeting
var wg sync.WaitGroup

// *** HANDLERS ***

func GreetingHandler(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	greetings = nil
	tmpGreeting := Greeting{
		ID:          uuid.New().String(),
		ServiceName: serviceName,
		Message:     message,
		CreatedAt:   time.Now().Local(),
		Hostname:    getHostname(),
	}
	greetings = append(greetings, tmpGreeting)
	// Call MongoDB to store the greeting asynchronously
	wg.Add(1)
	go callMongoDB(tmpGreeting, mongoConn)

	// Respond with the greeting
	err := json.NewEncoder(w).Encode(greetings)
	if err != nil {
		log.Error(err)
	}
}

func HealthCheckHandler(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	_, err := w.Write([]byte("{\"alive\": true}"))
	if err != nil {
		log.Error(err)
	}
}

// *** UTILITY FUNCTIONS ***

func getHostname() string {
	hostname, err := os.Hostname()
	if err != nil {
		log.Error(err)
	}
	return hostname
}

func callMongoDB(greeting Greeting, mongoConn string) {
	defer wg.Done()

	log.Info("Storing greeting in MongoDB:", greeting)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoConn))
	if err != nil {
		log.Error("Failed to connect to MongoDB:", err)
		return
	}
	defer client.Disconnect(ctx)

	// Ping MongoDB to ensure connection is healthy
	if err := client.Ping(ctx, nil); err != nil {
		log.Error("MongoDB connection failed:", err)
		return
	}

	collection := client.Database("service-f").Collection("messages")
	_, err = collection.InsertOne(ctx, greeting)
	if err != nil {
		log.Error("Failed to insert greeting into MongoDB:", err)
	}
}

func getMessages(rabbitMQConn string) {
	for {
		conn, err := amqp.Dial(rabbitMQConn)
		if err != nil {
			log.Error("Failed to connect to RabbitMQ:", err)
			time.Sleep(5 * time.Second) // Retry after a delay
			continue
		}
		defer conn.Close()

		ch, err := conn.Channel()
		if err != nil {
			log.Error("Failed to create a RabbitMQ channel:", err)
			time.Sleep(5 * time.Second) // Retry after a delay
			continue
		}
		defer ch.Close()

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
			time.Sleep(5 * time.Second) // Retry after a delay
			continue
		}

		msgs, err := ch.Consume(
			q.Name,
			"service-f",
			true,
			false,
			false,
			false,
			nil,
		)
		if err != nil {
			log.Error("Failed to consume messages:", err)
			time.Sleep(5 * time.Second) // Retry after a delay
			continue
		}

		for delivery := range msgs {
			log.Debug("Received message:", delivery.Body)
			tmpGreeting := deserialize(delivery.Body)
			// Call MongoDB asynchronously
			wg.Add(1)
			go callMongoDB(tmpGreeting, mongoConn)
		}
	}
}

func deserialize(b []byte) (t Greeting) {
	log.Debug("Deserializing message:", b)
	var tmpGreeting Greeting
	buf := bytes.NewBuffer(b)
	decoder := json.NewDecoder(buf)
	err := decoder.Decode(&tmpGreeting)
	if err != nil {
		log.Error("Failed to deserialize message:", err)
	}
	return tmpGreeting
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

func run() error {
	// Start receiving messages from RabbitMQ in a separate goroutine
	go getMessages(rabbitMQConn)

	// Initialize HTTP server
	router := mux.NewRouter()
	api := router.PathPrefix("/api").Subrouter()
	api.HandleFunc("/greeting", GreetingHandler).Methods("GET")
	api.HandleFunc("/health", HealthCheckHandler).Methods("GET")
	api.Handle("/metrics", promhttp.Handler())

	// Start HTTP server
	return http.ListenAndServe(port, router)
}

func init() {
	// Initialize log settings
	formatter := runtime.Formatter{ChildFormatter: &log.JSONFormatter{}}
	formatter.Line = true
	log.SetFormatter(&formatter)
	log.SetOutput(os.Stdout)
	level, err := log.ParseLevel(logLevel)
	if err != nil {
		log.Error("Invalid log level:", err)
	}
	log.SetLevel(level)
}

func main() {
	// Start the application and ensure graceful error handling
	if err := run(); err != nil {
		log.Fatal("Failed to start the server:", err)
		os.Exit(1)
	}
}
