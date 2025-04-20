package main

import (
	"log"
	"os"
	"public-service/handlers"
	"public-service/models"

	"github.com/gin-gonic/gin"
)

func main() {
	// Initialize DB connection
	db, err := models.ConnectDB()
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer db.Close()
	models.DB = db

	// Create a Gin router with HTML template rendering enabled.
	router := gin.Default()
	router.LoadHTMLGlob("templates/*.html")

	// Define routes.
	router.GET("/public", handlers.GetPublicPastes)

	// Start the server.
	port := os.Getenv("PORT")
	if port == "" {
		port = "3001"
	}
	log.Printf("Public service running on port %s\n", port)
	router.Run(":" + port)
}
