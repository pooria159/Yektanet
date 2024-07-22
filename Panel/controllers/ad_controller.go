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
	"time"
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

	// Generate new filename with timestamp and advertiser ID
	ext := filepath.Ext(file.Filename)
	name := strings.TrimSuffix(file.Filename, ext)
	timestamp := time.Now().Format("20060102150405")
	newFilename := fmt.Sprintf("%s_%s_%d%s", name, timestamp, id, ext)

	// Save the file to the media directory
	imagePath := filepath.Join("media", newFilename)
	// Convert the imagePath to use forward slashes instead of backslashes
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

func (ctrl AdController) ToggleActivation(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}
	ad, err := ctrl.Repo.FindByID(id)
	if err != nil {
		c.HTML(http.StatusNotFound, "notfound.html", gin.H{"error": "Ad not found"})
		return
	}
	ad.IsActive = !ad.IsActive

	ctrl.Repo.Update(ad)
	c.JSON(http.StatusOK, gin.H{})
}
