package controllers

import (
	"github.com/gin-gonic/gin"
	"go-ad-panel/models"
	"go-ad-panel/repositories"
	"net/http"
	"strconv"
)

type PublisherController struct {
	Repo repositories.PublisherRepository
}

func (ctrl PublisherController) CreatePublisher(c *gin.Context) {
	var publisher models.Publisher
	if err := c.ShouldBindJSON(&publisher); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := ctrl.Repo.Save(publisher); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, publisher)
}

// GetPublisherByID handles fetching a publisher by its ID
func (ctrl PublisherController) GetPublisherByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	publisher, err := ctrl.Repo.FindByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Publisher not found"})
		return
	}
	c.JSON(http.StatusOK, publisher)
}

// UpdatePublisher handles updating an existing publisher
func (ctrl PublisherController) UpdatePublisher(c *gin.Context) {
	var publisher models.Publisher
	if err := c.ShouldBindJSON(&publisher); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := ctrl.Repo.Update(publisher); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, publisher)
}

func (ctrl PublisherController) DeletePublisher(c *gin.Context) {
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

// GetAllPublishers handles fetching all publishers
func (ctrl PublisherController) GetAllPublishers(c *gin.Context) {
	publishers, err := ctrl.Repo.FindAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, publishers)
}
