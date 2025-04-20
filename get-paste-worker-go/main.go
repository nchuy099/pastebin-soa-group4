package main

import (
	"log"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	"get-paste-worker/config"
	"get-paste-worker/consumer"
	"get-paste-worker/db"
)

// Default worker configuration
const (
	DefaultNumWorkers = 10 // Default number of worker goroutines
)

func main() {
	log.Println("Starting get paste worker...")

	// Load environment variables and initialize database
	config.LoadEnv()
	db.InitDB()

	// Get number of workers from environment variable
	numWorkers := DefaultNumWorkers
	if numWorkersEnv := os.Getenv("NUM_WORKERS"); numWorkersEnv != "" {
		if n, err := strconv.Atoi(numWorkersEnv); err == nil && n > 0 {
			numWorkers = n
		} else if err != nil {
			log.Printf("Invalid NUM_WORKERS value: %s, using default: %d", numWorkersEnv, DefaultNumWorkers)
		}
	}
	log.Printf("Configured to use %d worker goroutines", numWorkers)

	// Initialize RabbitMQ
	err := consumer.InitRabbitMQ()
	if err != nil {
		log.Fatalf("Failed to initialize RabbitMQ: %v", err)
	}
	defer consumer.CloseRabbitMQ()

	// Start consuming messages
	messages, err := consumer.ConsumeMessages()
	if err != nil {
		log.Fatalf("Failed to consume messages: %v", err)
	}

	// Create and start worker pool
	workerPool := consumer.NewWorkerPool(numWorkers, messages)
	workerPool.Start()

	log.Printf("Get paste worker is running with %d goroutines", numWorkers)
	log.Printf("Press CTRL+C to exit")

	// Set up signal handling for graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	// Wait for shutdown signal
	<-quit
	log.Println("Shutting down worker...")

	// Stop the worker pool
	workerPool.Stop()

	log.Println("Worker shutdown complete")
}
