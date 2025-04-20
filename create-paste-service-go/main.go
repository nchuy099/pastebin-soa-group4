package main

import (
	"log"
	"net/http"
	"os"

	"create-paste-service-go/cache"
	"create-paste-service-go/config"
	"create-paste-service-go/queue"
	"create-paste-service-go/router"
)

func main() {
	log.Println("Starting paste creation service...")

	// Load environment variables
	config.LoadEnv()

	// Initialize Redis
	err := cache.InitRedis()
	if err != nil {
		log.Printf("Warning: Redis initialization failed: %v", err)
		log.Println("Service will continue without Redis caching")
	} else {
		defer cache.CloseRedis()
	}

	// Initialize RabbitMQ
	err = queue.InitRabbitMQ()
	if err != nil {
		log.Printf("Warning: RabbitMQ initialization failed: %v", err)
		log.Println("Service will continue without async processing")
	} else {
		defer queue.CloseRabbitMQ()
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8081" // Using a different port than the other services
	}

	r := router.SetupRouter()

	log.Printf("Paste Creation Service started at :%s", port)
	log.Fatal(http.ListenAndServe(":"+port, r))
}
