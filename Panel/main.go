package main

import (
	"github.com/gin-contrib/cors"
	"go-ad-panel/config"
	"go-ad-panel/models"
	"go-ad-panel/routes"
	"log"
)

func main() {

	config.CreateDB()
	config.Connect()
	config.Ping()

	config.Migrate(&models.Publisher{}, &models.Advertiser{}, &models.Ad{})
	router := routes.SetupRouter(config.DB)
	router.Use(cors.Default())

	if err := router.Run(":8082"); err != nil {
		log.Fatal("Failed to start server: ", err)
	}
}
