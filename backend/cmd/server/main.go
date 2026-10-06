package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/vosram/filesender/backend/internal/api"
)

func main() {
	// Create a Gin router with default middleware (logger and recovery)
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	dbConnString := os.Getenv("DB_CONNECTION_STRING")
	if dbConnString == "" {
		log.Fatal("DB connection not present")
	}
	jwtSecretStr := os.Getenv("JWT_SECRET")
	if jwtSecretStr == "" {
		log.Fatal("JWT_SECRET is not present")
	}
	platform := os.Getenv("PLATFORM")
	if platform != "production" {
		platform = "development"
	}
	r := gin.Default()

	apiConfig, err := api.New(dbConnString, jwtSecretStr, "")

	// Auth Routes
	r.POST("/api/auth/email/login", apiConfig.EmailLogin)
	r.POST("/api/auth/email/verify", apiConfig.EmailLoginVerify)
	// Start server on port 8080 (default)
	// Server will listen on 0.0.0.0:8080 (localhost:8080 on Windows)
	r.Run()
}
