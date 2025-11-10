package adapters

import (
	"encoding/json"
	"testing"
	"time"
)

func TestSberbankAdapter_SendApplication(t *testing.T) {
	adapter := NewSberbankAdapter()

	// Тестовые данные заявки
	applicationData := ApplicationData{
		ID:         "test_123",
		Type:       "credit",
		Amount:     1000000,
		ClientData: json.RawMessage(`{"firstName": "Иван", "lastName": "Иванов"}`),
		CreatedAt:  time.Now(),
	}

	// Отправка заявки
	response, err := adapter.SendApplication(applicationData)
	if err != nil {
		t.Fatalf("Ошибка отправки заявки: %v", err)
	}

	// Проверки ответа
	if response.ApplicationID == "" {
		t.Error("ID заявки не должен быть пустым")
	}

	if response.BankID != "sberbank" {
		t.Errorf("Ожидался BankID 'sberbank', получен '%s'", response.BankID)
	}

	if response.Timestamp.IsZero() {
		t.Error("Временная метка не установлена")
	}
}

func TestSberbankAdapter_GetApplicationStatus(t *testing.T) {
	adapter := NewSberbankAdapter()

	// Получение статуса заявки
	status, err := adapter.GetApplicationStatus("test_123")
	if err != nil {
		t.Fatalf("Ошибка получения статуса: %v", err)
	}

	// Проверки статуса
	if status.ApplicationID != "test_123" {
		t.Errorf("Ожидался ApplicationID 'test_123', получен '%s'", status.ApplicationID)
	}

	if status.BankID != "sberbank" {
		t.Errorf("Ожидался BankID 'sberbank', получен '%s'", status.BankID)
	}

	if status.Status == "" {
		t.Error("Статус не должен быть пустым")
	}

	if status.UpdatedAt.IsZero() {
		t.Error("Временная метка обновления не установлена")
	}
}

func TestSberbankAdapter_GetBankInfo(t *testing.T) {
	adapter := NewSberbankAdapter()

	bankInfo := adapter.GetBankInfo()

	if bankInfo.ID != "sberbank" {
		t.Errorf("Ожидался ID 'sberbank', получен '%s'", bankInfo.ID)
	}

	if bankInfo.Name != "Сбербанк" {
		t.Errorf("Ожидалось название 'Сбербанк', получено '%s'", bankInfo.Name)
	}

	if !bankInfo.IsActive {
		t.Error("Банк должен быть активным")
	}

	if len(bankInfo.SupportedTypes) == 0 {
		t.Error("Должны быть поддерживаемые типы заявок")
	}

	if bankInfo.MaxAmount <= bankInfo.MinAmount {
		t.Error("Максимальная сумма должна быть больше минимальной")
	}
}

func TestVTBAdapter_SendApplication(t *testing.T) {
	adapter := NewVTBAdapter()

	// Тестовые данные заявки
	applicationData := ApplicationData{
		ID:         "test_456",
		Type:       "guarantee",
		Amount:     500000,
		ClientData: json.RawMessage(`{"firstName": "Петр", "lastName": "Петров"}`),
		CreatedAt:  time.Now(),
	}

	// Отправка заявки
	response, err := adapter.SendApplication(applicationData)
	if err != nil {
		t.Fatalf("Ошибка отправки заявки: %v", err)
	}

	// Проверки ответа
	if response.ApplicationID == "" {
		t.Error("ID заявки не должен быть пустым")
	}

	if response.BankID != "vtb" {
		t.Errorf("Ожидался BankID 'vtb', получен '%s'", response.BankID)
	}
}

func TestVTBAdapter_GetApplicationStatus(t *testing.T) {
	adapter := NewVTBAdapter()

	// Получение статуса заявки
	status, err := adapter.GetApplicationStatus("test_456")
	if err != nil {
		t.Fatalf("Ошибка получения статуса: %v", err)
	}

	// Проверки статуса
	if status.ApplicationID != "test_456" {
		t.Errorf("Ожидался ApplicationID 'test_456', получен '%s'", status.ApplicationID)
	}

	if status.BankID != "vtb" {
		t.Errorf("Ожидался BankID 'vtb', получен '%s'", status.BankID)
	}
}

