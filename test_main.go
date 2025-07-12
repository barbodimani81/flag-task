// main_test.go
package main

import (
	"database/sql"
	"log"
	"os"
	"testing"

	_ "github.com/go-sql-driver/mysql"
)

// shared test DB connection
var testDB *sql.DB

func TestMain(m *testing.M) {
	dsn := os.Getenv("TEST_DB_DSN")
	if dsn == "" {
		log.Fatal("TEST_DB_DSN env var not set")
	}

	var err error
	testDB, err = sql.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("failed to open DB: %v", err)
	}

	if err := testDB.Ping(); err != nil {
		log.Fatalf("failed to ping DB: %v", err)
	}

	code := m.Run()

	testDB.Close()
	os.Exit(code)
}

// helper to clear tables before each test
func resetTestDB(t *testing.T) {
	t.Helper()
	tables := []string{"audit_logs", "dependencies", "feature_flags"}
	for _, tbl := range tables {
		_, err := testDB.Exec("DELETE FROM " + tbl)
		if err != nil {
			t.Fatalf("failed to reset table %s: %v", tbl, err)
		}
	}
}

func TestCreateFeatureFlag(t *testing.T) {
	resetTestDB(t)

	err := createFeatureFlag(testDB, "auth", nil)
	if err != nil {
		t.Fatalf("failed to create flag: %v", err)
	}

	var count int
	err = testDB.QueryRow(`SELECT COUNT(*) FROM feature_flags WHERE name = ?`, "auth").Scan(&count)
	if err != nil {
		t.Fatalf("query error: %v", err)
	}
	if count != 1 {
		t.Fatalf("expected 1 flag named 'auth', got %d", count)
	}
}
