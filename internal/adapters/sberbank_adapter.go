package adapters

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"time"
)

// SberbankAdapter адаптер для Сбербанка (песочница)
type SberbankAdapter struct {
	*BaseAdapter
}

// NewSberbankAdapter создает новый адаптер Сбербанка
func NewSberbankAdapter() *SberbankAdapter {
	bankInfo := BankInfo{
		ID:             "sberbank",
		Name:           "Сбербанк",
		Code:           "SBER",
		APIEndpoint:    "https://api.sberbank.ru/sandbox",
		IsActive:       true,
		SupportedTypes: []string{"credit", "guarantee"},
		MaxAmount:      50000000, // 50 млн рублей
		MinAmount:      100000,   // 100 тыс рублей
		ProcessingTime: "1-3 рабочих дня",
	}

	config := map[string]interface{}{
		"timeout":     30,
		"retry_count": 3,
		"api_key":     "sandbox_key_sberbank",
		"version":     "v1.0",
	}

	baseAdapter := NewBaseAdapter(bankInfo, config)

	return &SberbankAdapter{
		BaseAdapter: baseAdapter,
	}
}

// SendApplication отправляет заявку в Сбербанк
func (sa *SberbankAdapter) SendApplication(application ApplicationData) (*BankResponse, error) {
	// Валидация заявки
	if err := sa.ValidateApplication(application); err != nil {
		return &BankResponse{
			Success:   false,
			Message:   fmt.Sprintf("Ошибка валидации: %s", err.Error()),
			ErrorCode: "VALIDATION_ERROR",
			Timestamp: time.Now(),
		}, nil
	}

	// Имитация отправки в банк
	time.Sleep(100 * time.Millisecond) // Имитация сетевой задержки

	// Генерация ответа (песочница)
	response := sa.generateMockResponse(application)

	return response, nil
}

// GetApplicationStatus получает статус заявки
func (sa *SberbankAdapter) GetApplicationStatus(applicationID string) (*ApplicationStatus, error) {
	// Имитация запроса к банку
	time.Sleep(50 * time.Millisecond)

	// Генерация случайного статуса
	statuses := []string{"processing", "approved", "rejected", "pending_documents"}
	status := statuses[rand.Intn(len(statuses))]

	response := &ApplicationStatus{
		ApplicationID: applicationID,
		BankID:        sa.BankInfo.ID,
		Status:        status,
		UpdatedAt:     time.Now(),
	}

	// Добавление деталей в зависимости от статуса
	switch status {
	case "approved":
		response.Decision = "Одобрено"
		response.Amount = 1000000 + rand.Float64()*5000000 // 1-6 млн
		response.Rate = 12.5 + rand.Float64()*5            // 12.5-17.5%
		response.Term = 12 + rand.Intn(60)                 // 1-5 лет
		response.Message = "Заявка одобрена. Ожидается подписание документов."
	case "rejected":
		response.Decision = "Отклонено"
		response.Message = "Заявка отклонена по результатам скоринга"
	case "processing":
		response.Message = "Заявка находится на рассмотрении"
	case "pending_documents":
		response.Message = "Требуются дополнительные документы"
	}

	return response, nil
}

// generateMockResponse генерирует мок-ответ от банка
func (sa *SberbankAdapter) generateMockResponse(application ApplicationData) *BankResponse {
	// Случайный успех (80% вероятность)
	success := rand.Float64() < 0.8

	response := &BankResponse{
		ApplicationID: fmt.Sprintf("SBER_%d", time.Now().Unix()),
		BankID:        sa.BankInfo.ID,
		Success:       success,
		Timestamp:     time.Now(),
	}

	if success {
		response.Status = "received"
		response.Message = "Заявка успешно принята к рассмотрению"
	} else {
		response.Status = "rejected"
		response.Message = "Заявка отклонена на этапе первичной проверки"
		response.ErrorCode = "PRIMARY_CHECK_FAILED"
	}

	return response
}

// TransformApplicationData преобразует данные заявки в формат Сбербанка
func (sa *SberbankAdapter) TransformApplicationData(data ApplicationData) map[string]interface{} {
	// Извлечение данных клиента
	var clientData map[string]interface{}
	json.Unmarshal(data.ClientData, &clientData)

	// Преобразование в формат Сбербанка
	sberbankData := map[string]interface{}{
		"application_id": data.ID,
		"product_type":   data.Type,
		"amount":         data.Amount,
		"currency":       "RUB",
		"client": map[string]interface{}{
			"personal": map[string]interface{}{
				"surname":     clientData["lastName"],
				"name":        clientData["firstName"],
				"patronymic":  clientData["middleName"],
				"birth_date":  clientData["birthDate"],
				"birth_place": clientData["birthPlace"],
				"gender":      clientData["gender"],
				"citizenship": clientData["citizenship"],
			},
			"contact": map[string]interface{}{
				"phone":       clientData["primaryPhone"],
				"email":       clientData["email"],
				"reg_address": clientData["registrationAddress"],
			},
			"employment": map[string]interface{}{
				"employer": clientData["currentJob"].(map[string]interface{})["companyName"],
				"position": clientData["currentJob"].(map[string]interface{})["position"],
				"income":   clientData["currentJob"].(map[string]interface{})["monthlyIncome"],
			},
		},
		"files":   data.Files,
		"scoring": data.ScoringData,
	}

	return sberbankData
}
