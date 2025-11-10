package adapters

import (
	"fmt"
	"sync"
	"time"
)

// AdapterManager менеджер адаптеров банков
type AdapterManager struct {
	adapters map[string]BankAdapter
	mutex    sync.RWMutex
}

// NewAdapterManager создает новый менеджер адаптеров
func NewAdapterManager() *AdapterManager {
	manager := &AdapterManager{
		adapters: make(map[string]BankAdapter),
	}

	// Инициализация адаптеров-песочниц
	manager.RegisterAdapter(NewSberbankAdapter())
	manager.RegisterAdapter(NewVTBAdapter())

	return manager
}

// RegisterAdapter регистрирует адаптер банка
func (am *AdapterManager) RegisterAdapter(adapter BankAdapter) {
	am.mutex.Lock()
	defer am.mutex.Unlock()

	bankInfo := adapter.GetBankInfo()
	am.adapters[bankInfo.ID] = adapter
}

// GetAdapter возвращает адаптер по ID банка
func (am *AdapterManager) GetAdapter(bankID string) (BankAdapter, error) {
	am.mutex.RLock()
	defer am.mutex.RUnlock()

	adapter, exists := am.adapters[bankID]
	if !exists {
		return nil, fmt.Errorf("адаптер для банка %s не найден", bankID)
	}

	return adapter, nil
}

// GetAllAdapters возвращает все доступные адаптеры
func (am *AdapterManager) GetAllAdapters() map[string]BankAdapter {
	am.mutex.RLock()
	defer am.mutex.RUnlock()

	// Создаем копию для безопасности
	result := make(map[string]BankAdapter)
	for id, adapter := range am.adapters {
		result[id] = adapter
	}

	return result
}

// GetActiveAdapters возвращает только активные адаптеры
func (am *AdapterManager) GetActiveAdapters() map[string]BankAdapter {
	am.mutex.RLock()
	defer am.mutex.RUnlock()

	result := make(map[string]BankAdapter)
	for id, adapter := range am.adapters {
		bankInfo := adapter.GetBankInfo()
		if bankInfo.IsActive {
			result[id] = adapter
		}
	}

	return result
}

// SendToAllBanks отправляет заявку во все активные банки
func (am *AdapterManager) SendToAllBanks(application ApplicationData) ([]BankResponse, error) {
	activeAdapters := am.GetActiveAdapters()
	responses := make([]BankResponse, 0, len(activeAdapters))

	for bankID, adapter := range activeAdapters {
		response, err := adapter.SendApplication(application)
		if err != nil {
			// Логируем ошибку, но продолжаем с другими банками
			response = &BankResponse{
				Success:   false,
				BankID:    bankID,
				Message:   fmt.Sprintf("Ошибка отправки: %s", err.Error()),
				ErrorCode: "ADAPTER_ERROR",
				Timestamp: time.Now(),
			}
		}

		responses = append(responses, *response)
	}

	return responses, nil
}

// SendToSpecificBanks отправляет заявку в указанные банки
func (am *AdapterManager) SendToSpecificBanks(application ApplicationData, bankIDs []string) ([]BankResponse, error) {
	responses := make([]BankResponse, 0, len(bankIDs))

	for _, bankID := range bankIDs {
		adapter, err := am.GetAdapter(bankID)
		if err != nil {
			responses = append(responses, BankResponse{
				Success:   false,
				BankID:    bankID,
				Message:   err.Error(),
				ErrorCode: "ADAPTER_NOT_FOUND",
				Timestamp: time.Now(),
			})
			continue
		}

		response, err := adapter.SendApplication(application)
		if err != nil {
			response = &BankResponse{
				Success:   false,
				BankID:    bankID,
				Message:   fmt.Sprintf("Ошибка отправки: %s", err.Error()),
				ErrorCode: "SEND_ERROR",
				Timestamp: time.Now(),
			}
		}

		responses = append(responses, *response)
	}

	return responses, nil
}

// GetBankSummary возвращает сводку по банкам
func (am *AdapterManager) GetBankSummary() []BankInfo {
	activeAdapters := am.GetActiveAdapters()
	summary := make([]BankInfo, 0, len(activeAdapters))

	for _, adapter := range activeAdapters {
		bankInfo := adapter.GetBankInfo()
		summary = append(summary, *bankInfo)
	}

	return summary
}

// CheckBankAvailability проверяет доступность банка
func (am *AdapterManager) CheckBankAvailability(bankID string) bool {
	adapter, err := am.GetAdapter(bankID)
	if err != nil {
		return false
	}

	bankInfo := adapter.GetBankInfo()
	return bankInfo.IsActive
}

// GetSupportedBanksForType возвращает банки, поддерживающие определенный тип заявки
func (am *AdapterManager) GetSupportedBanksForType(applicationType string) []BankInfo {
	activeAdapters := am.GetActiveAdapters()
	supportedBanks := make([]BankInfo, 0)

	for _, adapter := range activeAdapters {
		bankInfo := adapter.GetBankInfo()
		for _, supportedType := range bankInfo.SupportedTypes {
			if supportedType == applicationType {
				supportedBanks = append(supportedBanks, *bankInfo)
				break
			}
		}
	}

	return supportedBanks
}
