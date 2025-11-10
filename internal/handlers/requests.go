package handlers

import (
	"net/http"
	"strconv"
	"tenderhelp/internal/database"
	"tenderhelp/internal/models"

	"github.com/gin-gonic/gin"
)

func GetRequests(c *gin.Context) {
	db := database.InitDB()
	var requests []models.Request

	// Получаем параметры фильтрации
	status := c.Query("status")
	clientID := c.Query("client_id")

	query := db.Preload("Client").Preload("Agent")

	if status != "" {
		query = query.Where("status = ?", status)
	}
	if clientID != "" {
		query = query.Where("client_id = ?", clientID)
	}

	query.Find(&requests)

	c.JSON(http.StatusOK, gin.H{
		"data":  requests,
		"total": len(requests),
	})
}

func CreateRequest(c *gin.Context) {
	db := database.InitDB()
	var request models.Request

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	db.Create(&request)
	c.JSON(http.StatusCreated, request)
}

func UpdateRequest(c *gin.Context) {
	db := database.InitDB()
	id, _ := strconv.Atoi(c.Param("id"))

	var request models.Request
	if err := db.First(&request, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Заявка не найдена"})
		return
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	db.Save(&request)
	c.JSON(http.StatusOK, request)
}

func DeleteRequest(c *gin.Context) {
	db := database.InitDB()
	id, _ := strconv.Atoi(c.Param("id"))

	db.Delete(&models.Request{}, id)
	c.JSON(http.StatusOK, gin.H{"message": "Заявка удалена"})
}
