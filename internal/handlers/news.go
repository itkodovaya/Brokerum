package handlers

import (
	"net/http"
	"tenderhelp/internal/database"
	"tenderhelp/internal/models"

	"github.com/gin-gonic/gin"
)

func GetNews(c *gin.Context) {
	db := database.InitDB()
	var news []models.News

	// Получаем параметры фильтрации
	category := c.Query("category")
	search := c.Query("search")

	query := db.Order("date DESC")

	if category != "" && category != "all" {
		query = query.Where("category = ?", category)
	}

	if search != "" {
		query = query.Where("title LIKE ? OR content LIKE ?",
			"%"+search+"%", "%"+search+"%")
	}

	query.Find(&news)

	c.JSON(http.StatusOK, gin.H{
		"data":  news,
		"total": len(news),
	})
}

func CreateNews(c *gin.Context) {
	db := database.InitDB()
	var news models.News

	if err := c.ShouldBindJSON(&news); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	db.Create(&news)
	c.JSON(http.StatusCreated, news)
}
