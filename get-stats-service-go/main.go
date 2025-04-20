package main

import (
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"get-stats-service/cache"
	"get-stats-service/config"
	"get-stats-service/db"
	"get-stats-service/router"
)

func main() {
	config.LoadEnv()
	db.InitDB()
	cache.InitRedis()

	port := os.Getenv("PORT")
	if port == "" {
		port = "8083" // Using a different port than the other service
	}

	r := router.SetupRouter()

	// Set up graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	// Start server in a goroutine
	go func() {
		log.Printf("Get Stats Service started at :%s", port)
		if err := http.ListenAndServe(":"+port, r); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server error: %v", err)
		}
	}()

	// Wait for shutdown signal
	<-quit
	log.Println("Shutting down server...")
	log.Println("Server shutdown complete")
}
