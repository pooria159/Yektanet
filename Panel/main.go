package main

import (
    // "go-ad-panel/config"
    // "go-ad-panel/models"
    "go-ad-panel/routes"
)

func main() {
    // config.Connect()
    // config.Migrate(&models.Advertiser{}, &models.Publisher{}, &models.Ad{}, &models.Impression{}, &models.Click{})
    r := routes.SetupRouter()
    r.Run(":8080")
}
