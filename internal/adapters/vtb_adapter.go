package adapters

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"time"
)

// VTBAdapter адаптер для ВТБ (песочница)
type VTBAdapter struct {
	*BaseAdapter
}

// NewVTBAdapter создает новый адаптер ВТБ
func NewVTBAdapter() *VTBAdapter {
	bankInfo := BankInfo{
		ID:             "vtb",
		Name:           "ВТБ",
		Code:           "VTB",
		APIEndpoint:    "https://api.vtb.ru/sandbox",
		IsActive:       true,
		SupportedTypes: []string{"credit", "guarantee"},
		MaxAmount:      30000000, // 30 млн рублей
		MinAmount:      50000,    // 50 тыс рублей
		ProcessingTime: "2-5 рабочих дней",
	}

	config := map[string]interface{}{
		"timeout":     45,
		"retry_count": 5,
		"api_key":     "sandbox_key_vtb",
		"version":     "v2.1",
	}

	baseAdapter := NewBaseAdapter(bankInfo, config)

	return &VTBAdapter{
		BaseAdapter: baseAdapter,
	}
}

// SendApplication отправляет заявку в ВТБ
func (va *VTBAdapter) SendApplication(application ApplicationData) (*BankResponse, error) {
	// Валидация заявки
	if err := va.ValidateApplication(application); err != nil {
		return &BankResponse{
			Success:   false,
			Message:   fmt.Sprintf("Ошибка валидации: %s", err.Error()),
			ErrorCode: "VALIDATION_ERROR",
			Timestamp: time.Now(),
		}, nil
	}

	// Имитация отправки в банк
	time.Sleep(150 * time.Millisecond) // Имитация сетевой задержки

	// Генерация ответа (песочница)
	response := va.generateMockResponse(application)

	return response, nil
}

// GetApplicationStatus получает статус заявки
func (va *VTBAdapter) GetApplicationStatus(applicationID string) (*ApplicationStatus, error) {
	// Имитация запроса к банку
	time.Sleep(75 * time.Millisecond)

	// Генерация случайного статуса
	statuses := []string{"under_review", "approved", "declined", "additional_info_required"}
	status := statuses[rand.Intn(len(statuses))]

	response := &ApplicationStatus{
		ApplicationID: applicationID,
		BankID:        va.BankInfo.ID,
		Status:        status,
		UpdatedAt:     time.Now(),
	}

	// Добавление деталей в зависимости от статуса
	switch status {
	case "approved":
		response.Decision = "Одобрено"
		response.Amount = 500000 + rand.Float64()*3000000 // 0.5-3.5 млн
		response.Rate = 11.9 + rand.Float64()*4           // 11.9-15.9%
		response.Term = 6 + rand.Intn(48)                 // 0.5-4 года
		response.Message = "Заявка одобрена. Менеджер свяжется с вами в течение 24 часов."
	case "declined":
		response.Decision = "Отклонено"
		response.Message = "Заявка отклонена по результатам анализа"
	case "under_review":
		response.Message = "Заявка находится на рассмотрении кредитного комитета"
	case "additional_info_required":
		response.Message = "Требуется предоставить дополнительные документы"
	}

	return response, nil
}

// generateMockResponse генерирует мок-ответ от банка
func (va *VTBAdapter) generateMockResponse(application ApplicationData) *BankResponse {
	// Случайный успех (75% вероятность)
	success := rand.Float64() < 0.75

	response := &BankResponse{
		ApplicationID: fmt.Sprintf("VTB_%d", time.Now().Unix()),
		BankID:        va.BankInfo.ID,
		Success:       success,
		Timestamp:     time.Now(),
	}

	if success {
		response.Status = "accepted"
		response.Message = "Заявка принята к рассмотрению"
	} else {
		response.Status = "rejected"
		response.Message = "Заявка не прошла первичную проверку"
		response.ErrorCode = "INITIAL_SCREENING_FAILED"
	}

	return response
}

// TransformApplicationData преобразует данные заявки в формат ВТБ
func (va *VTBAdapter) TransformApplicationData(data ApplicationData) map[string]interface{} {
	// Извлечение данных клиента
	var clientData map[string]interface{}
	json.Unmarshal(data.ClientData, &clientData)

	// Преобразование в формат ВТБ
	vtbData := map[string]interface{}{
		"request_id": data.ID,
		"product":    data.Type,
		"amount":     data.Amount,
		"currency":   "RUB",
		"applicant": map[string]interface{}{
			"personal_info": map[string]interface{}{
				"last_name":      clientData["lastName"],
				"first_name":     clientData["firstName"],
				"middle_name":    clientData["middleName"],
				"date_of_birth":  clientData["birthDate"],
				"place_of_birth": clientData["birthPlace"],
				"sex":            clientData["gender"],
				"nationality":    clientData["citizenship"],
			},
			"contact_info": map[string]interface{}{
				"mobile_phone":       clientData["primaryPhone"],
				"email_address":      clientData["email"],
				"registered_address": clientData["registrationAddress"],
			},
			"employment_info": map[string]interface{}{
				"company_name": clientData["currentJob"].(map[string]interface{})["companyName"],
				"job_title":    clientData["currentJob"].(map[string]interface{})["position"],
				"salary":       clientData["currentJob"].(map[string]interface{})["monthlyIncome"],
			},
		},
		"attachments":    data.Files,
		"scoring_result": data.ScoringData,
	}

	return vtbData
}
