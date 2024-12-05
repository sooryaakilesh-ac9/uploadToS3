package db

import (
	"os"
	"testing"
)

// TestInitDB tests the InitDB function
func TestInitDB(t *testing.T) {
	// Set up a test database connection string
	os.Setenv("DB_CONN", "host=localhost user=test dbname=test password=test sslmode=disable")

	err := InitDB()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	// Check if dbInstance is initialized
	if dbInstance == nil {
		t.Fatal("expected dbInstance to be initialized")
	}
}

// TestGetDB tests the GetDB function
func TestGetDB(t *testing.T) {
	_, err := GetDB()
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	// Initialize the database
	os.Setenv("DB_CONN", "host=localhost user=test dbname=test password=test sslmode=disable")
	_ = InitDB()

	db, err := GetDB()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if db == nil {
		t.Fatal("expected db to be initialized")
	}
}

// TestConnectToDB tests the ConnectToDB function
func TestConnectToDB(t *testing.T) {
	os.Setenv("DB_CONN", "host=localhost user=test dbname=test password=test sslmode=disable")

	db, err := ConnectToDB()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if db == nil {
		t.Fatal("expected db to be initialized")
	}
}
