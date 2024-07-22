package main

import (
	"go-ad-panel/config"
	"go-ad-panel/models"
	"go-ad-panel/routes"
	"log"
)

func main() {
	config.Connect()
	config.Ping()
	config.Migrate(&models.Publisher{}, &models.Advertiser{}, &models.Ad{})
	router := routes.SetupRouter(config.DB)

	if err := router.Run(":8080"); err != nil {
		log.Fatal("Failed to start server: ", err)
	}
}
