package handlers

import (
	"net/http"
	"strings"
	"tenderhelp/internal/database"
	"tenderhelp/internal/models"

	"github.com/gin-gonic/gin"
)

func GetHelpCategories(c *gin.Context) {
	db := database.InitDB()
	var categories []models.HelpCategory

	db.Find(&categories)

	c.JSON(http.StatusOK, gin.H{
		"data": categories,
	})
}

func GetHelpArticles(c *gin.Context) {
	db := database.InitDB()
	var articles []models.HelpArticle

	categoryID := c.Query("category_id")

	query := db.Preload("Category")

	if categoryID != "" {
		query = query.Where("category_id = ?", categoryID)
	}

	query.Find(&articles)

	c.JSON(http.StatusOK, gin.H{
		"data": articles,
	})
}

func SearchHelp(c *gin.Context) {
	db := database.InitDB()
	searchQuery := c.Query("q")

	if searchQuery == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Поисковый запрос не может быть пустым"})
		return
	}

	var articles []models.HelpArticle
	searchTerm := "%" + strings.ToLower(searchQuery) + "%"

	db.Preload("Category").Where(
		"LOWER(title) LIKE ? OR LOWER(content) LIKE ?",
		searchTerm, searchTerm,
	).Find(&articles)

	c.JSON(http.StatusOK, gin.H{
		"data":  articles,
		"query": searchQuery,
	})
}
