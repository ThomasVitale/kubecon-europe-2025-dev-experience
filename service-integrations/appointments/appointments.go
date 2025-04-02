package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"

	dapr "github.com/dapr/go-sdk/client"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	_ "github.com/lib/pq"
)

var (
	VERSION            = getEnv("VERSION", "1.0.0")
	SOURCE             = getEnv("SOURCE", "https://github.com/")
	APP_PORT           = getEnv("APP_PORT", "8081")
	PostgresqlHost     = getEnv("POSTGRES_HOST", "localhost")
	PostgresqlPort     = getEnv("POSTGRES_PORT", "5432")
	PostgresqlUsername = getEnv("POSTGRES_USERNAME", "postgres")
	PostgresqlPassword = getEnv("POSTGRES_PASSWORD", "postgres")
	PUBSUB_NAME        = getEnv("PUBSUB_NAME", "pubsub")
	PUBSUB_TOPIC_NAME  = getEnv("PUBSUB_TOPIC_NAME", "topic")
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
	log.Printf("Starting Appointments Service in Port: %s", APP_PORT)
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

type Event struct {
	Id      string `json:"id"`
	Payload string `json:"payload"`
	Type    string `json:"type"`
}

// NewChiServer creates a new Chi server.
func NewChiServer() *chi.Mux {
	log.Printf("Starting Appointments Service in Port: %s", APP_PORT)

	// create new router
	r := chi.NewRouter()

	// add middlewares
	r.Use(middleware.Logger)

	apiClient, err := dapr.NewClient()
	if err != nil {
		panic(err)
	}

	// connect to database
	db := NewDB()

	// check if database is alive
	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Connected to PostgreSQL.")

	// create new server
	server := NewServer(db, apiClient)

	// add routes
	r.Get("/", server.Welcome)
	r.Get("/appointments/", server.GetAllAppointments)
	r.Post("/appointments/", server.CreateAppointment)
	r.Delete("/appointments/", server.DeleteAllAppointments)

	return r
}

// server is the API server struct
type server struct {
	DB        *sql.DB
	APIClient dapr.Client
}

// NewServer creates a newServer.
func NewServer(db *sql.DB, client dapr.Client) *server {
	return &server{
		DB:        db,
		APIClient: client,
	}
}

// GetAllAppointments returns all appointments.
func (s *server) GetAllAppointments(w http.ResponseWriter, r *http.Request) {
	var query = "SELECT id, patientId, category, appointmentDate FROM Appointments a"
	var rows *sql.Rows
	var err error

	rows, err = s.DB.Query(query)

	if err != nil {
		log.Printf("There was an error executing the query %v", err)
	}

	defer rows.Close()
	appointments := []Appointment{}
	for rows.Next() {

		var appointment Appointment
		err = rows.Scan(&appointment.Id, &appointment.PatientId, &appointment.Category, &appointment.AppointmentDate)
		if err != nil {
			log.Printf("There was an error scanning the sql rows: %v", err)
		}
		appointments = append(appointments, appointment)

	}

	log.Printf("Appointments retrieved from Database: %d", len(appointments))
	respondWithJSON(w, http.StatusOK, appointments)
}

// DeleteAllAppointments delete all appointments.
func (s *server) DeleteAllAppointments(w http.ResponseWriter, r *http.Request) {
	var deleteStmt = "DELETE FROM Appointments"

	var err error

	_, err = s.DB.Exec(deleteStmt)

	if err != nil {
		log.Printf("There was an error executing the query %v", err)
	}

	log.Printf("All Appointments deleted from Database.")
	respondWithJSON(w, http.StatusOK, "")
}

// CreateAppointment creates a new appointment.
func (s *server) CreateAppointment(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	var appointment Appointment
	err := json.NewDecoder(r.Body).Decode(&appointment)
	if err != nil {
		log.Printf("There was an error decoding the request body into the struct: %v", err)
		respondWithJSON(w, http.StatusInternalServerError, err)
		return
	}

	appointment.Id = uuid.New().String()

	insertStmt := `insert into Appointments(id, patientId, category, appointmentDate) values($1, $2, $3, $4)`

	_, err = s.DB.Exec(insertStmt, appointment.Id, appointment.PatientId, appointment.Category, appointment.AppointmentDate)

	if err != nil {
		log.Printf("An error occurred while executing query: %v", err)
		respondWithJSON(w, http.StatusInternalServerError, err)
		return
	}

	log.Printf("Appointment Stored in Database: %v", appointment)

	appointmentJson, err := json.Marshal(appointment)
	if err != nil {
		log.Printf("An error occured while marshalling the appointment to json: %v", err)
		respondWithJSON(w, http.StatusInternalServerError, err)
		return
	}

	event := Event{
		Id:      uuid.New().String(),
		Type:    "appointment-created",
		Payload: string(appointmentJson),
	}
	eventJson, err := json.Marshal(event)
	if err != nil {
		log.Printf("An error occured while marshalling the event to json: %v", err)
		respondWithJSON(w, http.StatusInternalServerError, err)
		return
	}

	if err := s.APIClient.PublishEvent(ctx, PUBSUB_NAME, PUBSUB_TOPIC_NAME, eventJson); err != nil {
		log.Printf("An error occured while publishing the event: %v", err)
		respondWithJSON(w, http.StatusInternalServerError, err)
		return
	}

	log.Printf("New Event Published with Appointment: %v", appointment)

	respondWithJSON(w, http.StatusOK, appointment)

}

// Welcome returns a welcome message from the Appointments Service
func (s *server) Welcome(w http.ResponseWriter, r *http.Request) {
	var welcome Welcome = Welcome{
		Message: "Welcome: you are now talking to the Appointments Service",
	}
	w.Header().Set(ContentType, ApplicationJson)
	json.NewEncoder(w).Encode(welcome)
}

func NewDB() *sql.DB {
	connStr := "postgresql://" + PostgresqlUsername + ":" + PostgresqlPassword + "@" + PostgresqlHost + ":" + PostgresqlPort + "/postgres?sslmode=disable"
	log.Printf("Connecting to Database: %s.", connStr)
	// Connect to database

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	return db
}

const (
	ApplicationJson = "application/json"
	ContentType     = "Content-Type"
)

type Appointment struct {
	Id              string    `json:"id"`
	PatientId       string    `json:"patientId"`
	Category        string    `json:"category"`
	AppointmentDate time.Time `json:"appointmentDate"`
}

type Welcome struct {
	Message string `json:"message"`
}

func (s Appointment) MarshalBinary() ([]byte, error) {
	return json.Marshal(s)
}
