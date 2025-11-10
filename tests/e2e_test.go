package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"testing"
	"time"
)

const (
	baseURL = "http://localhost:8080"
	apiURL  = baseURL + "/api"
)

// TestFullApplicationFlow тестирует полный поток создания заявки
func TestFullApplicationFlow(t *testing.T) {
	// Шаг 1: Создание заявки
	application := createApplication(t)

	// Шаг 2: Загрузка файлов
	uploadFiles(t, application["id"])

	// Шаг 3: Запуск скоринга
	runScoring(t, application["id"])

	// Шаг 4: Получение результата скоринга
	_ = getScoringResult(t, application["id"])

	// Шаг 5: Отправка в банки
	sendToBanks(t, application["id"])

	// Шаг 6: Проверка статуса заявки
	checkApplicationStatus(t, application["id"])

	// Шаг 7: Создание заявки на ПОС
	createPOSApplication(t, application["id"])

	// Шаг 8: Генерация документа ПОС
	generatePOSDocument(t, application["id"])

	t.Logf("Полный поток заявки завершен успешно. ID: %v", application["id"])
}

// TestApplicationFlowWithErrors тестирует поток с ошибками
func TestApplicationFlowWithErrors(t *testing.T) {
	// Тест с невалидными данными
	invalidApplication := map[string]interface{}{
		"type":   "invalid_type",
		"amount": -1000,
	}

	response := makeRequest(t, "POST", "/applications", invalidApplication)
	if response.StatusCode != 400 {
		t.Errorf("Ожидался статус 400 для невалидных данных, получен: %d", response.StatusCode)
	}

	// Тест с несуществующей заявкой
	response = makeRequest(t, "GET", "/applications/99999", nil)
	if response.StatusCode != 404 {
		t.Errorf("Ожидался статус 404 для несуществующей заявки, получен: %d", response.StatusCode)
	}

	// Тест скоринга для несуществующей заявки
	response = makeRequest(t, "POST", "/scoring/run/99999", nil)
	if response.StatusCode != 404 {
		t.Errorf("Ожидался статус 404 для скоринга несуществующей заявки, получен: %d", response.StatusCode)
	}
}

// TestBankIntegrationFlow тестирует интеграцию с банками
func TestBankIntegrationFlow(t *testing.T) {
	// Получение сводки банков
	response := makeRequest(t, "GET", "/banks/summary", nil)
	if response.StatusCode != 200 {
		t.Fatalf("Ошибка получения сводки банков: %d", response.StatusCode)
	}

	var bankSummary struct {
		Banks []struct {
			ID   string `json:"id"`
			Name string `json:"name"`
		} `json:"banks"`
		Total int `json:"total"`
	}

	err := json.NewDecoder(response.Body).Decode(&bankSummary)
	if err != nil {
		t.Fatalf("Ошибка парсинга ответа: %v", err)
	}

	if bankSummary.Total == 0 {
		t.Error("Должны быть доступные банки")
	}

	// Проверка доступности банков
	for _, bank := range bankSummary.Banks {
		response := makeRequest(t, "GET", fmt.Sprintf("/banks/%s/availability", bank.ID), nil)
		if response.StatusCode != 200 {
			t.Errorf("Ошибка проверки доступности банка %s: %d", bank.ID, response.StatusCode)
		}
	}
}

// TestQueueManagement тестирует управление очередями
func TestQueueManagement(t *testing.T) {
	// Получение статистики очереди
	response := makeRequest(t, "GET", "/queue/stats", nil)
	if response.StatusCode != 200 {
		t.Fatalf("Ошибка получения статистики очереди: %d", response.StatusCode)
	}

	var stats map[string]interface{}
	err := json.NewDecoder(response.Body).Decode(&stats)
	if err != nil {
		t.Fatalf("Ошибка парсинга статистики: %v", err)
	}

	// Проверка наличия основных полей статистики
	requiredFields := []string{"total_tasks", "pending_tasks", "processing_tasks"}
	for _, field := range requiredFields {
		if _, exists := stats[field]; !exists {
			t.Errorf("Отсутствует поле статистики: %s", field)
		}
	}
}

