package main

import (
	"cleanup-expired-service-go/config"
	"cleanup-expired-service-go/db"
	"cleanup-expired-service-go/scheduler"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"
)

func main() {
	// Load environment variables
	config.LoadEnv()

	// Initialize database connection
	db.InitDB()
	defer db.DB.Close()

	// Get port from environment variable
	port := os.Getenv("PORT")
	if port == "" {
		port = "8084" // Default port
	}

	// Get cleanup interval from environment variable (in minutes)
	cleanupIntervalStr := os.Getenv("CLEANUP_INTERVAL_MINS")
	if cleanupIntervalStr == "" {
		cleanupIntervalStr = "0.25" // Default to 60 minutes if not specified
	}

	cleanupIntervalFloat, err := strconv.ParseFloat(cleanupIntervalStr, 64)
	if err != nil {
		log.Printf("Invalid CLEANUP_INTERVAL_MINUTES value: %v. Using default of 60 minutes", err)
		cleanupIntervalFloat = 0.25
	}

	cleanupDuration := time.Duration(cleanupIntervalFloat * float64(time.Minute))

	// Create a ticker for periodic cleanup
	ticker := time.NewTicker(cleanupDuration)
	defer ticker.Stop()

	// Setup HTTP server
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Cleanup service is running"))
	})

	// Start the HTTP server in a goroutine
	go func() {
		log.Printf("Starting HTTP server on port %s", port)
		if err := http.ListenAndServe(":"+port, nil); err != nil {
			log.Fatalf("HTTP server error: %v", err)
		}
	}()

	// Run cleanup immediately on startup
	log.Println("Starting expired pastes cleanup service...")
	count, err := scheduler.CleanupExpiredPastes()
	if err != nil {
		log.Printf("Error during initial cleanup: %v", err)
	} else {
		log.Printf("Initial cleanup completed: %d expired pastes deleted", count)
	}

	// Setup signal handling for graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// Main loop
	log.Printf("Cleanup scheduler running. Will check for expired pastes every %.2f minutes", cleanupIntervalFloat)
	for {
		select {
		case <-ticker.C:
			count, err := scheduler.CleanupExpiredPastes()
			if err != nil {
				log.Printf("Error during scheduled cleanup: %v", err)
			} else {
				log.Printf("Scheduled cleanup completed: %d expired pastes deleted", count)
			}
		case sig := <-sigChan:
			log.Printf("Received signal %v, shutting down...", sig)
			return
		}
	}
}
