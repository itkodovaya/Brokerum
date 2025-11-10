package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

// POSApplication заявка на кредит для ПОС
type POSApplication struct {
	ID                 uint      `json:"id" gorm:"primaryKey"`
	ApplicationID      uint      `json:"application_id"`
	MerchantID         string    `json:"merchant_id"`
	TerminalID         string    `json:"terminal_id"`
	MerchantName       string    `json:"merchant_name"`
	MerchantCategory   string    `json:"merchant_category"`
	MonthlyTurnover    float64   `json:"monthly_turnover"`
	AverageTransaction float64   `json:"average_transaction"`
	TransactionCount   int       `json:"transaction_count"`
	CreditLimit        float64   `json:"credit_limit"`
	InterestRate       float64   `json:"interest_rate"`
	Term               int       `json:"term"`
	IsActive           bool      `json:"is_active"`
	CreatedAt          time.Time `json:"created_at"`
	UpdatedAt          time.Time `json:"updated_at"`
}

// GuaranteeApplication заявка на банковскую гарантию
type GuaranteeApplication struct {
	ID                 uint      `json:"id" gorm:"primaryKey"`
	ApplicationID      uint      `json:"application_id"`
	GuaranteeAmount    float64   `json:"guarantee_amount"`
	GuaranteePurpose   string    `json:"guarantee_purpose"`
	GuaranteePeriod    int       `json:"guarantee_period"` // в днях
	BeneficiaryName    string    `json:"beneficiary_name"`
	BeneficiaryINN     string    `json:"beneficiary_inn"`
	BeneficiaryAddress string    `json:"beneficiary_address"`
	ContractNumber     string    `json:"contract_number"`
	ContractDate       time.Time `json:"contract_date"`
	ContractAmount     float64   `json:"contract_amount"`
	GuaranteeType      string    `json:"guarantee_type"` // performance, payment, advance
	IsActive           bool      `json:"is_active"`
	CreatedAt          time.Time `json:"created_at"`
	UpdatedAt          time.Time `json:"updated_at"`
}

// CreatePOSApplication создает заявку на кредит для ПОС
func CreatePOSApplication(c *gin.Context) {
	applicationIDStr := c.Param("id")
	applicationID, err := strconv.ParseUint(applicationIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Некорректный ID заявки"})
		return
	}

	var request struct {
		MerchantID         string  `json:"merchant_id" binding:"required"`
		TerminalID         string  `json:"terminal_id" binding:"required"`
		MerchantName       string  `json:"merchant_name" binding:"required"`
		MerchantCategory   string  `json:"merchant_category" binding:"required"`
		MonthlyTurnover    float64 `json:"monthly_turnover" binding:"required"`
		AverageTransaction float64 `json:"average_transaction" binding:"required"`
		TransactionCount   int     `json:"transaction_count" binding:"required"`
		CreditLimit        float64 `json:"credit_limit" binding:"required"`
		InterestRate       float64 `json:"interest_rate" binding:"required"`
		Term               int     `json:"term" binding:"required"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Создание заявки на ПОС
	posApplication := POSApplication{
		ApplicationID:      uint(applicationID),
		MerchantID:         request.MerchantID,
		TerminalID:         request.TerminalID,
		MerchantName:       request.MerchantName,
		MerchantCategory:   request.MerchantCategory,
		MonthlyTurnover:    request.MonthlyTurnover,
		AverageTransaction: request.AverageTransaction,
		TransactionCount:   request.TransactionCount,
		CreditLimit:        request.CreditLimit,
		InterestRate:       request.InterestRate,
		Term:               request.Term,
		IsActive:           true,
		CreatedAt:          time.Now(),
		UpdatedAt:          time.Now(),
	}

	// Сохранение в базу данных
	if err := db.Create(&posApplication).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка создания заявки на ПОС"})
		return
	}

	// Обновление основной заявки
	var application Application
	if err := db.First(&application, applicationID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Заявка не найдена"})
		return
	}

	application.Type = "pos_credit"
	application.UpdatedAt = time.Now()
	db.Save(&application)

	c.JSON(http.StatusCreated, gin.H{
		"message":         "Заявка на ПОС создана",
		"pos_application": posApplication,
	})
}

// CreateGuaranteeApplication создает заявку на банковскую гарантию
func CreateGuaranteeApplication(c *gin.Context) {
	applicationIDStr := c.Param("id")
	applicationID, err := strconv.ParseUint(applicationIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Некорректный ID заявки"})
		return
	}

	var request struct {
		GuaranteeAmount    float64   `json:"guarantee_amount" binding:"required"`
		GuaranteePurpose   string    `json:"guarantee_purpose" binding:"required"`
		GuaranteePeriod    int       `json:"guarantee_period" binding:"required"`
		BeneficiaryName    string    `json:"beneficiary_name" binding:"required"`
		BeneficiaryINN     string    `json:"beneficiary_inn" binding:"required"`
		BeneficiaryAddress string    `json:"beneficiary_address" binding:"required"`
		ContractNumber     string    `json:"contract_number" binding:"required"`
		ContractDate       time.Time `json:"contract_date" binding:"required"`
		ContractAmount     float64   `json:"contract_amount" binding:"required"`
		GuaranteeType      string    `json:"guarantee_type" binding:"required"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Создание заявки на банковскую гарантию
	guaranteeApplication := GuaranteeApplication{
		ApplicationID:      uint(applicationID),
		GuaranteeAmount:    request.GuaranteeAmount,
		GuaranteePurpose:   request.GuaranteePurpose,
		GuaranteePeriod:    request.GuaranteePeriod,
		BeneficiaryName:    request.BeneficiaryName,
		BeneficiaryINN:     request.BeneficiaryINN,
		BeneficiaryAddress: request.BeneficiaryAddress,
		ContractNumber:     request.ContractNumber,
		ContractDate:       request.ContractDate,
		ContractAmount:     request.ContractAmount,
		GuaranteeType:      request.GuaranteeType,
		IsActive:           true,
		CreatedAt:          time.Now(),
		UpdatedAt:          time.Now(),
	}

	// Сохранение в базу данных
	if err := db.Create(&guaranteeApplication).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка создания заявки на банковскую гарантию"})
		return
	}

	// Обновление основной заявки
	var application Application
	if err := db.First(&application, applicationID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Заявка не найдена"})
		return
	}

	application.Type = "guarantee"
	application.UpdatedAt = time.Now()
	db.Save(&application)

	c.JSON(http.StatusCreated, gin.H{
		"message":               "Заявка на банковскую гарантию создана",
		"guarantee_application": guaranteeApplication,
	})
}

