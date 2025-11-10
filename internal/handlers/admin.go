package handlers

import (
	"io"
	"net/http"
	"tenderhelp/internal/database"
	"tenderhelp/internal/models"

	"github.com/gin-gonic/gin"
)

// Bank represents a bank offer
type Bank struct {
	ID      uint   `json:"id"`
	Name    string `json:"name"`
	Product string `json:"product"`
	Amount  string `json:"amount"`
	Rate    string `json:"rate"`
	Term    string `json:"term"`
	Income  string `json:"income"`
	HasAPI  bool   `json:"has_api"`
	Logo    string `json:"logo"`
}

// GetBanks returns all banks
func GetBanks(c *gin.Context) {
	// For now, return mock data
	mockBanks := []Bank{
		{
			ID:      1,
			Name:    "МСП Банк",
			Product: "«Экспресс поддержка»",
			Amount:  "от 50 тыс. до 8 млн.",
			Rate:    "от 21.5%",
			Term:    "от 3 месяцев до 5 лет",
			Income:  "0.6% от суммы кредита",
			HasAPI:  true,
			Logo:    "🏦",
		},
		{
			ID:      2,
			Name:    "Альфа-Банк",
			Product: "«Бизнес кредит»",
			Amount:  "от 300 тыс. до 150 млн.",
			Rate:    "от 24.5%",
			Term:    "от 1 месяца до 6 лет",
			Income:  "0.1% - 0.6% от суммы кредита",
			HasAPI:  false,
			Logo:    "🏛️",
		},
		{
			ID:      3,
			Name:    "Просто Банк",
			Product: "«Кредитная линия»",
			Amount:  "от 100 тыс. до 10 млн.",
			Rate:    "от 17.9%",
			Term:    "от 6 месяцев до 2 лет",
			Income:  "0.9% от суммы кредита",
			HasAPI:  false,
			Logo:    "🏢",
		},
		{
			ID:      4,
			Name:    "ПСБ банк - КИК",
			Product: "«Без бумаг Контрактный»",
			Amount:  "от 50 тыс. до 100 млн.",
			Rate:    "от 21.1%",
			Term:    "от 3 месяцев до 2 лет",
			Income:  "0.6% от суммы кредита",
			HasAPI:  true,
			Logo:    "🏪",
		},
		{
			ID:      5,
			Name:    "Совком банк",
			Product: "«Оборотный кредит»",
			Amount:  "от 300 тыс. до 10 млн.",
			Rate:    "от 29%",
			Term:    "от 6 месяцев до 2 лет",
			Income:  "0.6% от суммы кредита",
			HasAPI:  false,
			Logo:    "🏬",
		},
	}

	c.JSON(http.StatusOK, mockBanks)
}

// UploadBanks handles Excel file upload for banks
func UploadBanks(c *gin.Context) {
	file, header, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No file uploaded"})
		return
	}
	defer file.Close()

	// Read file content
	fileContent, err := io.ReadAll(file)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read file"})
		return
	}

	// For now, just return success
	// In real implementation, you would parse the Excel file and save to database
	c.JSON(http.StatusOK, gin.H{
		"message":  "File uploaded successfully",
		"filename": header.Filename,
		"size":     len(fileContent),
	})
}

// GetAnalytics returns admin analytics
func GetAdminAnalytics(c *gin.Context) {
	db := database.InitDB()

	// Get statistics
	var totalRequests int64
	var approvedRequests int64
	var rejectedRequests int64
	var inProgressRequests int64

	db.Model(&models.Request{}).Count(&totalRequests)
	db.Model(&models.Request{}).Where("status = ?", "Одобрено").Count(&approvedRequests)
	db.Model(&models.Request{}).Where("status = ?", "Отклонено").Count(&rejectedRequests)
	db.Model(&models.Request{}).Where("status = ?", "В работе").Count(&inProgressRequests)

	analytics := gin.H{
		"total_requests": totalRequests,
		"approved":       approvedRequests,
		"rejected":       rejectedRequests,
		"in_progress":    inProgressRequests,
	}

	c.JSON(http.StatusOK, analytics)
}

// GetUsers returns all users
func GetUsers(c *gin.Context) {
	// Mock data for now
	users := []gin.H{
		{
			"id":     1,
			"name":   "Максим Искин Олегович",
			"email":  "agent@tenderhelp.ru",
			"role":   "Агент",
			"status": "Активен",
		},
	}
	c.JSON(http.StatusOK, users)
}

// CreateUser creates a new user
func CreateUser(c *gin.Context) {
	var user gin.H
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, user)
}

// UpdateUser updates a user
func UpdateUser(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "User updated"})
}

// DeleteUser deletes a user
func DeleteUser(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "User deleted"})
}

// CreateBank creates a new bank
func CreateBank(c *gin.Context) {
	var bank Bank
	if err := c.ShouldBindJSON(&bank); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, bank)
}

// UpdateBank updates a bank
func UpdateBank(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "Bank updated"})
}

// DeleteBank deletes a bank
func DeleteBank(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "Bank deleted"})
}

// GetAdminStats returns admin statistics
func GetAdminStats(c *gin.Context) {
	stats := gin.H{
		"total_users":     1,
		"total_banks":     5,
		"active_requests": 10,
		"system_uptime":   "2 дня",
	}
	c.JSON(http.StatusOK, stats)
}

// GetAdminRequests returns all requests for admin
func GetAdminRequests(c *gin.Context) {
	// Mock data
	requests := []gin.H{
		{
			"id":     1,
			"number": "6322",
			"client": "ООО \"КОНТУР\"",
			"status": "Черновик",
			"amount": "1,000,000 ₽",
		},
	}
	c.JSON(http.StatusOK, requests)
}

// UpdateRequestStatus updates request status
func UpdateRequestStatus(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "Request status updated"})
}

// GetAdminClients returns all clients for admin
func GetAdminClients(c *gin.Context) {
	// Mock data
	clients := []gin.H{
		{
			"id":     1,
			"name":   "ООО \"КОНТУР\"",
			"inn":    "1234567890",
			"email":  "info@kontur.ru",
			"status": "Активен",
		},
	}
	c.JSON(http.StatusOK, clients)
}
