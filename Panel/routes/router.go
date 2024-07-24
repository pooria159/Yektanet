package routes

import (
	"github.com/gin-gonic/gin"
	"go-ad-panel/controllers"
	"go-ad-panel/repositories"
	"gorm.io/gorm"
)

func SetupRouter(db *gorm.DB) *gin.Engine {
	router := gin.Default()
	router.LoadHTMLGlob("templates/*")
	router.Static("/static", "./static")
	router.Static("/media", "./media") // Serve uploaded files from the media directory

	// Publisher setup
	publisherRepo := repositories.PublisherRepository{Db: db}
	publisherController := controllers.PublisherController{Repo: publisherRepo}

	// Advertiser setup
	advertiserRepo := repositories.AdvertiserRepository{Db: db}
	advertiserController := controllers.AdvertiserController{Repo: advertiserRepo}

	// Ad setup
	adRepo := repositories.AdRepository{Db: db}
	adController := controllers.AdController{Repo: adRepo}

	router.GET("/publishers/:id", publisherController.PublisherPanel)
	router.GET("/advertisers/:id", advertiserController.AdvertiserPanel)
	router.POST("/advertisers/:id/ad", adController.CreateAd)

	router.POST("/advertisers/:id/charge", advertiserController.ChargeAdvertiser)
	router.POST("/publisher/:id/withdraw", publisherController.PublisherWithdraw)
	router.POST("/ads/:id/toggle", adController.ToggleActivation)

	v1 := router.Group("/api/v1")
	{
		// Publisher routes
		publishers := v1.Group("/publishers")
		{
			publishers.POST("", publisherController.CreatePublisher)
			publishers.PUT("/:id", publisherController.UpdatePublisher)
			publishers.DELETE("/:id", publisherController.DeletePublisher)
			publishers.GET("", publisherController.GetAllPublishers)
		}

		// Advertiser routes
		advertisers := v1.Group("/advertisers")
		{
			advertisers.POST("", advertiserController.CreateAdvertiser)
			advertisers.PUT("/:id", advertiserController.UpdateAdvertiser)
			advertisers.DELETE("/:id", advertiserController.DeleteAdvertiser)
			advertisers.GET("", advertiserController.GetAllAdvertisers)
		}

		// Ad routes
		ads := v1.Group("/ads")
		{
			ads.GET("/active", adController.GetAllActiveAds)
		}
	}

	return router
}
