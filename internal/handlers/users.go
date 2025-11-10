package handlers

import (
	"math/rand"
	"net/http"
	"strings"
	"tenderhelp/internal/models"
	"time"

	"github.com/gin-gonic/gin"
)

type registerAgentRequest struct {
	Name  string `json:"name"`
	Phone string `json:"phone"`
	Email string `json:"email"`
}

// RegisterAgent обрабатывает регистрацию агента и создает пользователя с ролью agent.
func RegisterAgent(c *gin.Context) {
	var req registerAgentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Некорректные данные"})
		return
	}

	req.Email = strings.TrimSpace(strings.ToLower(req.Email))
	req.Name = strings.TrimSpace(req.Name)
	req.Phone = strings.TrimSpace(req.Phone)

	if req.Name == "" || req.Phone == "" || req.Email == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ФИО, телефон и почта обязательны"})
		return
	}

	// Генерация временного пароля
	tempPassword := generateTempPassword(10)

	user := models.User{
		Email:    req.Email,
		Password: tempPassword,
		Name:     req.Name,
		Role:     "agent",
		// Phone сохранится в поле модели (добавлено в модели)
		Phone: req.Phone,
	}

	if err := db.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Не удалось создать пользователя"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"id":                user.ID,
		"email":             user.Email,
		"name":              user.Name,
		"phone":             user.Phone,
		"role":              user.Role,
		"temporaryPassword": tempPassword,
	})
}

func generateTempPassword(length int) string {
	const letters = "abcdefghjkmnpqrstuvwxyzABCDEFGHJKMNPQRSTUVWXYZ23456789"
	rand.Seed(time.Now().UnixNano())
	b := make([]byte, length)
	for i := 0; i < length; i++ {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

type registerClientRequest struct {
	Name          string `json:"name"`
	INN           string `json:"inn"`
	OGRN          string `json:"ogrn"`
	AccountNumber string `json:"account_number"`
	BIK           string `json:"bik"`
	CorrAccount   string `json:"corr_account"`
	BankName      string `json:"bank_name"`
	Address       string `json:"address"`
	Email         string `json:"email"`
	Phone         string `json:"phone"`
	City          string `json:"city"`
}

// RegisterClient создает клиента по предоставленным реквизитам
func RegisterClient(c *gin.Context) {
	var req registerClientRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Некорректные данные"})
		return
	}

	// Простейшая валидация обязательных полей
	if strings.TrimSpace(req.Name) == "" || strings.TrimSpace(req.INN) == "" || strings.TrimSpace(req.OGRN) == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Название, ИНН и ОГРН обязательны"})
		return
	}

	client := models.Client{
		Name:          strings.TrimSpace(req.Name),
		INN:           strings.TrimSpace(req.INN),
		OGRN:          strings.TrimSpace(req.OGRN),
		Phone:         strings.TrimSpace(req.Phone),
		Email:         strings.TrimSpace(req.Email),
		City:          strings.TrimSpace(req.City),
		Address:       strings.TrimSpace(req.Address),
		BankName:      strings.TrimSpace(req.BankName),
		AccountNumber: strings.TrimSpace(req.AccountNumber),
		BIK:           strings.TrimSpace(req.BIK),
		CorrAccount:   strings.TrimSpace(req.CorrAccount),
	}

	if err := db.Create(&client).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Не удалось создать клиента"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"id": client.ID, "name": client.Name})
}
