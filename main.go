package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	dsn := os.Getenv("DB_DSN")
	if dsn == "" {
		log.Fatal("DB_DSN is not set")
	}

	var db *sql.DB
	var err error

	// Retry loop
	for i := 0; i < 10; i++ {
		db, err = sql.Open("mysql", dsn)
		if err == nil && db.Ping() == nil {
			break
		}
		fmt.Println("⏳ Waiting for MySQL to be ready...")
		time.Sleep(2 * time.Second)
	}
	if err != nil {
		log.Fatalf("❌ Could not connect to MySQL: %v", err)
	}

	fmt.Println("✅ Connected to MySQL")

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	fmt.Printf("🚀 Server running on http://localhost:%s\n", port)

	http.ListenAndServe(":"+port, Routes(db))
}
