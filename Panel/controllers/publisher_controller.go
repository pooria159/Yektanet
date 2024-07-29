package controllers

import (
	// "fmt"
	"github.com/gin-gonic/gin"
	"go-ad-panel/models"
	"go-ad-panel/repositories"
	"net/http"
	"strconv"
)

// type PublisherController struct {
// 	Repo repositories.PublisherRepository
// }

type PublisherController struct {
	Repo repositories.PublisherRepositoryInterface
}

// Is Okey
func (ctrl PublisherController) PublisherPanel(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil || id <= 0 {
		c.HTML(http.StatusBadRequest, "publisher.html", gin.H{"error": "Invalid ID"})
		return
	}
	publisher, err := ctrl.Repo.FindByID(uint(id))
	if err != nil || publisher.ID == 0 {
		c.HTML(http.StatusNotFound, "publisher.html", gin.H{"error": "Publisher Not Found"})
		return
	}
	c.HTML(http.StatusOK, "publisher.html", gin.H{"publisher": publisher})
}

// IS Okey
func (ctrl PublisherController) PublisherWithdraw(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil || id <= 0 {
		c.HTML(http.StatusBadRequest, "publisher.html", gin.H{"error": "Invalid ID"})
		return
	}
	publisher, err := ctrl.Repo.FindByID(uint(id))
	if err != nil || publisher.ID == 0 {
		c.HTML(http.StatusNotFound, "publisher.html", gin.H{"error": "Publisher Not Found"})
		return
	}

	amountStr := c.PostForm("amount")
	amount, err := strconv.ParseFloat(amountStr, 64)
	if err != nil || amount <= 0 {
		c.HTML(http.StatusBadRequest, "publisher.html", gin.H{"publisher": publisher, "error": "Invalid Amount" })
		return
	}

	if amount != float64(int(amount)) {
        c.HTML(http.StatusBadRequest, "publisher.html", gin.H{"publisher": publisher, "error": "Amount must be an integer"})
        return
    }

	intAmount := int(amount)

	if publisher.Credit >= intAmount  {
		publisher.Credit -= intAmount
		err := ctrl.Repo.Update(&publisher)
		if err != nil {
			c.HTML(http.StatusInternalServerError, "publisher.html", gin.H{"publisher": publisher, "error": "Internal Server Error"})
			return
		}
		c.HTML(http.StatusOK, "publisher.html", gin.H{"publisher": publisher, "success": "Withdrawal Successful"})
	} else {
		c.HTML(http.StatusBadRequest, "publisher.html", gin.H{"publisher": publisher, "error": "Insufficient Balance"})
	}
}

// IS Okey
func (ctrl PublisherController) CreatePublisher(c *gin.Context) {
	var publisher models.Publisher
	if err := c.ShouldBindJSON(&publisher); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := ctrl.Repo.Save(&publisher); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, publisher)
}

// IS Okey
func (ctrl PublisherController) UpdatePublisher(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	if _, err := ctrl.Repo.FindByID(uint(id)); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "publisher not found"})
		return
	}
	var publisher models.Publisher
	if err := c.ShouldBindJSON(&publisher); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	publisher.ID = uint(id)
	if err := ctrl.Repo.Update(&publisher); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, publisher)
}

// IS Okey
func (ctrl PublisherController) DeletePublisher(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	if err := ctrl.Repo.Delete(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusNoContent, gin.H{"message": "publisher deleted successfully"})
}

// IS Okey
func (ctrl PublisherController) GetAllPublishers(c *gin.Context) {
	publishers, err := ctrl.Repo.FindAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, publishers)
}
