package main

import (
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"get-paste-service/cache"
	"get-paste-service/config"
	"get-paste-service/db"
	"get-paste-service/producer"
	"get-paste-service/router"
)

func main() {
	// Load configuration
	config.LoadEnv()

	// Initialize components
	db.InitDB()
	cache.InitRedis()
	producer.InitRabbitMQ()

	// Set up graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	// Set up HTTP server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8082"
	}

	r := router.SetupRouter()

	// Start the server in a goroutine
	go func() {
		log.Printf("Get Paste Service started at :%s", port)
		if err := http.ListenAndServe(":"+port, r); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server error: %v", err)
		}
	}()

	// Wait for shutdown signal
	<-quit
	log.Println("Shutting down server...")

	// Close connections
	producer.CloseRabbitMQ()

	log.Println("Server shutdown complete")
}
