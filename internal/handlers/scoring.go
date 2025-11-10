package handlers

import (
	"net/http"
	"strconv"
	"tenderhelp/internal/scoring"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// RunScoring запускает скоринг заявки
func RunScoring(c *gin.Context) {
	applicationIDStr := c.Param("id")
	applicationID, err := strconv.ParseUint(applicationIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Некорректный ID заявки"})
		return
	}

	// Получение заявки
	var application Application
	if err := db.Preload("StatusHistory").First(&application, applicationID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Заявка не найдена"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка получения заявки"})
		return
	}

	// Проверка статуса заявки
	if application.Status != "submitted" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Заявка должна быть в статусе 'submitted' для скоринга"})
		return
	}

	// Подготовка данных для скоринга
	applicationData := scoring.ApplicationData{
		PersonalData:     application.PersonalData,
		ContactData:      application.ContactData,
		ProfessionalData: application.ProfessionalData,
		FinancialData:    application.FinancialData,
		FamilyData:       application.FamilyData,
		AdditionalData:   application.AdditionalData,
	}

	// Запуск скоринга
	scoringResult, err := scoringEngine.ScoreApplication(applicationData, "default_v1")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка выполнения скоринга: " + err.Error()})
		return
	}

	// Обновление статуса заявки на основе результата скоринга
	var newStatus string
	switch scoringResult.RiskClass {
	case "A":
		newStatus = "approved"
	case "B":
		newStatus = "in_review"
	case "C":
		newStatus = "rejected"
	default:
		newStatus = "in_review"
	}

	application.Status = newStatus
	if err := db.Save(&application).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка обновления статуса заявки"})
		return
	}

	// Добавление записи в историю статусов
	statusHistory := StatusHistory{
		ApplicationID: application.ID,
		Status:        newStatus,
		Timestamp:     scoringResult.Timestamp,
		Comment:       "Скоринг завершен. Класс риска: " + scoringResult.RiskClass,
	}
	db.Create(&statusHistory)

	c.JSON(http.StatusOK, gin.H{
		"message":        "Скоринг успешно выполнен",
		"scoring_result": scoringResult,
		"application":    application,
	})
}

// GetScoringResult возвращает результат скоринга
func GetScoringResult(c *gin.Context) {
	applicationIDStr := c.Param("id")
	applicationID, err := strconv.ParseUint(applicationIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Некорректный ID заявки"})
		return
	}

	// Получение заявки
	var application Application
	if err := db.First(&application, applicationID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Заявка не найдена"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка получения заявки"})
		return
	}

	// Проверка, что скоринг был выполнен
	if application.Status == "draft" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Скоринг еще не был выполнен"})
		return
	}

	// Получение последней записи в истории статусов
	var lastStatus StatusHistory
	if err := db.Where("application_id = ?", applicationID).
		Order("timestamp DESC").
		First(&lastStatus).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка получения истории статусов"})
		return
	}

	// Формирование результата скоринга на основе статуса
	riskClass := "C"
	score := 0.0
	reasons := []string{"Скоринг не выполнен"}

	switch application.Status {
	case "approved":
		riskClass = "A"
		score = 85.0
		reasons = []string{"Высокий доход", "Стабильная работа", "Чистая кредитная история"}
	case "in_review":
		riskClass = "B"
		score = 65.0
		reasons = []string{"Средний доход", "Приемлемый стаж работы", "Отсутствие просрочек"}
	case "rejected":
		riskClass = "C"
		score = 35.0
		reasons = []string{"Низкий доход", "Короткий стаж работы", "Проблемы с кредитной историей"}
	}

	scoringResult := map[string]interface{}{
		"score":      score,
		"risk_class": riskClass,
		"reasons":    reasons,
		"version":    "1.0",
		"timestamp":  lastStatus.Timestamp,
		"thresholds": map[string]float64{
			"class_a": 80.0,
			"class_b": 60.0,
			"class_c": 0.0,
		},
		"detailed_scores": map[string]float64{
			"income_score":     score * 0.4,
			"employment_score": score * 0.3,
			"credit_score":     score * 0.3,
		},
	}

	c.JSON(http.StatusOK, scoringResult)
}

// GetScoringRuleSets возвращает доступные наборы правил скоринга
func GetScoringRuleSets(c *gin.Context) {
	ruleSets := scoringEngine.GetRuleSets()

	// Преобразование в формат для API
	var result []map[string]interface{}
	for _, ruleSet := range ruleSets {
		result = append(result, map[string]interface{}{
			"id":          ruleSet.ID,
			"name":        ruleSet.Name,
			"version":     ruleSet.Version,
			"is_active":   ruleSet.IsActive,
			"created_at":  ruleSet.CreatedAt,
			"rules_count": len(ruleSet.Rules),
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"rule_sets": result,
	})
}

// CreateScoringRuleSet создает новый набор правил скоринга
func CreateScoringRuleSet(c *gin.Context) {
	var ruleSet scoring.RuleSet
	if err := c.ShouldBindJSON(&ruleSet); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Валидация набора правил
	if ruleSet.ID == "" || ruleSet.Name == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID и название набора правил обязательны"})
		return
	}

	// Добавление набора правил в движок
	scoringEngine.AddRuleSet(&ruleSet)

	c.JSON(http.StatusCreated, gin.H{
		"message":  "Набор правил успешно создан",
		"rule_set": ruleSet,
	})
}

// UpdateScoringRuleSet обновляет набор правил скоринга
func UpdateScoringRuleSet(c *gin.Context) {
	ruleSetID := c.Param("id")

	var updateData map[string]interface{}
	if err := c.ShouldBindJSON(&updateData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Получение существующего набора правил
	ruleSets := scoringEngine.GetRuleSets()
	ruleSet, exists := ruleSets[ruleSetID]
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "Набор правил не найден"})
		return
	}

	// Обновление полей
	if name, ok := updateData["name"].(string); ok {
		ruleSet.Name = name
	}
	if isActive, ok := updateData["is_active"].(bool); ok {
		ruleSet.IsActive = isActive
	}

	// Сохранение обновленного набора правил
	scoringEngine.AddRuleSet(ruleSet)

	c.JSON(http.StatusOK, gin.H{
		"message":  "Набор правил успешно обновлен",
		"rule_set": ruleSet,
	})
}
