package handlers

import (
	"net/http"
	"strconv"
	"tenderhelp/internal/database"
	"tenderhelp/internal/models"

	"github.com/gin-gonic/gin"
)

func GetClients(c *gin.Context) {
	db := database.InitDB()
	var clients []models.Client

	// Получаем параметры поиска
	search := c.Query("search")

	query := db

	if search != "" {
		query = query.Where("name LIKE ? OR inn LIKE ? OR email LIKE ?",
			"%"+search+"%", "%"+search+"%", "%"+search+"%")
	}

	query.Find(&clients)

	c.JSON(http.StatusOK, gin.H{
		"data":  clients,
		"total": len(clients),
	})
}

func CreateClient(c *gin.Context) {
	db := database.InitDB()
	var client models.Client

	if err := c.ShouldBindJSON(&client); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	db.Create(&client)
	c.JSON(http.StatusCreated, client)
}

func UpdateClient(c *gin.Context) {
	db := database.InitDB()
	id, _ := strconv.Atoi(c.Param("id"))

	var client models.Client
	if err := db.First(&client, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Клиент не найден"})
		return
	}

	if err := c.ShouldBindJSON(&client); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	db.Save(&client)
	c.JSON(http.StatusOK, client)
}

func DeleteClient(c *gin.Context) {
	db := database.InitDB()
	id, _ := strconv.Atoi(c.Param("id"))

	db.Delete(&models.Client{}, id)
	c.JSON(http.StatusOK, gin.H{"message": "Клиент удален"})
}
