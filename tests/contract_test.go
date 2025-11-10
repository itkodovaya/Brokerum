package tests

import (
	"encoding/json"
	"tenderhelp/internal/adapters"
	"testing"
	"time"
)

// TestBankAdapterContract тестирует контракт адаптеров банков
func TestBankAdapterContract(t *testing.T) {
	// Тестирование Сбербанка
	testBankAdapter(t, "sberbank", adapters.NewSberbankAdapter())

	// Тестирование ВТБ
	testBankAdapter(t, "vtb", adapters.NewVTBAdapter())
}

// testBankAdapter тестирует конкретный адаптер банка
func testBankAdapter(t *testing.T, bankName string, adapter adapters.BankAdapter) {
	t.Run("GetBankInfo", func(t *testing.T) {
		bankInfo := adapter.GetBankInfo()

		// Проверка обязательных полей
		if bankInfo.ID == "" {
			t.Errorf("%s: ID банка не должен быть пустым", bankName)
		}

		if bankInfo.Name == "" {
			t.Errorf("%s: Название банка не должно быть пустым", bankName)
		}

		if bankInfo.Code == "" {
			t.Errorf("%s: Код банка не должен быть пустым", bankName)
		}

		if bankInfo.APIEndpoint == "" {
			t.Errorf("%s: API endpoint не должен быть пустым", bankName)
		}

		if len(bankInfo.SupportedTypes) == 0 {
			t.Errorf("%s: Должны быть поддерживаемые типы заявок", bankName)
		}

		if bankInfo.MaxAmount <= bankInfo.MinAmount {
			t.Errorf("%s: Максимальная сумма должна быть больше минимальной", bankName)
		}

		if bankInfo.ProcessingTime == "" {
			t.Errorf("%s: Время обработки не должно быть пустым", bankName)
		}
	})

	t.Run("SendApplication", func(t *testing.T) {
		// Тестовые данные заявки
		applicationData := adapters.ApplicationData{
			ID:         "contract_test_123",
			Type:       "credit",
			Amount:     1000000,
			ClientData: json.RawMessage(`{"firstName": "Тест", "lastName": "Тестов"}`),
			CreatedAt:  time.Now(),
		}

		response, err := adapter.SendApplication(applicationData)
		if err != nil {
			t.Errorf("%s: Ошибка отправки заявки: %v", bankName, err)
			return
		}

		// Проверка обязательных полей ответа
		if response.ApplicationID == "" {
			t.Errorf("%s: ID заявки в ответе не должен быть пустым", bankName)
		}

		if response.BankID == "" {
			t.Errorf("%s: ID банка в ответе не должен быть пустым", bankName)
		}

		if response.Status == "" {
			t.Errorf("%s: Статус в ответе не должен быть пустым", bankName)
		}

		if response.Message == "" {
			t.Errorf("%s: Сообщение в ответе не должно быть пустым", bankName)
		}

		if response.Timestamp.IsZero() {
			t.Errorf("%s: Временная метка в ответе не должна быть нулевой", bankName)
		}

		// Проверка валидных статусов
		validStatuses := []string{"received", "accepted", "rejected", "processing"}
		isValidStatus := false
		for _, status := range validStatuses {
			if response.Status == status {
				isValidStatus = true
				break
			}
		}
		if !isValidStatus {
			t.Errorf("%s: Некорректный статус: %s", bankName, response.Status)
		}
	})

	t.Run("GetApplicationStatus", func(t *testing.T) {
		status, err := adapter.GetApplicationStatus("contract_test_123")
		if err != nil {
			t.Errorf("%s: Ошибка получения статуса заявки: %v", bankName, err)
			return
		}

		// Проверка обязательных полей статуса
		if status.ApplicationID == "" {
			t.Errorf("%s: ID заявки в статусе не должен быть пустым", bankName)
		}

		if status.BankID == "" {
			t.Errorf("%s: ID банка в статусе не должен быть пустым", bankName)
		}

		if status.Status == "" {
			t.Errorf("%s: Статус не должен быть пустым", bankName)
		}

		if status.Message == "" {
			t.Errorf("%s: Сообщение в статусе не должно быть пустым", bankName)
		}

		if status.UpdatedAt.IsZero() {
			t.Errorf("%s: Временная метка обновления не должна быть нулевой", bankName)
		}

		// Проверка валидных статусов
		validStatuses := []string{"processing", "approved", "rejected", "pending_documents", "under_review", "declined", "additional_info_required"}
		isValidStatus := false
		for _, validStatus := range validStatuses {
			if status.Status == validStatus {
				isValidStatus = true
				break
			}
		}
		if !isValidStatus {
			t.Errorf("%s: Некорректный статус: %s", bankName, status.Status)
		}
	})

	t.Run("SendApplicationWithInvalidData", func(t *testing.T) {
		// Тест с невалидными данными
		invalidApplicationData := adapters.ApplicationData{
			ID:         "contract_test_invalid",
			Type:       "invalid_type",
			Amount:     -1000,
			ClientData: json.RawMessage(`{}`),
			CreatedAt:  time.Now(),
		}

		response, err := adapter.SendApplication(invalidApplicationData)
		if err != nil {
			t.Errorf("%s: Ошибка не должна возникать при невалидных данных: %v", bankName, err)
			return
		}

		// Проверка, что ответ содержит информацию об ошибке
		if response.Success && response.ErrorCode == "" {
			t.Errorf("%s: При невалидных данных должен быть указан код ошибки", bankName)
		}
	})
}

