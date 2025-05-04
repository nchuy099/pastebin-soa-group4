package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

// DB is the global database connection
var DB *sql.DB

// InitDB initializes the database connection
func InitDB() {
	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?parseTime=true",
		os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_HOST"), os.Getenv("DB_NAME"))

	var err error
	DB, err = sql.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("Failed to connect to MySQL: %v", err.Error())
	}

<<<<<<<< HEAD:microservice/create-paste-service-go/db/mysql.go
	DB.SetMaxOpenConns(200)    // tổng số kết nối tối đa (active hoặc idle)
	DB.SetMaxIdleConns(80)     // số kết nối nhàn rỗi (idle) giữ lại
	DB.SetConnMaxLifetime(300) // lifetime (seconds) của mỗi connection, tránh timeout ngẫu nhiên
========
	// DB.SetMaxOpenConns(36)     // Maximum number of open connections (active or idle)
	// DB.SetMaxIdleConns(18)     // Maximum number of idle connections
	// DB.SetConnMaxLifetime(300) // Maximum lifetime (seconds) of each connection to avoid random timeouts
>>>>>>>> ms-rabbitmq-redis:create-paste-worker-go/db/mysql.go

	if err = DB.Ping(); err != nil {
		log.Fatalf("Database unreachable: %v", err.Error())
	}

	log.Println("Connected to MySQL successfully")
}
