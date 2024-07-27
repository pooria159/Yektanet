package controllers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"go-ad-panel/models"
	"go-ad-panel/repositories"
	"gorm.io/gorm"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

type AdController struct {
	Repo repositories.AdRepository
	//temp
	RepoAdvertiser repositories.AdvertiserRepository
	RepoPublisher  repositories.PublisherRepository
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

	ctrl.Repo.Update(&ad)
	c.JSON(http.StatusOK, gin.H{})
}

//temp : only for phase 1

type EventRequest struct {
	EventType   string `json:"event_type" binding:"required"`
	PublisherID string `json:"publisher_id" binding:"required"`
}

//	func (ctrl AdController) HandleEvent(c *gin.Context) {
//		// Convert the id parameter from string to int
//		id, err := strconv.Atoi(c.Param("id"))
//		if err != nil {
//			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
//			return
//		}
//
//		ad, err := ctrl.Repo.FindByID(id)
//		if err != nil {
//			c.JSON(http.StatusNotFound, gin.H{"error": "Ad not found"})
//			return
//		}
//
//		// Define a struct to bind JSON fields
//		var eventRequest EventRequest
//
//		// Bind the JSON fields to the struct and check for errors
//		if err := c.ShouldBindJSON(&eventRequest); err != nil {
//			c.JSON(http.StatusBadRequest, gin.H{"error": "here"})
//			return
//		}
//		advertiser, err := ctrl.RepoAdvertiser.FindByID(uint(eventRequest.AdvertiserID))
//		if err != nil {
//			c.JSON(http.StatusNotFound, gin.H{"error": "Advertiser not found"})
//			return
//		}
//		publisher, err := ctrl.RepoPublisher.FindByID(uint(eventRequest.PublisherID))
//		if err != nil {
//			c.JSON(http.StatusNotFound, gin.H{"error": "Publisher not found"})
//			return
//		}
//
//		switch eventRequest.EventType {
//		case "click":
//			ad.Clicks += 1
//			advertiser.Credit -= ad.BidValue
//			publisher.Credit += ad.BidValue
//		case "impression":
//			ad.Impressions += 1
//		default:
//			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid event type"})
//			return
//		}
//		err = ctrl.Repo.Update(&ad)
//		err = ctrl.RepoAdvertiser.Update(&advertiser)
//		err = ctrl.RepoPublisher.Update(&publisher)
//		if err != nil {
//			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to handle error"})
//			return
//		}
//
//		c.JSON(http.StatusOK, gin.H{"message": "Event successfully processed"})
//	}
func (ctrl AdController) HandleEventAtomic(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	err = ctrl.Repo.WithTransaction(func(tx *gorm.DB) error {
		ad, err := ctrl.Repo.FindByIDTx(tx, id)
		if err != nil {
			return err
		}
		advertiser_id := ad.AdvertiserID

		var eventRequest EventRequest
		if err := c.ShouldBindJSON(&eventRequest); err != nil {
			return err
		}

		advertiser, err := ctrl.RepoAdvertiser.FindByIDTx(tx, int(uint(advertiser_id)))
		if err != nil {
			return err
		}
		publisher_id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			return err
		}

		publisher, err := ctrl.RepoPublisher.FindByIDTx(tx, publisher_id)
		if err != nil {
			return err
		}

		switch eventRequest.EventType {
		case "click":
			if err := ctrl.RepoAdvertiser.DecreaseCredit(tx, &advertiser, ad.BidValue); err != nil {
				return err
			}
			if err := ctrl.RepoPublisher.IncreaseCredit(tx, &publisher, ad.BidValue); err != nil {
				return err
			}
			if err := ctrl.Repo.IncrementClicksTx(tx, &ad); err != nil {
				return err
			}
		case "impression":
			if err := ctrl.Repo.IncrementImpressionsTx(tx, &ad); err != nil {
				return err
			}
		}

		return nil
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process event"})
		fmt.Println(err.Error())
	} else {
		c.JSON(http.StatusOK, gin.H{"message": "Event successfully processed"})
	}
}
