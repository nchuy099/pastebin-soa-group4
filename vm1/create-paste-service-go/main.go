package main

import (
	"log"
	"net/http"
	"os"

	"create-paste-service/cache"
	"create-paste-service/config"
	"create-paste-service/producer"
	"create-paste-service/router"
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
	err = producer.InitRabbitMQ()
	if err != nil {
		log.Printf("Warning: RabbitMQ initialization failed: %v", err)
		log.Println("Service will continue without async processing")
	} else {
		defer producer.CloseRabbitMQ()
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8081" // Using a different port than the other services
	}

	r := router.SetupRouter()

	log.Printf("Paste Creation Service started at :%s", port)
	log.Fatal(http.ListenAndServe(":"+port, r))
}