// TestAdapterManagerContract тестирует контракт менеджера адаптеров
func TestAdapterManagerContract(t *testing.T) {
	manager := adapters.NewAdapterManager()

	t.Run("GetAllAdapters", func(t *testing.T) {
		allAdapters := manager.GetAllAdapters()

		if len(allAdapters) == 0 {
			t.Error("Должны быть зарегистрированные адаптеры")
		}

		// Проверка, что все адаптеры имеют корректную информацию
		for bankID, adapter := range allAdapters {
			bankInfo := adapter.GetBankInfo()
			if bankInfo.ID != bankID {
				t.Errorf("ID банка в менеджере (%s) не совпадает с ID в информации (%s)", bankID, bankInfo.ID)
			}
		}
	})

	t.Run("GetActiveAdapters", func(t *testing.T) {
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
	})

	t.Run("GetBankSummary", func(t *testing.T) {
		summary := manager.GetBankSummary()

		if len(summary) == 0 {
			t.Error("Должна быть сводка по банкам")
		}

		// Проверка, что все банки в сводке активны
		for _, bankInfo := range summary {
			if !bankInfo.IsActive {
				t.Errorf("Банк %s в сводке должен быть активным", bankInfo.ID)
			}
		}
	})

	t.Run("CheckBankAvailability", func(t *testing.T) {
		// Проверка доступности существующих банков
		activeAdapters := manager.GetActiveAdapters()
		for bankID := range activeAdapters {
			available := manager.CheckBankAvailability(bankID)
			if !available {
				t.Errorf("Банк %s должен быть доступен", bankID)
			}
		}

		// Проверка недоступности несуществующего банка
		available := manager.CheckBankAvailability("nonexistent_bank")
		if available {
			t.Error("Несуществующий банк не должен быть доступен")
		}
	})

	t.Run("GetSupportedBanksForType", func(t *testing.T) {
		// Тест для кредитов
		creditBanks := manager.GetSupportedBanksForType("credit")
		if len(creditBanks) == 0 {
			t.Error("Должны быть банки, поддерживающие кредиты")
		}

		// Проверка, что все банки действительно поддерживают кредиты
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

		// Тест для банковских гарантий
		guaranteeBanks := manager.GetSupportedBanksForType("guarantee")
		if len(guaranteeBanks) == 0 {
			t.Error("Должны быть банки, поддерживающие банковские гарантии")
		}
	})

	t.Run("SendToAllBanks", func(t *testing.T) {
		applicationData := adapters.ApplicationData{
			ID:         "contract_test_all_banks",
			Type:       "credit",
			Amount:     1000000,
			ClientData: json.RawMessage(`{"firstName": "Тест", "lastName": "Тестов"}`),
			CreatedAt:  time.Now(),
		}

		responses, err := manager.SendToAllBanks(applicationData)
		if err != nil {
			t.Errorf("Ошибка отправки во все банки: %v", err)
			return
		}

		if len(responses) == 0 {
			t.Error("Должны быть ответы от банков")
		}

		// Проверка, что все ответы содержат необходимые поля
		for i, response := range responses {
			if response.BankID == "" {
				t.Errorf("Ответ %d: BankID не должен быть пустым", i)
			}
			if response.Timestamp.IsZero() {
				t.Errorf("Ответ %d: Временная метка не должна быть нулевой", i)
			}
		}
	})

	t.Run("SendToSpecificBanks", func(t *testing.T) {
		applicationData := adapters.ApplicationData{
			ID:         "contract_test_specific_banks",
			Type:       "credit",
			Amount:     1000000,
			ClientData: json.RawMessage(`{"firstName": "Тест", "lastName": "Тестов"}`),
			CreatedAt:  time.Now(),
		}

		// Получение списка доступных банков
		activeAdapters := manager.GetActiveAdapters()
		bankIDs := make([]string, 0, len(activeAdapters))
		for bankID := range activeAdapters {
			bankIDs = append(bankIDs, bankID)
		}

		if len(bankIDs) == 0 {
			t.Skip("Нет доступных банков для тестирования")
		}

		// Отправка в первые два банка
		testBankIDs := bankIDs[:min(2, len(bankIDs))]
		responses, err := manager.SendToSpecificBanks(applicationData, testBankIDs)
		if err != nil {
			t.Errorf("Ошибка отправки в конкретные банки: %v", err)
			return
		}

		if len(responses) != len(testBankIDs) {
			t.Errorf("Ожидалось %d ответов, получено %d", len(testBankIDs), len(responses))
		}

		// Проверка, что все запрошенные банки ответили
		respondedBanks := make(map[string]bool)
		for _, response := range responses {
			respondedBanks[response.BankID] = true
		}

		for _, bankID := range testBankIDs {
			if !respondedBanks[bankID] {
				t.Errorf("Банк %s не ответил", bankID)
			}
		}
	})
}