func TestAdapterManager_RegisterAdapter(t *testing.T) {
	manager := NewAdapterManager()

	// Создание тестового адаптера
	testAdapter := NewSberbankAdapter()

	// Регистрация адаптера
	manager.RegisterAdapter(testAdapter)

	// Проверка, что адаптер зарегистрирован
	adapter, err := manager.GetAdapter("sberbank")
	if err != nil {
		t.Fatalf("Ошибка получения адаптера: %v", err)
	}

	if adapter == nil {
		t.Error("Адаптер не должен быть nil")
	}
}

func TestAdapterManager_GetActiveAdapters(t *testing.T) {
	manager := NewAdapterManager()

	activeAdapters := manager.GetActiveAdapters()

	if len(activeAdapters) == 0 {
		t.Error("Должны быть активные адаптеры")
	}

	// Проверка, что все адаптеры активны
	for bankID, adapter := range activeAdapters {
		bankInfo := adapter.GetBankInfo()
		if !bankInfo.IsActive {
			t.Errorf("Адаптер %s должен быть активным", bankID)
		}
	}
}

func TestAdapterManager_SendToAllBanks(t *testing.T) {
	manager := NewAdapterManager()

	// Тестовые данные заявки
	applicationData := ApplicationData{
		ID:         "test_all_banks",
		Type:       "credit",
		Amount:     1000000,
		ClientData: json.RawMessage(`{"firstName": "Тест", "lastName": "Тестов"}`),
		CreatedAt:  time.Now(),
	}

	// Отправка во все банки
	responses, err := manager.SendToAllBanks(applicationData)
	if err != nil {
		t.Fatalf("Ошибка отправки во все банки: %v", err)
	}

	if len(responses) == 0 {
		t.Error("Должны быть ответы от банков")
	}

	// Проверка, что все ответы содержат необходимые поля
	for _, response := range responses {
		if response.BankID == "" {
			t.Error("BankID не должен быть пустым")
		}

		if response.Timestamp.IsZero() {
			t.Error("Временная метка не установлена")
		}
	}
}

func TestAdapterManager_SendToSpecificBanks(t *testing.T) {
	manager := NewAdapterManager()

	// Тестовые данные заявки
	applicationData := ApplicationData{
		ID:         "test_specific_banks",
		Type:       "credit",
		Amount:     1000000,
		ClientData: json.RawMessage(`{"firstName": "Тест", "lastName": "Тестов"}`),
		CreatedAt:  time.Now(),
	}

	// Отправка в конкретные банки
	bankIDs := []string{"sberbank", "vtb"}
	responses, err := manager.SendToSpecificBanks(applicationData, bankIDs)
	if err != nil {
		t.Fatalf("Ошибка отправки в конкретные банки: %v", err)
	}

	if len(responses) != len(bankIDs) {
		t.Errorf("Ожидалось %d ответов, получено %d", len(bankIDs), len(responses))
	}

	// Проверка, что все запрошенные банки ответили
	respondedBanks := make(map[string]bool)
	for _, response := range responses {
		respondedBanks[response.BankID] = true
	}

	for _, bankID := range bankIDs {
		if !respondedBanks[bankID] {
			t.Errorf("Банк %s не ответил", bankID)
		}
	}
}

func TestAdapterManager_GetBankSummary(t *testing.T) {
	manager := NewAdapterManager()

	summary := manager.GetBankSummary()

	if len(summary) == 0 {
		t.Error("Должна быть сводка по банкам")
	}

	// Проверка, что все банки в сводке активны
	for _, bankInfo := range summary {
		if !bankInfo.IsActive {
			t.Errorf("Банк %s должен быть активным", bankInfo.ID)
		}

		if bankInfo.Name == "" {
			t.Error("Название банка не должно быть пустым")
		}

		if bankInfo.MaxAmount <= bankInfo.MinAmount {
			t.Error("Максимальная сумма должна быть больше минимальной")
		}
	}
}

