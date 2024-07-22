package controllers

import (
	"github.com/gin-gonic/gin"
	"go-ad-panel/models"
	"go-ad-panel/repositories"
	"net/http"
	"strconv"
	"fmt"

)

type PublisherController struct {
	Repo repositories.PublisherRepository
}

// Is Okey
func (ctrl PublisherController) PublisherPanel(c *gin.Context) {
    id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}
	publisher, err := ctrl.Repo.FindByID(uint(id))
    if err != nil {
        c.HTML(http.StatusNotFound, "notfound.html", gin.H{"error": "Publisher not found"})
        return
    }
    c.HTML(http.StatusOK, "publisher.html", gin.H{"publisher": publisher})
}


func (ctrl PublisherController) PublisherWithdraw(c *gin.Context) {
    id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}
	publisher, err := ctrl.Repo.FindByID(uint(id))
    if err != nil {
        c.HTML(http.StatusNotFound, "notfound.html", gin.H{"error": "Publisher not found"})
        return
    }
	amountStr := c.PostForm("amount")
	amount, err := strconv.ParseFloat(amountStr, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid amount"})
		return
	}
	if publisher.Credit >= int(amount) {
        publisher.Credit -= int(amount)
        ctrl.Repo.Update(publisher)
    } else {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Insufficient balance"})
        return
    }
	fmt.Println(id)
    c.Redirect(http.StatusSeeOther, fmt.Sprintf("/publishers/%d", id))
}


// func publisherWithdraw(c *gin.Context) {
//     id := c.Param("id")
//     var publisher Publisher
//     if err := DB.First(&publisher, id).Error; err != nil {
//         c.JSON(http.StatusNotFound, gin.H{"error": "Publisher not found"})
//         return
//     }

//     amountStr := c.PostForm("amount")
//     amount, err := strconv.ParseFloat(amountStr, 64)
//     if err != nil {
//         c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid amount"})
//         return
//     }

    // if publisher.Credit >= int(amount) {
    //     publisher.Credit -= int(amount)
    //     DB.Save(&publisher)
    // } else {
    //     c.JSON(http.StatusBadRequest, gin.H{"error": "Insufficient balance"})
    //     return
    // }

//     c.Redirect(http.StatusSeeOther, "/publisher/"+id)
// }


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

// func (ctrl PublisherController) GetPublisherByID(c *gin.Context) {
// 	id, err := strconv.Atoi(c.Param("id"))
// 	if err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
// 		return
// 	}

// 	publisher, err := ctrl.Repo.FindByID(uint(id))
// 	if err != nil {
// 		c.JSON(http.StatusNotFound, gin.H{"error": "Publisher not found"})
// 		return
// 	}
// 	c.JSON(http.StatusOK, publisher)
// }

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

func (ctrl PublisherController) GetAllPublishers(c *gin.Context) {
	publishers, err := ctrl.Repo.FindAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, publishers)
}
