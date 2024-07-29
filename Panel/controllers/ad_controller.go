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
	RepoAdvertiser repositories.AdvertiserRepository
	RepoPublisher  repositories.PublisherRepository
}

//IS Okey
func (ctrl AdController) GetAllActiveAds(c *gin.Context) {
	ads, err := ctrl.Repo.FindAllActiveAds()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, ads)
}

//IS Okey
func (ctrl AdController) CreateAd(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil || id <= 0 {
		c.HTML(http.StatusBadRequest, "advertiser.html", gin.H{"notfounderror": "Invalid Advertiser ID"})
		return
	}

	title := c.PostForm("title")
	bid ,_ := strconv.Atoi(c.PostForm("bid"))
	file, err := c.FormFile("image")
	if err != nil {
		c.HTML(http.StatusBadRequest, "advertiser.html", gin.H{"notfounderror": "Image Upload Failed"})
		return
	}

	// Generate new filename with timestamp and advertiser ID
	ext := filepath.Ext(file.Filename)
	name := strings.TrimSuffix(file.Filename, ext)
	timestamp := time.Now().Format("20060102150405")
	newFilename := fmt.Sprintf("%s_%s_%d%s", name, timestamp, id, ext)

	// Save the file to the media directory
	imagePath := filepath.Join("media", newFilename)
	imagePath = strings.ReplaceAll(imagePath, "\\", "/")
	if err := c.SaveUploadedFile(file, imagePath); err != nil {
		c.HTML(http.StatusInternalServerError, "advertiser.html", gin.H{"notfounderror": "Failed To Save Image"})
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
		c.HTML(http.StatusInternalServerError, "advertiser.html", gin.H{"notfounderror": "The Ad Was Not Created"})
		return
	}

	c.HTML(http.StatusOK, "advertiser.html", gin.H{"adsuccess": "Ad Created Successfully", "ad": ad})
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

type EventRequest struct {
	EventType   string `json:"event_type" binding:"required"`
	PublisherID string `json:"publisher_id" binding:"required"`
}

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
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		fmt.Println(err.Error())
	} else {
		c.JSON(http.StatusOK, gin.H{"message": "Event successfully processed"})
	}
}
