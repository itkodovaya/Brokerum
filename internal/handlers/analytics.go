package handlers

import (
	"net/http"
	"tenderhelp/internal/database"
	"tenderhelp/internal/models"
	"time"

	"github.com/gin-gonic/gin"
)

func GetAnalytics(c *gin.Context) {
	db := database.InitDB()

	// Получаем статистику
	var totalRequests int64
	var newClients int64
	var inProgressRequests int64
	var approvedRequests int64

	db.Model(&models.Request{}).Count(&totalRequests)
	db.Model(&models.Client{}).Where("created_at > ?", time.Now().AddDate(0, 0, -7)).Count(&newClients)
	db.Model(&models.Request{}).Where("status = ?", "В работе").Count(&inProgressRequests)
	db.Model(&models.Request{}).Where("status = ?", "Одобрено").Count(&approvedRequests)

	// Расчет процентов
	var approvalRate float64
	var issuanceRate float64

	if totalRequests > 0 {
		approvalRate = float64(approvedRequests) / float64(totalRequests) * 100
		issuanceRate = float64(approvedRequests) / float64(totalRequests) * 100
	}

	// Данные для графика по месяцам
	var monthlyData []gin.H
	now := time.Now()

	for i := 0; i < 12; i++ {
		month := now.AddDate(0, -i, 0)
		startOfMonth := time.Date(month.Year(), month.Month(), 1, 0, 0, 0, 0, month.Location())
		endOfMonth := startOfMonth.AddDate(0, 1, 0)

		var count int64
		db.Model(&models.Request{}).Where("created_at BETWEEN ? AND ?", startOfMonth, endOfMonth).Count(&count)

		monthlyData = append([]gin.H{
			gin.H{
				"month": month.Format("2006-01"),
				"count": int(count),
			},
		}, monthlyData...)
	}

	c.JSON(http.StatusOK, gin.H{
		"metrics": gin.H{
			"total_requests": totalRequests,
			"new_clients":    newClients,
			"earned_amount":  0, // TODO: рассчитать сумму
			"in_progress":    inProgressRequests,
			"approval_rate":  approvalRate,
			"issuance_rate":  issuanceRate,
		},
		"monthly_data": monthlyData,
	})
}
