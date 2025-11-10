package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"tenderhelp/internal/database"
	"tenderhelp/internal/scoring"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// Application представляет заявку
type Application struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	ClientID  uint      `json:"client_id"`
	Type      string    `json:"type"` // credit, guarantee
	Amount    float64   `json:"amount"`
	Status    string    `json:"status"`
	Bank      string    `json:"bank"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	// Данные анкеты
	PersonalData     json.RawMessage `json:"personal_data" gorm:"type:jsonb"`
	ContactData      json.RawMessage `json:"contact_data" gorm:"type:jsonb"`
	ProfessionalData json.RawMessage `json:"professional_data" gorm:"type:jsonb"`
	FinancialData    json.RawMessage `json:"financial_data" gorm:"type:jsonb"`
	FamilyData       json.RawMessage `json:"family_data" gorm:"type:jsonb"`
	AdditionalData   json.RawMessage `json:"additional_data" gorm:"type:jsonb"`

	// История статусов
	StatusHistory []StatusHistory `json:"status_history" gorm:"foreignKey:ApplicationID"`
}

// StatusHistory представляет историю статусов
type StatusHistory struct {
	ID            uint      `json:"id" gorm:"primaryKey"`
	ApplicationID uint      `json:"application_id"`
	Status        string    `json:"status"`
	Timestamp     time.Time `json:"timestamp"`
	Comment       string    `json:"comment"`
}

// CreateApplicationRequest запрос на создание заявки
type CreateApplicationRequest struct {
	Type             string          `json:"type" binding:"required"`
	Amount           float64         `json:"amount"`
	PersonalData     json.RawMessage `json:"personalData" binding:"required"`
	ContactData      json.RawMessage `json:"contactData" binding:"required"`
	ProfessionalData json.RawMessage `json:"professionalData" binding:"required"`
	FinancialData    json.RawMessage `json:"financialData" binding:"required"`
	FamilyData       json.RawMessage `json:"familyData" binding:"required"`
	AdditionalData   json.RawMessage `json:"additionalData" binding:"required"`
}

// UpdateApplicationRequest запрос на обновление заявки
type UpdateApplicationRequest struct {
	Step string          `json:"step" binding:"required"`
	Data json.RawMessage `json:"data" binding:"required"`
}

// GetApplicationsResponse ответ со списком заявок
type GetApplicationsResponse struct {
	Applications []Application `json:"applications"`
	Total        int64         `json:"total"`
	Page         int           `json:"page"`
	Limit        int           `json:"limit"`
}

// CreateApplication создает новую заявку
func CreateApplication(c *gin.Context) {
	var req CreateApplicationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Валидация данных анкеты
	if err := validateApplicationData(req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Создание заявки
	application := Application{
		ClientID:         1, // TODO: получить из JWT токена
		Type:             req.Type,
		Amount:           req.Amount,
		Status:           "draft",
		PersonalData:     req.PersonalData,
		ContactData:      req.ContactData,
		ProfessionalData: req.ProfessionalData,
		FinancialData:    req.FinancialData,
		FamilyData:       req.FamilyData,
		AdditionalData:   req.AdditionalData,
		CreatedAt:        time.Now(),
		UpdatedAt:        time.Now(),
	}

	// Сохранение в базу данных
	if err := db.Create(&application).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка создания заявки"})
		return
	}

	// Добавление записи в историю статусов
	statusHistory := StatusHistory{
		ApplicationID: application.ID,
		Status:        "draft",
		Timestamp:     time.Now(),
		Comment:       "Заявка создана",
	}
	db.Create(&statusHistory)

	c.JSON(http.StatusCreated, gin.H{
		"message":     "Заявка успешно создана",
		"application": application,
	})
}

// GetApplications возвращает список заявок с фильтрацией
func GetApplications(c *gin.Context) {
	// Параметры пагинации
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	offset := (page - 1) * limit

	// Фильтры
	status := c.Query("status")
	dateFrom := c.Query("dateFrom")
	dateTo := c.Query("dateTo")

	// Построение запроса
	query := db.Model(&Application{})

	// Применение фильтров
	if status != "" && status != "all" {
		query = query.Where("status = ?", status)
	}
	if dateFrom != "" {
		query = query.Where("created_at >= ?", dateFrom)
	}
	if dateTo != "" {
		query = query.Where("created_at <= ?", dateTo)
	}

	// Подсчет общего количества
	var total int64
	query.Count(&total)

	// Получение заявок с пагинацией
	var applications []Application
	if err := query.
		Preload("StatusHistory").
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&applications).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка получения заявок"})
		return
	}

	response := GetApplicationsResponse{
		Applications: applications,
		Total:        total,
		Page:         page,
		Limit:        limit,
	}

	c.JSON(http.StatusOK, response)
}

// GetApplication возвращает заявку по ID
func GetApplication(c *gin.Context) {
	id := c.Param("id")

	var application Application
	if err := db.Preload("StatusHistory").First(&application, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Заявка не найдена"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка получения заявки"})
		return
	}

	c.JSON(http.StatusOK, application)
}

// UpdateApplication обновляет заявку
func UpdateApplication(c *gin.Context) {
	id := c.Param("id")

	var req UpdateApplicationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Получение заявки
	var application Application
	if err := db.First(&application, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Заявка не найдена"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка получения заявки"})
		return
	}

	// Валидация данных шага
	if err := validateStepData(req.Step, req.Data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Обновление данных в зависимости от шага
	switch req.Step {
	case "personal":
		application.PersonalData = req.Data
	case "contact":
		application.ContactData = req.Data
	case "professional":
		application.ProfessionalData = req.Data
	case "financial":
		application.FinancialData = req.Data
	case "family":
		application.FamilyData = req.Data
	case "additional":
		application.AdditionalData = req.Data
	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неизвестный шаг"})
		return
	}

	application.UpdatedAt = time.Now()

	// Сохранение изменений
	if err := db.Save(&application).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка обновления заявки"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":     "Заявка успешно обновлена",
		"application": application,
	})
}

// SubmitApplication отправляет заявку на рассмотрение
func SubmitApplication(c *gin.Context) {
	id := c.Param("id")

	// Получение заявки
	var application Application
	if err := db.First(&application, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Заявка не найдена"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка получения заявки"})
		return
	}

	// Проверка, что заявка в статусе "draft"
	if application.Status != "draft" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Заявка уже отправлена или не может быть изменена"})
		return
	}

	// Валидация всех данных
	if err := validateApplicationData(CreateApplicationRequest{
		PersonalData:     application.PersonalData,
		ContactData:      application.ContactData,
		ProfessionalData: application.ProfessionalData,
		FinancialData:    application.FinancialData,
		FamilyData:       application.FamilyData,
		AdditionalData:   application.AdditionalData,
	}); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Не все обязательные поля заполнены: " + err.Error()})
		return
	}

	// Обновление статуса
	application.Status = "submitted"
	application.UpdatedAt = time.Now()

	if err := db.Save(&application).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка обновления заявки"})
		return
	}

	// Добавление записи в историю статусов
	statusHistory := StatusHistory{
		ApplicationID: application.ID,
		Status:        "submitted",
		Timestamp:     time.Now(),
		Comment:       "Заявка отправлена на рассмотрение",
	}
	db.Create(&statusHistory)

	// TODO: Запуск скоринга и отправка в банки
	go processApplication(application.ID)

	c.JSON(http.StatusOK, gin.H{
		"message":     "Заявка успешно отправлена",
		"application": application,
	})
}

// validateApplicationData валидирует данные анкеты
func validateApplicationData(req CreateApplicationRequest) error {
	// Проверка обязательных полей
	if len(req.PersonalData) == 0 {
		return fmt.Errorf("личные данные не заполнены")
	}
	if len(req.ContactData) == 0 {
		return fmt.Errorf("контактные данные не заполнены")
	}
	if len(req.ProfessionalData) == 0 {
		return fmt.Errorf("профессиональные данные не заполнены")
	}
	if len(req.FinancialData) == 0 {
		return fmt.Errorf("финансовые данные не заполнены")
	}
	if len(req.FamilyData) == 0 {
		return fmt.Errorf("семейные данные не заполнены")
	}
	if len(req.AdditionalData) == 0 {
		return fmt.Errorf("дополнительные данные не заполнены")
	}

	// TODO: Добавить детальную валидацию каждого шага
	return nil
}

// validateStepData валидирует данные конкретного шага
func validateStepData(step string, data json.RawMessage) error {
	if len(data) == 0 {
		return fmt.Errorf("данные шага не заполнены")
	}

	// TODO: Добавить детальную валидацию для каждого шага
	return nil
}

// processApplication обрабатывает заявку (скоринг, отправка в банки)
func processApplication(applicationID uint) {
	// TODO: Реализовать скоринг
	// TODO: Отправить в банки
	// TODO: Обновить статус заявки
	fmt.Printf("Обработка заявки %d\n", applicationID)
}

// SetScoringEngine устанавливает движок скоринга
func SetScoringEngine(engine *scoring.ScoringEngine) {
	scoringEngine = engine
}

// Глобальная переменная для движка скоринга
var scoringEngine *scoring.ScoringEngine

// Глобальная переменная для базы данных
var db *gorm.DB

// InitDB инициализирует подключение к базе данных
func InitDB() {
	db = database.InitDB()
}
