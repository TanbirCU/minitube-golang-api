package controllers

import (
	"net/http"

	"minitube-api/config"
	"minitube-api/models"

	"github.com/gin-gonic/gin"
)

func GetCategories(c *gin.Context) {

	var categories []models.Category

	config.DB.Order("name ASC").Find(&categories)

	c.JSON(http.StatusOK, gin.H{
		"data": categories,
	})
}

func CreateCategory(c *gin.Context) {

	var request struct {
		Name string `json:"name"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid request",
		})
		return
	}

	category := models.Category{
		Name: request.Name,
	}

	if err := config.DB.Create(&category).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Category creation failed",
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"data": category,
	})
}