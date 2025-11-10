package adapters

import (
	"encoding/json"
	"fmt"
	"time"
)

// BankAdapter интерфейс для адаптеров банков
type BankAdapter interface {
	// SendApplication отправляет заявку в банк
	SendApplication(application ApplicationData) (*BankResponse, error)

	// GetApplicationStatus получает статус заявки
	GetApplicationStatus(applicationID string) (*ApplicationStatus, error)

	// GetBankInfo возвращает информацию о банке
	GetBankInfo() *BankInfo
}

// ApplicationData данные заявки для отправки в банк
type ApplicationData struct {
	ID          string          `json:"id"`
	Type        string          `json:"type"` // credit, guarantee
	Amount      float64         `json:"amount"`
	ClientData  json.RawMessage `json:"client_data"`
	Files       []FileData      `json:"files"`
	ScoringData json.RawMessage `json:"scoring_data"`
	CreatedAt   time.Time       `json:"created_at"`
}

// FileData данные файла
type FileData struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Type     string `json:"type"`
	Size     int64  `json:"size"`
	Content  []byte `json:"content"`
	MimeType string `json:"mime_type"`
}

// BankResponse ответ от банка
type BankResponse struct {
	Success       bool      `json:"success"`
	ApplicationID string    `json:"application_id"`
	BankID        string    `json:"bank_id"`
	Status        string    `json:"status"`
	Message       string    `json:"message"`
	ErrorCode     string    `json:"error_code,omitempty"`
	Timestamp     time.Time `json:"timestamp"`
}

// ApplicationStatus статус заявки в банке
type ApplicationStatus struct {
	ApplicationID string    `json:"application_id"`
	BankID        string    `json:"bank_id"`
	Status        string    `json:"status"`
	Decision      string    `json:"decision,omitempty"`
	Amount        float64   `json:"amount,omitempty"`
	Rate          float64   `json:"rate,omitempty"`
	Term          int       `json:"term,omitempty"`
	Message       string    `json:"message"`
	UpdatedAt     time.Time `json:"updated_at"`
}

// BankInfo информация о банке
type BankInfo struct {
	ID             string   `json:"id"`
	Name           string   `json:"name"`
	Code           string   `json:"code"`
	APIEndpoint    string   `json:"api_endpoint"`
	IsActive       bool     `json:"is_active"`
	SupportedTypes []string `json:"supported_types"`
	MaxAmount      float64  `json:"max_amount"`
	MinAmount      float64  `json:"min_amount"`
	ProcessingTime string   `json:"processing_time"`
}

// BaseAdapter базовая структура для адаптеров
type BaseAdapter struct {
	BankInfo BankInfo
	Config   map[string]interface{}
}

// NewBaseAdapter создает новый базовый адаптер
func NewBaseAdapter(bankInfo BankInfo, config map[string]interface{}) *BaseAdapter {
	return &BaseAdapter{
		BankInfo: bankInfo,
		Config:   config,
	}
}

// GetBankInfo возвращает информацию о банке
func (ba *BaseAdapter) GetBankInfo() *BankInfo {
	return &ba.BankInfo
}

// ValidateApplication валидирует заявку перед отправкой
func (ba *BaseAdapter) ValidateApplication(data ApplicationData) error {
	// Проверка типа заявки
	if data.Type != "credit" && data.Type != "guarantee" {
		return fmt.Errorf("неподдерживаемый тип заявки: %s", data.Type)
	}

	// Проверка суммы
	if data.Amount < ba.BankInfo.MinAmount {
		return fmt.Errorf("сумма заявки меньше минимальной: %.2f", ba.BankInfo.MinAmount)
	}

	if data.Amount > ba.BankInfo.MaxAmount {
		return fmt.Errorf("сумма заявки больше максимальной: %.2f", ba.BankInfo.MaxAmount)
	}

	// Проверка поддержки типа заявки
	supported := false
	for _, t := range ba.BankInfo.SupportedTypes {
		if t == data.Type {
			supported = true
			break
		}
	}

	if !supported {
		return fmt.Errorf("банк не поддерживает тип заявки: %s", data.Type)
	}

	return nil
}
