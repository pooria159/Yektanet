package controllers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"go-ad-panel/models"
	"go-ad-panel/repositories"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"
)

type AdController struct {
	Repo repositories.AdRepository
}

func (ctrl AdController) GetAllActiveAds(c *gin.Context) {
	ads, err := ctrl.Repo.FindAllActiveAds()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, ads)
}

// CreateAd handles the creation of a new ad with image upload.
func (ctrl AdController) CreateAd(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid advertiser ID"})
		return
	}

	title := c.PostForm("title")
	bid, err := strconv.Atoi(c.PostForm("bid"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid bid value"})
		return
	}

	// Handle file upload
	file, err := c.FormFile("image")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Image upload failed"})
		return
	}

	// Save the file to the media directory
	imagePath := filepath.Join("media", filepath.Base(file.Filename))
	imagePath = strings.ReplaceAll(imagePath, "\\", "/")

	if err := c.SaveUploadedFile(file, imagePath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save image"})
		return
	}

	ad := models.Ad{
		Title:        title,
		ImagePath:    imagePath,
		BidValue:     bid,
		IsActive:     true,
		AdvertiserID: id,
	}

	if err := ctrl.Repo.Save(ad); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Redirect(http.StatusFound, fmt.Sprintf("/advertisers/%d", id))
}
