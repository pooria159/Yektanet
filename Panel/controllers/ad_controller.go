package controllers

import (
	"github.com/gin-gonic/gin"
	"go-ad-panel/repositories"
	"net/http"
)

type AdController struct {
	Repo repositories.AdRepository
}

// GetAllActiveAds handles fetching all active ads
func (ctrl AdController) GetAllActiveAds(c *gin.Context) {
	ads, err := ctrl.Repo.FindAllActiveAds()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, ads)
}
