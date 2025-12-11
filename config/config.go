package config

import (
    "fmt"
    "log"
    "os"

    "gorm.io/driver/postgres"
    "gorm.io/gorm"

	"github.com/followCode/djjs-event-reporting-backend/app/models"
)


var JWTSecret []byte

func LoadJWTSecret() {
    secret := os.Getenv("JWT_SECRET")
    if secret == "" {
        log.Fatal("JWT_SECRET is not set in environment")
    }
    JWTSecret = []byte(secret)
}

var DB *gorm.DB

func ConnectDB() {
    dbUser := os.Getenv("POSTGRES_USER")
    dbPass := os.Getenv("POSTGRES_PASSWORD")
    dbName := os.Getenv("POSTGRES_DB")
    dbPort := os.Getenv("PG_PORT")
    dbHost := os.Getenv("POSTGRES_HOST")

    // Validate required environment variables
    if dbHost == "" {
        log.Fatal("POSTGRES_HOST is required in .env or environment variables")
    }
    if dbUser == "" {
        log.Fatal("POSTGRES_USER is required in .env or environment variables")
    }
    if dbPass == "" {
        log.Fatal("POSTGRES_PASSWORD is required in .env or environment variables")
    }
    if dbName == "" {
        log.Fatal("POSTGRES_DB is required in .env or environment variables")
    }
    if dbPort == "" {
        dbPort = "5432" // Default PostgreSQL port
    }

	log.Printf("Connecting to DB -> host=%s port=%s user=%s dbname=%s", dbHost, dbPort, dbUser, dbName)

    dsn := fmt.Sprintf(
        "host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
        dbHost, dbUser, dbPass, dbName, dbPort,
    )

    db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
    if err != nil {
        log.Fatal("Failed to connect to DB:", err)
    }

    DB = db
}

func AutoMigrate() {
    DB.AutoMigrate(&models.Role{}, &models.User{})
}
