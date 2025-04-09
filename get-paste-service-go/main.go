package main

import (
	"log"
	"net/http"
	"os"

	"get-paste-service/config"
	"get-paste-service/db"
	"get-paste-service/router"
)

func main() {
	config.LoadEnv()
	db.InitDB()

	port := os.Getenv("PORT")
	if port == "" {
		port = "8082"
	}

	r := router.SetupRouter()

	log.Printf("Get Paste Service started at :%s", port)
	log.Fatal(http.ListenAndServe(":"+port, r))
}
