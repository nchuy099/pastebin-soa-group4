package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

var DB *sql.DB

func InitDB() {
	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?parseTime=true",
		os.Getenv("DB_USER"), os.Getenv("DB_PASS"), os.Getenv("DB_HOST"), os.Getenv("DB_NAME"))

	var err error
	DB, err = sql.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("Failed to connect to MySQL: %v", err)
	}

	// DB.SetMaxOpenConns(12)     // tổng số kết nối tối đa (active hoặc idle)
	// DB.SetMaxIdleConns(6)      // số kết nối nhàn rỗi (idle) giữ lại
	// DB.SetConnMaxLifetime(300) // lifetime (seconds) của mỗi connection, tránh timeout ngẫu nhiên

	if err = DB.Ping(); err != nil {
		log.Fatalf("Database unreachable: %v", err)
	}

	log.Println("Connected to MySQL successfully")
}
