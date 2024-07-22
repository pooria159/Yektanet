package routes

import (
	"github.com/gin-gonic/gin"
	"go-ad-panel/controllers"
	"go-ad-panel/repositories"
	"gorm.io/gorm"
)

func SetupRouter(db *gorm.DB) *gin.Engine {
	router := gin.Default()

	// Publisher setup
	publisherRepo := repositories.PublisherRepository{Db: db}
	publisherController := controllers.PublisherController{Repo: publisherRepo}

	// Advertiser setup
	advertiserRepo := repositories.AdvertiserRepository{Db: db}
	advertiserController := controllers.AdvertiserController{Repo: advertiserRepo}

	v1 := router.Group("/api/v1")
	{
		// Publisher routes
		publishers := v1.Group("/publishers")
		{
			publishers.POST("", publisherController.CreatePublisher)
			publishers.GET("/:id", publisherController.GetPublisherByID)
			publishers.PUT("/:id", publisherController.UpdatePublisher)
			publishers.DELETE("/:id", publisherController.DeletePublisher)
			publishers.GET("", publisherController.GetAllPublishers)
		}

		// Advertiser routes
		advertisers := v1.Group("/advertisers")
		{
			advertisers.POST("", advertiserController.CreateAdvertiser)
			advertisers.GET("/:id", advertiserController.GetAdvertiserByID)
			advertisers.PUT("/:id", advertiserController.UpdateAdvertiser)
			advertisers.DELETE("/:id", advertiserController.DeleteAdvertiser)
			advertisers.GET("", advertiserController.GetAllAdvertisers)
		}
	}

	return router
}
