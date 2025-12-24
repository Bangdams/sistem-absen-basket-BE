package main

import (
	"absen-qr-backend/internal/config"
	"log"
	"os"

	"github.com/go-playground/validator/v10"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	db := config.NewDatabase()
	app := config.NewFiber()
	validate := validator.New()

	config.Bootstrap(&config.BootstrapConfig{
		DB:       db,
		App:      app,
		Validate: validate,
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	if err := app.Listen(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
