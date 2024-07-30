package controllers

import (
	"github.com/gin-gonic/gin"
	"go-ad-panel/models"
	"go-ad-panel/repositories"
	"net/http"
	"strconv"
)

// type AdvertiserController struct {
// 	Repo repositories.AdvertiserRepository
// }

type AdvertiserController struct {
	Repo repositories.AdvertiserRepositoryInterface
}

// IS Okey
func (ctrl AdvertiserController) AdvertiserPanel(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil || id <= 0 {
		c.HTML(http.StatusBadRequest, "advertiser.html", gin.H{"notfounderror": "Invalid ID"})
		return
	}
	advertiser, ads, err := ctrl.Repo.FindByIDWithAds(uint(id))
	if err != nil || advertiser.ID == 0 {
		c.HTML(http.StatusNotFound, "advertiser.html", gin.H{"notfounderror": "Advertiser not found"})
		return
	}
	c.HTML(http.StatusOK, "advertiser.html", gin.H{"advertiser": advertiser, "ads": ads})
}

// IS Okey
func (ctrl AdvertiserController) CreateAdvertiser(c *gin.Context) {
	var advertiser models.Advertiser
	if err := c.ShouldBindJSON(&advertiser); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := ctrl.Repo.Save(&advertiser); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, advertiser)
}

// IS Okey
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

// IS Okey
func (ctrl AdvertiserController) UpdateAdvertiser(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}
	var advertiser models.Advertiser
	if err := c.ShouldBindJSON(&advertiser); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	advertiser.ID = uint(id)
	if _, err := ctrl.Repo.FindByID(uint(id)); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Advertiser not found"})
		return
	}
	if err := ctrl.Repo.Update(&advertiser); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, advertiser)
}

// IS Okey
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

// IS Okey
func (ctrl AdvertiserController) GetAllAdvertisers(c *gin.Context) {
	advertisers, err := ctrl.Repo.FindAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, advertisers)
}

// IS Okey
func (ctrl AdvertiserController) ChargeAdvertiser(c *gin.Context) {
    id, err := strconv.Atoi(c.Param("id"))
    if err != nil || id <= 0 {
        c.HTML(http.StatusBadRequest, "advertiser.html", gin.H{"notfounderror": "Invalid ID"})
        return
    }

    advertiser, err := ctrl.Repo.FindByID(uint(id))
    if err != nil || advertiser.ID == 0 {
        c.HTML(http.StatusNotFound, "advertiser.html", gin.H{"notfounderror": "Advertiser Not Found"})
        return
    }

    amountStr := c.PostForm("amount")
    amount, err := strconv.ParseFloat(amountStr, 64)
    if err != nil || amount <= 0 {
        c.HTML(http.StatusBadRequest, "advertiser.html", gin.H{"advertiser": advertiser, "error": "Invalid Amount"})
        return
    }

	if amount != float64(int(amount)) {
        c.HTML(http.StatusBadRequest, "advertiser.html", gin.H{"advertiser": advertiser, "error": "Charge must be an integer"})
        return
    }

	intAmount := int(amount)

    advertiser.Credit += intAmount
    err = ctrl.Repo.Update(&advertiser)
    if err != nil {
        c.HTML(http.StatusInternalServerError, "advertiser.html", gin.H{"advertiser": advertiser, "error": "Internal Server Error"})
        return
    }

    c.HTML(http.StatusOK, "advertiser.html", gin.H{"advertiser": advertiser, "success": "Charge Successful"})
}
