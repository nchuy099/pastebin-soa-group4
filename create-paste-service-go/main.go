package main

import (
	"log"
	"net/http"
	"os"

	"create-paste-service-go/config"
	"create-paste-service-go/db"
	"create-paste-service-go/router"
)

func main() {
	config.LoadEnv()
	db.InitDB()

	port := os.Getenv("PORT")
	if port == "" {
		port = "8081" // Using a different port than the other services
	}

	r := router.SetupRouter()

	log.Printf("Paste Creation Service started at :%s", port)
	log.Fatal(http.ListenAndServe(":"+port, r))
}
