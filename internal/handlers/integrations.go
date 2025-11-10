package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"tenderhelp/internal/adapters"
	"tenderhelp/internal/queue"
	"time"

	"github.com/gin-gonic/gin"
)

// AdapterManager глобальный менеджер адаптеров
var adapterManager *adapters.AdapterManager

// QueueManager глобальный менеджер очередей
var queueManager *queue.QueueManager

// InitIntegrations инициализирует систему интеграций
func InitIntegrations() {
	// Инициализация менеджера адаптеров
	adapterManager = adapters.NewAdapterManager()

	// Инициализация менеджера очередей
	queueManager = queue.NewQueueManager(3) // 3 воркера

	// Создание обработчика задач
	processor := queue.NewBankTaskProcessor(adapterManager)

	// Запуск воркеров
	queueManager.StartWorkers(processor)
}

// SendApplicationToBanks отправляет заявку в банки
func SendApplicationToBanks(c *gin.Context) {
	applicationID := c.Param("id")

	// Получение заявки из базы данных
	var application Application
	if err := db.First(&application, applicationID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Заявка не найдена"})
		return
	}

	// Подготовка данных для отправки
	applicationData := adapters.ApplicationData{
		ID:          fmt.Sprintf("%d", application.ID),
		Type:        application.Type,
		Amount:      application.Amount,
		ClientData:  application.PersonalData, // Упрощение для примера
		ScoringData: json.RawMessage(`{"score": 85.5, "risk_class": "A"}`),
		CreatedAt:   application.CreatedAt,
	}

	// Получение списка банков для отправки
	bankIDs := c.QueryArray("banks")
	if len(bankIDs) == 0 {
		// Отправка во все активные банки
		activeAdapters := adapterManager.GetActiveAdapters()
		for bankID := range activeAdapters {
			bankIDs = append(bankIDs, bankID)
		}
	}

	// Отправка в указанные банки
	responses, err := adapterManager.SendToSpecificBanks(applicationData, bankIDs)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка отправки в банки: " + err.Error()})
		return
	}

	// Обновление статуса заявки
	application.Status = "sent_to_banks"
	application.UpdatedAt = time.Now()
	db.Save(&application)

	// Добавление записи в историю
	statusHistory := StatusHistory{
		ApplicationID: application.ID,
		Status:        "sent_to_banks",
		Timestamp:     time.Now(),
		Comment:       "Заявка отправлена в банки",
	}
	db.Create(&statusHistory)

	c.JSON(http.StatusOK, gin.H{
		"message":   "Заявка отправлена в банки",
		"responses": responses,
	})
}

// GetBankSummary возвращает сводку по банкам
func GetBankSummary(c *gin.Context) {
	summary := adapterManager.GetBankSummary()

	c.JSON(http.StatusOK, gin.H{
		"banks": summary,
		"total": len(summary),
	})
}

// GetSupportedBanks возвращает банки, поддерживающие определенный тип заявки
func GetSupportedBanks(c *gin.Context) {
	applicationType := c.Query("type")
	if applicationType == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Тип заявки не указан"})
		return
	}

	supportedBanks := adapterManager.GetSupportedBanksForType(applicationType)

	c.JSON(http.StatusOK, gin.H{
		"type":  applicationType,
		"banks": supportedBanks,
	})
}

// CheckBankAvailability проверяет доступность банка
func CheckBankAvailability(c *gin.Context) {
	bankID := c.Param("bankId")

	available := adapterManager.CheckBankAvailability(bankID)

	c.JSON(http.StatusOK, gin.H{
		"bank_id":   bankID,
		"available": available,
	})
}

// GetApplicationStatusFromBank получает статус заявки из банка
func GetApplicationStatusFromBank(c *gin.Context) {
	applicationID := c.Param("id")
	bankID := c.Param("bankId")

	// Получение адаптера банка
	adapter, err := adapterManager.GetAdapter(bankID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Банк не найден"})
		return
	}

	// Получение статуса заявки
	status, err := adapter.GetApplicationStatus(applicationID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка получения статуса: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, status)
}

// GetQueueStats возвращает статистику очереди
func GetQueueStats(c *gin.Context) {
	stats := queueManager.GetQueueStats()

	c.JSON(http.StatusOK, stats)
}

// GetTaskStatus возвращает статус задачи
func GetTaskStatus(c *gin.Context) {
	taskID := c.Param("taskId")

	task, err := queueManager.GetTaskStatus(taskID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Задача не найдена"})
		return
	}

	c.JSON(http.StatusOK, task)
}

// CreateTask создает новую задачу
func CreateTask(c *gin.Context) {
	var request struct {
		Type          string          `json:"type" binding:"required"`
		ApplicationID string          `json:"application_id" binding:"required"`
		BankID        string          `json:"bank_id" binding:"required"`
		Data          json.RawMessage `json:"data" binding:"required"`
		Priority      int             `json:"priority"`
		MaxRetries    int             `json:"max_retries"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Создание задачи
	task := &queue.Task{
		ID:            "",
		Type:          request.Type,
		ApplicationID: request.ApplicationID,
		BankID:        request.BankID,
		Data:          request.Data,
		Priority:      request.Priority,
		MaxRetries:    request.MaxRetries,
		CreatedAt:     time.Now(),
	}

	// Добавление задачи в очередь
	if err := queueManager.AddTask(task); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка создания задачи: " + err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Задача создана",
		"task_id": task.ID,
	})
}

// GetBankResponses возвращает ответы от банков для заявки
func GetBankResponses(c *gin.Context) {
	applicationID := c.Param("id")

	// Получение заявки
	var application Application
	if err := db.First(&application, applicationID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Заявка не найдена"})
		return
	}

	// Получение истории статусов
	var statusHistory []StatusHistory
	db.Where("application_id = ?", applicationID).Order("timestamp DESC").Find(&statusHistory)

	// Формирование ответов от банков
	responses := make([]map[string]interface{}, 0)
	for _, status := range statusHistory {
		if status.Status == "sent_to_banks" || status.Status == "bank_response" {
			responses = append(responses, map[string]interface{}{
				"status":    status.Status,
				"comment":   status.Comment,
				"timestamp": status.Timestamp,
			})
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"application_id": applicationID,
		"responses":      responses,
		"total":          len(responses),
	})
}