// TestBaseAdapterContract тестирует контракт базового адаптера
func TestBaseAdapterContract(t *testing.T) {
	bankInfo := adapters.BankInfo{
		ID:             "test_bank",
		Name:           "Тестовый банк",
		Code:           "TEST",
		IsActive:       true,
		SupportedTypes: []string{"credit", "guarantee"},
		MaxAmount:      1000000,
		MinAmount:      10000,
	}

	adapter := adapters.NewBaseAdapter(bankInfo, map[string]interface{}{})

	t.Run("GetBankInfo", func(t *testing.T) {
		info := adapter.GetBankInfo()

		if info.ID != bankInfo.ID {
			t.Errorf("ID банка не совпадает: ожидался %s, получен %s", bankInfo.ID, info.ID)
		}

		if info.Name != bankInfo.Name {
			t.Errorf("Название банка не совпадает: ожидалось %s, получено %s", bankInfo.Name, info.Name)
		}
	})

	t.Run("ValidateApplication", func(t *testing.T) {
		// Тест валидной заявки
		validApplication := adapters.ApplicationData{
			Type:   "credit",
			Amount: 500000,
		}

		err := adapter.ValidateApplication(validApplication)
		if err != nil {
			t.Errorf("Валидная заявка не должна вызывать ошибку: %v", err)
		}

		// Тест невалидного типа заявки
		invalidTypeApplication := adapters.ApplicationData{
			Type:   "invalid_type",
			Amount: 500000,
		}

		err = adapter.ValidateApplication(invalidTypeApplication)
		if err == nil {
			t.Error("Невалидный тип заявки должен вызывать ошибку")
		}

		// Тест слишком маленькой суммы
		smallAmountApplication := adapters.ApplicationData{
			Type:   "credit",
			Amount: 5000,
		}

		err = adapter.ValidateApplication(smallAmountApplication)
		if err == nil {
			t.Error("Слишком маленькая сумма должна вызывать ошибку")
		}

		// Тест слишком большой суммы
		largeAmountApplication := adapters.ApplicationData{
			Type:   "credit",
			Amount: 2000000,
		}

		err = adapter.ValidateApplication(largeAmountApplication)
		if err == nil {
			t.Error("Слишком большая сумма должна вызывать ошибку")
		}
	})
}

// Вспомогательная функция для минимума
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
