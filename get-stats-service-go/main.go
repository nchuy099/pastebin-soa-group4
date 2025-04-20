package main

import (
	"log"
	"net/http"
	"os"

	"get-stats-service-go/config"
	"get-stats-service-go/db"
	"get-stats-service-go/router"
)

func main() {
	config.LoadEnv()
	db.InitDB()

	port := os.Getenv("PORT")
	if port == "" {
		port = "8083" // Using a different port than the other service
	}

	r := router.SetupRouter()

	log.Printf("Get Stats Service started at :%s", port)
	log.Fatal(http.ListenAndServe(":"+port, r))
}
