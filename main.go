package main

import (
	"AggregatorProject/internal/database"
	"database/sql"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"

	_ "github.com/lib/pq"
)

func main() {
	envErr := godotenv.Load()
	if envErr != nil {
		log.Fatal("Could not load .env file")
	}

	port := os.Getenv("PORT")
	dbURL := os.Getenv("DB_CONNECTION_STRING")

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal("Error connecting to the database", err)
		return
	}

	dbQueries := database.New(db)

	apicfg := apiConfig{
		DB: dbQueries,
	}

	serveMux := http.NewServeMux()

	serveMux.HandleFunc("/v1/readiness", readinessHandler)
	serveMux.HandleFunc("/v1/err", errorHandler)

	serveMux.HandleFunc("POST /v1/users", apicfg.createUser)

	corsMux := CORSMiddleware(serveMux)

	server := &http.Server{
		Addr:    ":" + port,
		Handler: corsMux,
	}

	log.Printf("Server starting on port: %s...", port)
	serverListenErr := server.ListenAndServe()
	if serverListenErr != nil {
		log.Fatal("Failed to start server")
	}
}
