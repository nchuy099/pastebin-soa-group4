package main

import (
	"log"
	"net/http"
	"os"

	"get-public-service/config"
	"get-public-service/db"
	"get-public-service/router"
)

func main() {
	config.LoadEnv()
	db.InitDB()

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	r := router.SetupRouter()

	log.Printf("Get Public Service started at :%s", port)
	log.Fatal(http.ListenAndServe(":"+port, r))
}