// GetPOSApplication получает заявку на ПОС
func GetPOSApplication(c *gin.Context) {
	applicationID := c.Param("id")

	var posApplication POSApplication
	if err := db.Where("application_id = ?", applicationID).First(&posApplication).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Заявка на ПОС не найдена"})
		return
	}

	c.JSON(http.StatusOK, posApplication)
}

// GetGuaranteeApplication получает заявку на банковскую гарантию
func GetGuaranteeApplication(c *gin.Context) {
	applicationID := c.Param("id")

	var guaranteeApplication GuaranteeApplication
	if err := db.Where("application_id = ?", applicationID).First(&guaranteeApplication).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Заявка на банковскую гарантию не найдена"})
		return
	}

	c.JSON(http.StatusOK, guaranteeApplication)
}

// GeneratePOSDocument генерирует документ для ПОС
func GeneratePOSDocument(c *gin.Context) {
	applicationID := c.Param("id")

	// Получение заявки на ПОС
	var posApplication POSApplication
	if err := db.Where("application_id = ?", applicationID).First(&posApplication).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Заявка на ПОС не найдена"})
		return
	}

	// Получение основной заявки
	var application Application
	if err := db.First(&application, applicationID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Заявка не найдена"})
		return
	}

	// Извлечение данных клиента
	var clientData map[string]interface{}
	json.Unmarshal(application.PersonalData, &clientData)

	// Генерация документа
	document := map[string]interface{}{
		"document_type":  "POS_CREDIT_APPLICATION",
		"application_id": applicationID,
		"generated_at":   time.Now(),
		"client_info": map[string]interface{}{
			"name":  clientData["firstName"].(string) + " " + clientData["lastName"].(string),
			"phone": clientData["primaryPhone"],
			"email": clientData["email"],
		},
		"pos_info": map[string]interface{}{
			"merchant_id":         posApplication.MerchantID,
			"terminal_id":         posApplication.TerminalID,
			"merchant_name":       posApplication.MerchantName,
			"category":            posApplication.MerchantCategory,
			"monthly_turnover":    posApplication.MonthlyTurnover,
			"average_transaction": posApplication.AverageTransaction,
		},
		"credit_terms": map[string]interface{}{
			"credit_limit":  posApplication.CreditLimit,
			"interest_rate": posApplication.InterestRate,
			"term":          posApplication.Term,
		},
	}

	c.JSON(http.StatusOK, gin.H{
		"message":  "Документ для ПОС сгенерирован",
		"document": document,
	})
}

// GenerateGuaranteeDocument генерирует документ для банковской гарантии
func GenerateGuaranteeDocument(c *gin.Context) {
	applicationID := c.Param("id")

	// Получение заявки на банковскую гарантию
	var guaranteeApplication GuaranteeApplication
	if err := db.Where("application_id = ?", applicationID).First(&guaranteeApplication).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Заявка на банковскую гарантию не найдена"})
		return
	}

	// Получение основной заявки
	var application Application
	if err := db.First(&application, applicationID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Заявка не найдена"})
		return
	}

	// Извлечение данных клиента
	var clientData map[string]interface{}
	json.Unmarshal(application.PersonalData, &clientData)

	// Генерация документа
	document := map[string]interface{}{
		"document_type":  "BANK_GUARANTEE_APPLICATION",
		"application_id": applicationID,
		"generated_at":   time.Now(),
		"client_info": map[string]interface{}{
			"name":  clientData["firstName"].(string) + " " + clientData["lastName"].(string),
			"phone": clientData["primaryPhone"],
			"email": clientData["email"],
		},
		"guarantee_info": map[string]interface{}{
			"amount":  guaranteeApplication.GuaranteeAmount,
			"purpose": guaranteeApplication.GuaranteePurpose,
			"period":  guaranteeApplication.GuaranteePeriod,
			"type":    guaranteeApplication.GuaranteeType,
		},
		"beneficiary_info": map[string]interface{}{
			"name":    guaranteeApplication.BeneficiaryName,
			"inn":     guaranteeApplication.BeneficiaryINN,
			"address": guaranteeApplication.BeneficiaryAddress,
		},
		"contract_info": map[string]interface{}{
			"number": guaranteeApplication.ContractNumber,
			"date":   guaranteeApplication.ContractDate,
			"amount": guaranteeApplication.ContractAmount,
		},
	}

	c.JSON(http.StatusOK, gin.H{
		"message":  "Документ для банковской гарантии сгенерирован",
		"document": document,
	})
}
