package main

import (
	"go-ad-panel/config"
	"go-ad-panel/models"
	"go-ad-panel/routes"
	"log"
)

func main() {
	config.Connect()

	// Ping the database to ensure the connection is alive
	config.Ping()

	// Automatically migrate your schema
	config.Migrate(&models.Publisher{})

	// Set up the router
	router := routes.SetupRouter(config.DB)

	// Start the server
	if err := router.Run(":8080"); err != nil {
		log.Fatal("Failed to start server: ", err)
	}
}
