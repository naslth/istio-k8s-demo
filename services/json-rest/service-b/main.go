package main

import (
	"encoding/json"
	runtime "github.com/banzaicloud/logrus-runtime-formatter"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	log "github.com/sirupsen/logrus"
)

var (
	logLevel    = getEnv("LOG_LEVEL", "debug")
	port        = getEnv("PORT", ":8080")
	serviceName = getEnv("SERVICE_NAME", "Service B")
	message     = getEnv("GREETING", "Namasté (नमस्ते), from Service B!")
	URLServiceD = getEnv("SERVICE_D_URL", "http://service-d")
	URLServiceE = getEnv("SERVICE_E_URL", "http://service-e")
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

	log.Debug(r)

	greetings = nil

	callNextServiceWithTrace(URLServiceD+"/api/greeting", r)
	callNextServiceWithTrace(URLServiceE+"/api/greeting", r)

	tmpGreeting := Greeting{
		ID:          uuid.New().String(),
		ServiceName: serviceName,
		Message:     message,
		CreatedAt:   time.Now().Local(),
		Hostname:    getHostname(),
	}

	greetings = append(greetings, tmpGreeting)

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

func callNextServiceWithTrace(url string, r *http.Request) {
	log.Debug(url)

	var tmpGreetings []Greeting

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Error(err)
	}

	incomingHeaders := []string{
		"x-b3-flags",
		"x-b3-parentspanid",
		"x-b3-sampled",
		"x-b3-spanid",
		"x-b3-traceid",
		"x-ot-span-context",
		"x-request-id",
	}

	for _, header := range incomingHeaders {
		if r.Header.Get(header) != "" {
			req.Header.Add(header, r.Header.Get(header))
		}
	}

	log.Info(req)

	client := &http.Client{
		Timeout: time.Second * 10,
	}
	response, err := client.Do(req)

	if err != nil {
		log.Error(err)
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Error(err)
		}
	}(response.Body)

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Error(err)
	}

	err = json.Unmarshal(body, &tmpGreetings)
	if err != nil {
		log.Error(err)
	}

	for _, r := range tmpGreetings {
		greetings = append(greetings, r)
	}
}

func getHostname() string {
	hostname, err := os.Hostname()
	if err != nil {
		log.Error(err)
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
		log.Error(err)
	}
	log.SetLevel(level)
}

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
}
