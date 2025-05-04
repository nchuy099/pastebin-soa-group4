package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

// DB is the global database connection
var DB *sql.DB

// InitDB initializes the database connection with ProxySQL configuration
func InitDB() {
	// Build DSN with ProxySQL parameters
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?timeout=10s&readTimeout=30s&writeTimeout=30s&parseTime=true",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_NAME"))

	var err error
	DB, err = sql.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("Failed to connect to MySQL via ProxySQL: %v", err)
	}

	// Set connection pool parameters optimized for ProxySQL
	DB.SetMaxOpenConns(30)
	DB.SetMaxIdleConns(15)
	DB.SetConnMaxLifetime(5 * time.Minute)
	DB.SetConnMaxIdleTime(3 * time.Minute)

	// Verify connection
	if err = DB.Ping(); err != nil {
		log.Fatalf("Database unreachable through ProxySQL: %v", err)
	}

	log.Println("Connected to MySQL via ProxySQL successfully")
}

// CloseDB closes the database connection
func CloseDB() {
	if DB != nil {
		DB.Close()
		log.Println("Database connection closed")
	}
}
