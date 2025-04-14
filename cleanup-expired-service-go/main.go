package main

import (
	"cleanup-expired-service-go/config"
	"cleanup-expired-service-go/db"
	"cleanup-expired-service-go/scheduler"
	"log"
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
