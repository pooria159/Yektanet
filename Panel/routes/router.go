package routes

import (
    // "go-ad-panel/controllers"
    "github.com/gin-gonic/gin"
    "net/http"
    // "go-ad-panel/config"
)

func SetupRouter() *gin.Engine {
    r := gin.Default()
    // db := config.DB
    r.GET("", func(context *gin.Context) {
		context.JSON(http.StatusOK, "welcome home")
	})

	r.NoRoute(func(c *gin.Context) {
		c.JSON(404, gin.H{"code": "PAGE_NOT_FOUND", "message": "Page not found"})
	})

    // advertiserRepo := repositories.NewAdvertiserRepository(db)
    // advertiserCtrl := &controllers.AdvertiserController{Repository: advertiserRepo}
    // r.GET("/panel/advertiser/:id", advertiserCtrl.GetAdvertiser)
    // r.POST("/panel/advertiser", advertiserCtrl.CreateAdvertiser)
    return r
}



// package router

// import (
// 	"github.com/gin-gonic/gin"
// 	"net/http"
// 	"practice-api-gin-one/controller"
// )

// func NewRouter(tagController *controller.TagController) *gin.Engine {
// 	service := gin.Default()

// 	service.GET("", func(context *gin.Context) {
// 		context.JSON(http.StatusOK, "welcome home")
// 	})

// 	service.NoRoute(func(c *gin.Context) {
// 		c.JSON(404, gin.H{"code": "PAGE_NOT_FOUND", "message": "Page not found"})
// 	})

// 	router := service.Group("/api")
// 	tagRouter := router.Group("/tag")
// 	tagRouter.GET("", tagController.FindAll)
// 	tagRouter.GET("/:tagId", tagController.FindById)
// 	tagRouter.POST("", tagController.Create)
// 	tagRouter.PATCH("/:tagId", tagController.Update)
// 	tagRouter.DELETE("/:tagId", tagController.Delete)

// 	return service
// }