package db

import (
	"backend/pkg/quotes"
	"fmt"
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	dbInstance *gorm.DB
)

// init db and the instance is reused
// call to ConnectToDB and get the connection instance
func InitDB() error {
	var err error
	dbInstance, err = ConnectToDB()
	if err != nil {
		return err
	}
	return nil
}

func GetDB() (*gorm.DB, error) {
	if dbInstance == nil {
		return nil, fmt.Errorf("database connection not initialized, call initDB() first")
	}
	return dbInstance, nil
}

// connect to DB and returns an instance of the DB
// ConnectToDB connects to the PostgreSQL database and ensures the necessary tables exist
func ConnectToDB() (*gorm.DB, error) {
	// Database connection string (todo: move to environment variables)
	dsn := os.Getenv("DB_CONN")

	// Open a DB connection
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to the database: %w", err)
	}

	// Migrate the 'quotes' table (if it doesn't exist, AutoMigrate will create it)
	if err := db.AutoMigrate(&quotes.Quote{}); err != nil {
		log.Printf("Failed to migrate 'quotes' table: %v", err)
		return nil, err
	} else {
		log.Println("Table 'quotes' migration successful or already exists.")
	}

	// Migrate the 'flyers' table (if it doesn't exist, AutoMigrate will create it)
	if err := db.AutoMigrate(&quotes.Quote{}); err != nil {
		log.Printf("Failed to migrate 'flyers' table: %v", err)
		return nil, err
	} else {
		log.Println("Table 'flyers' migration successful or already exists.")
	}

	// Test the connection by executing a simple query
	if err := db.Exec("SELECT 1").Error; err != nil {
		return nil, fmt.Errorf("failed to ping the database: %w", err)
	}

	// Successfully connected
	fmt.Println("Successfully connected to the PostgreSQL database!")

	// Return the DB instance for further usage
	return db, nil
}
