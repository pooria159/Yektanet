package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

var publisherData map[string]string

func init() {
	publisherData = map[string]string{
		"publisher1": "This is a unique text for Publisher 1 website.",
		"publisher2": "This is a unique text for Publisher 2 website.",
		"publisher3": "This is a unique text for Publisher 3 website.",
		"publisher4": "This is a unique text for Publisher 4 website.",
		"publisher5": "This is a unique text for Publisher 5 website.",
	}
}

func main() {
	router := gin.Default()
	//router.Use(cors.Default())
	router.LoadHTMLGlob("templates/*")
	router.Static("/static", "./static")
	router.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "Hello World"})
		return
	})
	router.GET("/:publisher", func(c *gin.Context) {
		publisher := c.Param("publisher")
		text, exists := publisherData[publisher]
		if !exists {
			c.JSON(http.StatusNotFound, gin.H{"error": "Publisher not found"})
			return
		}
		c.HTML(http.StatusOK, publisher+".html", gin.H{
			"text": text,
		})
	})

	router.Run(":9123")
}
