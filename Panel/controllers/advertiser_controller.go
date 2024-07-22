package controllers

import (
	"github.com/gin-gonic/gin"
	"go-ad-panel/models"
	"go-ad-panel/repositories"
	"net/http"
	"strconv"
)

type AdvertiserController struct {
	Repo repositories.AdvertiserRepository
}

func (ctrl AdvertiserController) CreateAdvertiser(c *gin.Context) {
	var advertiser models.Advertiser
	if err := c.ShouldBindJSON(&advertiser); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := ctrl.Repo.Save(advertiser); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, advertiser)
}

func (ctrl AdvertiserController) GetAdvertiserByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	advertiser, err := ctrl.Repo.FindByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Advertiser not found"})
		return
	}
	c.JSON(http.StatusOK, advertiser)
}

func (ctrl AdvertiserController) UpdateAdvertiser(c *gin.Context) {
	var advertiser models.Advertiser
	if err := c.ShouldBindJSON(&advertiser); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := ctrl.Repo.Update(advertiser); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, advertiser)
}

func (ctrl AdvertiserController) DeleteAdvertiser(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	if err := ctrl.Repo.Delete(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusNoContent, nil)
}

func (ctrl AdvertiserController) GetAllAdvertisers(c *gin.Context) {
	advertisers, err := ctrl.Repo.FindAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, advertisers)
}