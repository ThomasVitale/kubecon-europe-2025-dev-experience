package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
)

var (
	VERSION  = getEnv("VERSION", "1.0.0")
	SOURCE   = getEnv("SOURCE", "https://github.com/")
	APP_PORT = getEnv("APP_PORT", "8082")
)

// respondWithJSON is a helper function to write a JSON response.
func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

func main() {
	chiServer := NewChiServer()

	// Start the server; this is a blocking call
	err := http.ListenAndServe(":"+APP_PORT, chiServer)
	log.Printf("Starting Front Desk Service in Port: %s", APP_PORT)
	if err != http.ErrServerClosed {
		log.Panic(err)
	}
}

// getEnv returns the value of an environment variable, or a fallback value if not set.
func getEnv(key, fallback string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		value = fallback
	}
	return value
}

// NewChiServer creates a new Chi server.
func NewChiServer() *chi.Mux {
	log.Printf("Starting Front Desk Service in Port: %s", APP_PORT)

	// create new router
	r := chi.NewRouter()

	// add middlewares
	r.Use(middleware.Logger)

	// create new server
	server := NewServer()

	// add routes
	r.Get("/", server.Welcome)
	r.Get("/frontdesk/", server.GetAllFrontDeskEvents)
	r.Post("/frontdesk/events", server.ConsumeEvents)

	return r
}

// server is the API server struct
type server struct {
}

// NewServer creates a newServer.
func NewServer() *server {
	return &server{}
}

var events = []Event{}

// GetAllFrontDeskEvents events.
func (s *server) GetAllFrontDeskEvents(w http.ResponseWriter, r *http.Request) {
	log.Printf("Events received: %d", len(events))
	respondWithJSON(w, http.StatusOK, events)
}

// ConsumeEvents from other hospital services
func (s *server) ConsumeEvents(w http.ResponseWriter, r *http.Request) {

	var event Event
	err := json.NewDecoder(r.Body).Decode(&event)
	if err != nil {
		log.Printf("There was an error decoding the request body into the struct: %v", err)
	}

	events = append(events, event)

	log.Printf("Event received: %v", event)

	respondWithJSON(w, http.StatusOK, event)

}

// Welcome returns a welcome message from the Front Desk Service
func (s *server) Welcome(w http.ResponseWriter, r *http.Request) {
	var welcome Welcome = Welcome{
		Message: "Welcome: you are now talking to the Front Desk Service",
	}
	w.Header().Set(ContentType, ApplicationJson)
	json.NewEncoder(w).Encode(welcome)
}

const (
	ApplicationJson = "application/json"
	ContentType     = "Content-Type"
)

type Event struct {
	Id      string `json:"id"`
	Payload string `json:"payload"`
	Type    string `json:"type"`
}

type Welcome struct {
	Message string `json:"message"`
}

func (e Event) MarshalBinary() ([]byte, error) {
	return json.Marshal(e)
}
