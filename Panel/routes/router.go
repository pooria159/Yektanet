package routes

import (
	"github.com/gin-gonic/gin"
	"go-ad-panel/controllers"
	"go-ad-panel/repositories"
	"gorm.io/gorm"
)

func SetupRouter(db *gorm.DB) *gin.Engine {
	router := gin.Default()

	// Initialize the repository
	publisherRepo := repositories.PublisherRepositoryImpl{Db: db}

	publisherController := controllers.PublisherController{Repo: publisherRepo}

	v1 := router.Group("/api/v1")
	{
		publishers := v1.Group("/publishers")
		{
			publishers.POST("", publisherController.CreatePublisher)
			publishers.GET("/:id", publisherController.GetPublisherByID)
			publishers.PUT("/:id", publisherController.UpdatePublisher)
			publishers.DELETE("/:id", publisherController.DeletePublisher)
			publishers.GET("", publisherController.GetAllPublishers)
		}
	}

	return router
}