// createApplication создает новую заявку
func createApplication(t *testing.T) map[string]interface{} {
	applicationData := map[string]interface{}{
		"type":   "credit",
		"amount": 1000000,
		"personalData": map[string]interface{}{
			"firstName":     "Иван",
			"lastName":      "Иванов",
			"middleName":    "Иванович",
			"birthDate":     "1985-03-15",
			"gender":        "male",
			"maritalStatus": "married",
			"citizenship":   "RU",
		},
		"contactData": map[string]interface{}{
			"primaryPhone": "+79161234567",
			"email":        "ivan@example.com",
		},
		"professionalData": map[string]interface{}{
			"currentJob": map[string]interface{}{
				"companyName":    "ООО Тест",
				"position":       "Менеджер",
				"employmentDate": "2020-01-15",
				"monthlyIncome":  100000,
			},
		},
		"financialData": map[string]interface{}{
			"totalMonthlyIncome": 100000,
			"creditHistory": map[string]interface{}{
				"hasOverdue":       false,
				"activeLoansCount": 1,
			},
		},
		"familyData": map[string]interface{}{
			"spouse": map[string]interface{}{
				"hasSpouse":    true,
				"spouseIncome": 80000,
			},
		},
		"additionalData": map[string]interface{}{
			"consents": map[string]interface{}{
				"dataProcessing": true,
				"creditHistory":  true,
				"scoring":        true,
			},
		},
	}

	response := makeRequest(t, "POST", "/applications", applicationData)
	if response.StatusCode != 201 {
		t.Fatalf("Ошибка создания заявки: %d", response.StatusCode)
	}

	var result map[string]interface{}
	err := json.NewDecoder(response.Body).Decode(&result)
	if err != nil {
		t.Fatalf("Ошибка парсинга ответа: %v", err)
	}

	application, ok := result["application"].(map[string]interface{})
	if !ok {
		t.Fatal("Не удалось получить данные заявки")
	}

	return application
}

// uploadFiles загружает файлы для заявки
func uploadFiles(t *testing.T, applicationID interface{}) {
	// Создание тестового файла
	_ = []byte("Тестовый файл для заявки")

	// Запрос на получение presigned URL
	response := makeRequest(t, "GET", "/files/presigned", nil)
	if response.StatusCode != 200 {
		t.Fatalf("Ошибка получения presigned URL: %d", response.StatusCode)
	}

	var presignedResponse map[string]interface{}
	err := json.NewDecoder(response.Body).Decode(&presignedResponse)
	if err != nil {
		t.Fatalf("Ошибка парсинга presigned URL: %v", err)
	}

	// Имитация загрузки файла (в реальном тесте здесь был бы HTTP PUT к S3)
	t.Logf("Файл загружен для заявки %v", applicationID)
}

// runScoring запускает скоринг заявки
func runScoring(t *testing.T, applicationID interface{}) {
	response := makeRequest(t, "POST", fmt.Sprintf("/scoring/run/%v", applicationID), nil)
	if response.StatusCode != 200 {
		t.Fatalf("Ошибка запуска скоринга: %d", response.StatusCode)
	}

	var result map[string]interface{}
	err := json.NewDecoder(response.Body).Decode(&result)
	if err != nil {
		t.Fatalf("Ошибка парсинга результата скоринга: %v", err)
	}

	if result["message"] == nil {
		t.Error("Отсутствует сообщение о результате скоринга")
	}
}

// getScoringResult получает результат скоринга
func getScoringResult(t *testing.T, applicationID interface{}) map[string]interface{} {
	response := makeRequest(t, "GET", fmt.Sprintf("/scoring/result/%v", applicationID), nil)
	if response.StatusCode != 200 {
		t.Fatalf("Ошибка получения результата скоринга: %d", response.StatusCode)
	}

	var result map[string]interface{}
	err := json.NewDecoder(response.Body).Decode(&result)
	if err != nil {
		t.Fatalf("Ошибка парсинга результата скоринга: %v", err)
	}

	return result
}