func TestAdapterManager_GetSupportedBanksForType(t *testing.T) {
	manager := NewAdapterManager()

	// Тестирование для типа "credit"
	creditBanks := manager.GetSupportedBanksForType("credit")
	if len(creditBanks) == 0 {
		t.Error("Должны быть банки, поддерживающие кредиты")
	}

	// Проверка, что все банки поддерживают кредиты
	for _, bankInfo := range creditBanks {
		supportsCredit := false
		for _, supportedType := range bankInfo.SupportedTypes {
			if supportedType == "credit" {
				supportsCredit = true
				break
			}
		}
		if !supportsCredit {
			t.Errorf("Банк %s должен поддерживать кредиты", bankInfo.ID)
		}
	}

	// Тестирование для типа "guarantee"
	guaranteeBanks := manager.GetSupportedBanksForType("guarantee")
	if len(guaranteeBanks) == 0 {
		t.Error("Должны быть банки, поддерживающие гарантии")
	}
}

func TestAdapterManager_CheckBankAvailability(t *testing.T) {
	manager := NewAdapterManager()

	// Проверка доступности существующего банка
	available := manager.CheckBankAvailability("sberbank")
	if !available {
		t.Error("Сбербанк должен быть доступен")
	}

	// Проверка доступности несуществующего банка
	available = manager.CheckBankAvailability("nonexistent")
	if available {
		t.Error("Несуществующий банк не должен быть доступен")
	}
}

func TestBaseAdapter_ValidateApplication(t *testing.T) {
	bankInfo := BankInfo{
		ID:             "test_bank",
		Name:           "Тестовый банк",
		IsActive:       true,
		SupportedTypes: []string{"credit"},
		MaxAmount:      1000000,
		MinAmount:      10000,
	}

	adapter := NewBaseAdapter(bankInfo, map[string]interface{}{})

	// Тестирование валидной заявки
	validApplication := ApplicationData{
		Type:   "credit",
		Amount: 500000,
	}

	err := adapter.ValidateApplication(validApplication)
	if err != nil {
		t.Errorf("Валидная заявка не должна вызывать ошибку: %v", err)
	}

	// Тестирование невалидного типа заявки
	invalidTypeApplication := ApplicationData{
		Type:   "invalid_type",
		Amount: 500000,
	}

	err = adapter.ValidateApplication(invalidTypeApplication)
	if err == nil {
		t.Error("Невалидный тип заявки должен вызывать ошибку")
	}

	// Тестирование слишком маленькой суммы
	smallAmountApplication := ApplicationData{
		Type:   "credit",
		Amount: 5000,
	}

	err = adapter.ValidateApplication(smallAmountApplication)
	if err == nil {
		t.Error("Слишком маленькая сумма должна вызывать ошибку")
	}

	// Тестирование слишком большой суммы
	largeAmountApplication := ApplicationData{
		Type:   "credit",
		Amount: 2000000,
	}

	err = adapter.ValidateApplication(largeAmountApplication)
	if err == nil {
		t.Error("Слишком большая сумма должна вызывать ошибку")
	}
}

func BenchmarkSberbankAdapter_SendApplication(b *testing.B) {
	adapter := NewSberbankAdapter()

	applicationData := ApplicationData{
		ID:         "benchmark_test",
		Type:       "credit",
		Amount:     1000000,
		ClientData: json.RawMessage(`{"firstName": "Тест", "lastName": "Тестов"}`),
		CreatedAt:  time.Now(),
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, err := adapter.SendApplication(applicationData)
		if err != nil {
			b.Fatalf("Ошибка отправки заявки: %v", err)
		}
	}
}

func BenchmarkAdapterManager_SendToAllBanks(b *testing.B) {
	manager := NewAdapterManager()

	applicationData := ApplicationData{
		ID:         "benchmark_all_banks",
		Type:       "credit",
		Amount:     1000000,
		ClientData: json.RawMessage(`{"firstName": "Тест", "lastName": "Тестов"}`),
		CreatedAt:  time.Now(),
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, err := manager.SendToAllBanks(applicationData)
		if err != nil {
			b.Fatalf("Ошибка отправки во все банки: %v", err)
		}
	}
}
