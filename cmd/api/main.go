package main

import (
	"log"
	"golang_v1/internal/config"
	"golang_v1/internal/routes"
)

func main() {
	// Initialize MongoDB
	if err := config.InitMongoDB(); err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}
	defer config.CloseMongoDB()

	// Setup router
	r := routes.SetupRouter()

	// Start server
	log.Println("Server running on :8080")
	if err := r.Run(":8080"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}