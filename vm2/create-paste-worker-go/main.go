package main

import (
	"log"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	"create-paste-worker/config"
	"create-paste-worker/consumer"
	"create-paste-worker/db"
	"create-paste-worker/repository"
)

func main() {
	log.Println("Starting paste creation worker service...")

	// Load environment variables
	config.LoadEnv()

	// Initialize database connection
	db.InitDB()
	defer db.DB.Close()

	// Initialize RabbitMQ
	err := consumer.InitRabbitMQ()
	if err != nil {
		log.Fatalf("Failed to initialize RabbitMQ: %v", err)
	}
	defer consumer.CloseRabbitMQ()

	// Configure batch processing
	configureBatchSettings()

	// Get number of workers from environment
	numWorkersStr := os.Getenv("NUM_WORKERS")
	numWorkers := 5 // Default to 5 workers
	if numWorkersStr != "" {
		if n, err := strconv.Atoi(numWorkersStr); err == nil && n > 0 {
			numWorkers = n
		}
	}

	log.Printf("Starting worker pool with %d workers", numWorkers)

	// Start consuming messages
	messages, err := consumer.ConsumeMessages()
	if err != nil {
		log.Fatalf("Failed to consume messages: %v", err)
	}

	// Create and start worker pool
	workerPool := consumer.NewWorkerPool(numWorkers, messages)
	workerPool.Start()

	// Wait for termination signal
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	log.Println("Shutting down worker service...")
	workerPool.Stop()
}

func configureBatchSettings() {
	// Get batch size from environment
	batchSizeStr := os.Getenv("BATCH_SIZE")
	if batchSizeStr != "" {
		if batchSize, err := strconv.Atoi(batchSizeStr); err == nil && batchSize > 0 {
			repository.BatchSize = batchSize
			log.Printf("Configured batch size: %d", batchSize)
		}
	} else {
		log.Printf("Using default batch size: %d", repository.BatchSize)
	}

	// Get batch timeout from environment
	batchTimeoutStr := os.Getenv("BATCH_TIMEOUT_SECS")
	if batchTimeoutStr != "" {
		if timeout, err := strconv.Atoi(batchTimeoutStr); err == nil && timeout > 0 {
			repository.BatchTimeoutSecs = timeout
			log.Printf("Configured batch timeout: %d seconds", timeout)
		}
	} else {
		log.Printf("Using default batch timeout: %d seconds", repository.BatchTimeoutSecs)
	}
}
