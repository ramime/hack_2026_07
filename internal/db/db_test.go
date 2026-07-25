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

	// Test Client Portal Queries
	client, err := database.GetClientByToken("ritter-sport-8821")
	if err != nil || client == nil {
		t.Fatalf("Failed to fetch client by token: %v", err)
	}
	if client.Name != "Ritter Sport" || client.PinCode != "1234" {
		t.Errorf("Unexpected client data: %+v", client)
	}

	assets, err := database.GetClientContentAssets(client.ID)
	if err != nil {
		t.Fatalf("Failed to fetch client content assets: %v", err)
	}
	if len(assets) == 0 {
		t.Errorf("Expected content assets for Ritter Sport, got 0")
	}
}