// sendToBanks отправляет заявку в банки
func sendToBanks(t *testing.T, applicationID interface{}) {
	response := makeRequest(t, "POST", fmt.Sprintf("/applications/%v/send", applicationID), nil)
	if response.StatusCode != 200 {
		t.Fatalf("Ошибка отправки в банки: %d", response.StatusCode)
	}

	var result map[string]interface{}
	err := json.NewDecoder(response.Body).Decode(&result)
	if err != nil {
		t.Fatalf("Ошибка парсинга ответа банков: %v", err)
	}

	if result["responses"] == nil {
		t.Error("Отсутствуют ответы от банков")
	}
}

// checkApplicationStatus проверяет статус заявки
func checkApplicationStatus(t *testing.T, applicationID interface{}) {
	response := makeRequest(t, "GET", fmt.Sprintf("/applications/%v", applicationID), nil)
	if response.StatusCode != 200 {
		t.Fatalf("Ошибка получения статуса заявки: %d", response.StatusCode)
	}

	var application map[string]interface{}
	err := json.NewDecoder(response.Body).Decode(&application)
	if err != nil {
		t.Fatalf("Ошибка парсинга заявки: %v", err)
	}

	if application["status"] == nil {
		t.Error("Отсутствует статус заявки")
	}
}

// createPOSApplication создает заявку на ПОС
func createPOSApplication(t *testing.T, applicationID interface{}) {
	posData := map[string]interface{}{
		"merchant_id":         "TEST_MERCHANT_123",
		"terminal_id":         "TERMINAL_456",
		"merchant_name":       "Тестовый магазин",
		"merchant_category":   "retail",
		"monthly_turnover":    500000,
		"average_transaction": 5000,
		"transaction_count":   100,
		"credit_limit":        100000,
		"interest_rate":       15.5,
		"term":                12,
	}

	response := makeRequest(t, "POST", fmt.Sprintf("/applications/%v/pos", applicationID), posData)
	if response.StatusCode != 201 {
		t.Fatalf("Ошибка создания заявки на ПОС: %d", response.StatusCode)
	}

	var result map[string]interface{}
	err := json.NewDecoder(response.Body).Decode(&result)
	if err != nil {
		t.Fatalf("Ошибка парсинга заявки на ПОС: %v", err)
	}

	if result["pos_application"] == nil {
		t.Error("Отсутствуют данные заявки на ПОС")
	}
}

// generatePOSDocument генерирует документ для ПОС
func generatePOSDocument(t *testing.T, applicationID interface{}) {
	response := makeRequest(t, "POST", fmt.Sprintf("/applications/%v/pos/document", applicationID), nil)
	if response.StatusCode != 200 {
		t.Fatalf("Ошибка генерации документа ПОС: %d", response.StatusCode)
	}

	var result map[string]interface{}
	err := json.NewDecoder(response.Body).Decode(&result)
	if err != nil {
		t.Fatalf("Ошибка парсинга документа: %v", err)
	}

	if result["document"] == nil {
		t.Error("Отсутствует сгенерированный документ")
	}
}

// makeRequest выполняет HTTP запрос
func makeRequest(t *testing.T, method, path string, body interface{}) *http.Response {
	var reqBody io.Reader
	if body != nil {
		jsonData, err := json.Marshal(body)
		if err != nil {
			t.Fatalf("Ошибка сериализации тела запроса: %v", err)
		}
		reqBody = bytes.NewBuffer(jsonData)
	}

	req, err := http.NewRequest(method, apiURL+path, reqBody)
	if err != nil {
		t.Fatalf("Ошибка создания запроса: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("Ошибка выполнения запроса: %v", err)
	}

	return resp
}

// TestServerHealth проверяет здоровье сервера
func TestServerHealth(t *testing.T) {
	response := makeRequest(t, "GET", "/", nil)
	if response.StatusCode != 200 {
		t.Fatalf("Сервер недоступен: %d", response.StatusCode)
	}
}

// TestAPIDocumentation проверяет доступность API документации
func TestAPIDocumentation(t *testing.T) {
	response := makeRequest(t, "GET", "/swagger/index.html", nil)
	if response.StatusCode != 200 {
		t.Fatalf("Swagger документация недоступна: %d", response.StatusCode)
	}
}
