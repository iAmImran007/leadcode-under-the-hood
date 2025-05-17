package main

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)



// Struct to hold the DB instance
type Db struct {
	DB *gorm.DB
}

// LoadEnv loads .env file
func LoadEnv() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

// ConnectDb initializes the DB connection and sets it inside the struct
func ConnectDb(db *Db) {
	// Load environment variables
	LoadEnv()

	// Print env values for debugging (optional)
	fmt.Println("Connecting to DB with the following env values:")
	fmt.Println("DB_HOST:", os.Getenv("DB_HOST"))
	fmt.Println("DB_PORT:", os.Getenv("DB_PORT"))

	// Build DSN string
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_PORT"),
		os.Getenv("SSL_MODE"),
	)

	// Connect to PostgreSQL using GORM
	conn, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to the database: %v", err)
	}

	// Save the DB connection
	db.DB = conn

	// Auto migrate your models
	err = db.DB.AutoMigrate(&Problem{}, &TestCase{})
	if err != nil {
		log.Printf("Failed to auto-migrate database: %v", err)
	} else {
		fmt.Println("Database connected and migrated successfully")
	}
}

// GetDb returns the gorm.DB instance
func (d *Db) GetDb() *gorm.DB {
	return d.DB
}

