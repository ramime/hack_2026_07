package db

import (
	"os"
	"testing"
)

func TestInitDBAndSeed(t *testing.T) {
	dbFile := "test_agencypulse.db"
	defer os.Remove(dbFile)

	database, err := InitDB(dbFile)
	if err != nil {
		t.Fatalf("Failed to init db: %v", err)
	}
	defer database.Close()

	var clientCount int
	if err := database.QueryRow("SELECT COUNT(*) FROM clients").Scan(&clientCount); err != nil {
		t.Fatalf("Failed to query clients count: %v", err)
	}
	if clientCount != 3 {
		t.Errorf("Expected 3 clients, got %d", clientCount)
	}

	var logCount int
	if err := database.QueryRow("SELECT COUNT(*) FROM time_logs").Scan(&logCount); err != nil {
		t.Fatalf("Failed to query time_logs count: %v", err)
	}
	if logCount < 3 {
		t.Errorf("Expected at least 3 initial time logs, got %d", logCount)
	}
}
