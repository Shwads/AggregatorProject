package main

import (
	"AggregatorProject/internal/database"
	"AggregatorProject/internal/scraper"
	"context"
	"database/sql"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

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

	ctx, cancel := context.WithCancel(context.Background())

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		sig := <-sigChan
		log.Printf("Received signal %s: Shutting down gracefully...", sig)
		cancel()
	}()

	go scraper.ScraperWorker(ctx, dbQueries)

	serveMux := http.NewServeMux()

	serveMux.HandleFunc("/v1/readiness", readinessHandler)
	serveMux.HandleFunc("/v1/err", errorHandler)

	serveMux.HandleFunc("POST /v1/users", apicfg.createUser)
	serveMux.HandleFunc("GET /v1/users", apicfg.middlewareAuth(authedGetUser))

	serveMux.HandleFunc("POST /v1/feeds", apicfg.middlewareAuth(apicfg.authedCreateFeed))
	serveMux.HandleFunc("GET /v1/feeds", apicfg.getFeeds)

	serveMux.HandleFunc("GET /v1/feed_follows", apicfg.middlewareAuth(apicfg.authedGetFeedFollows))
	serveMux.HandleFunc("POST /v1/feed_follows", apicfg.middlewareAuth(apicfg.authedNewFeedFollow))
	serveMux.HandleFunc("DELETE /v1/feed_follows/{feedFollowID}", apicfg.middlewareAuth(apicfg.authedDeleteFeedFollow))

	corsMux := CORSMiddleware(serveMux)

	server := &http.Server{
		Addr:    ":" + port,
		Handler: corsMux,
	}

	go func() {
		log.Printf("Server starting on port: %s...", port)
		if serverListenErr := server.ListenAndServe(); serverListenErr != http.ErrServerClosed {
			log.Print(serverListenErr)
		}
	}()

	<-ctx.Done()

	ctxShutdown, cancelShutdown := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelShutdown()

	if err := server.Shutdown(ctxShutdown); err != nil {
		log.Fatalf("Server forced to shutdown: %s", err)
	}

	log.Println("Worker has been shutdown.")
}
