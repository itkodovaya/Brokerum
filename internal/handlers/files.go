package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// File представляет файл в системе
type File struct {
	ID            uint      `json:"id" gorm:"primaryKey"`
	ApplicationID uint      `json:"application_id"`
	FileName      string    `json:"file_name"`
	OriginalName  string    `json:"original_name"`
	FilePath      string    `json:"file_path"`
	FileSize      int64     `json:"file_size"`
	MimeType      string    `json:"mime_type"`
	FileHash      string    `json:"file_hash"`
	UploadedAt    time.Time `json:"uploaded_at"`
	IsDeleted     bool      `json:"is_deleted"`

	// Метаданные
	Metadata json.RawMessage `json:"metadata" gorm:"type:jsonb"`
}

// FileUploadRequest запрос на загрузку файла
type FileUploadRequest struct {
	ApplicationID uint   `json:"application_id" binding:"required"`
	Category      string `json:"category"` // passport, income, property, etc.
	Description   string `json:"description"`
}

// FileListResponse ответ со списком файлов
type FileListResponse struct {
	Files []File `json:"files"`
	Total int64  `json:"total"`
}

// PresignedUploadRequest запрос на получение presigned URL
type PresignedUploadRequest struct {
	FileName      string `json:"file_name" binding:"required"`
	FileSize      int64  `json:"file_size" binding:"required"`
	MimeType      string `json:"mime_type" binding:"required"`
	ApplicationID uint   `json:"application_id" binding:"required"`
}

// PresignedUploadResponse ответ с presigned URL
type PresignedUploadResponse struct {
	UploadURL string    `json:"upload_url"`
	FileID    uint      `json:"file_id"`
	ExpiresAt time.Time `json:"expires_at"`
}

// UploadFile загружает файл
func UploadFile(c *gin.Context) {
	// Получение файла из формы
	file, header, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Файл не найден"})
		return
	}
	defer file.Close()

	// Получение дополнительных параметров
	applicationIDStr := c.PostForm("application_id")
	category := c.PostForm("category")
	description := c.PostForm("description")

	applicationID, err := strconv.ParseUint(applicationIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Некорректный ID заявки"})
		return
	}

	// Валидация файла
	if err := validateFile(header); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Создание директории для файлов
	uploadDir := fmt.Sprintf("uploads/applications/%d", applicationID)
	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка создания директории"})
		return
	}

	// Генерация уникального имени файла
	fileExt := filepath.Ext(header.Filename)
	fileName := fmt.Sprintf("%d_%s%s", time.Now().Unix(), generateRandomString(8), fileExt)
	filePath := filepath.Join(uploadDir, fileName)

	// Сохранение файла
	dst, err := os.Create(filePath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка создания файла"})
		return
	}
	defer dst.Close()

	if _, err := io.Copy(dst, file); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка сохранения файла"})
		return
	}

	// Получение информации о файле
	fileInfo, err := dst.Stat()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка получения информации о файле"})
		return
	}

	// Вычисление хеша файла
	fileHash, err := calculateFileHash(filePath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка вычисления хеша файла"})
		return
	}

	// Создание записи в базе данных
	fileRecord := File{
		ApplicationID: uint(applicationID),
		FileName:      fileName,
		OriginalName:  header.Filename,
		FilePath:      filePath,
		FileSize:      fileInfo.Size(),
		MimeType:      header.Header.Get("Content-Type"),
		FileHash:      fileHash,
		UploadedAt:    time.Now(),
		IsDeleted:     false,
		Metadata: json.RawMessage(fmt.Sprintf(`{
			"category": "%s",
			"description": "%s",
			"upload_ip": "%s",
			"user_agent": "%s"
		}`, category, description, c.ClientIP(), c.GetHeader("User-Agent"))),
	}

	if err := db.Create(&fileRecord).Error; err != nil {
		// Удаление файла в случае ошибки БД
		os.Remove(filePath)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка сохранения информации о файле"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Файл успешно загружен",
		"file":    fileRecord,
	})
}

// GetPresignedUploadURL возвращает presigned URL для загрузки
func GetPresignedUploadURL(c *gin.Context) {
	var req PresignedUploadRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Валидация размера файла (максимум 10MB)
	if req.FileSize > 10*1024*1024 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Размер файла превышает 10MB"})
		return
	}

	// Валидация типа файла
	if !isAllowedFileType(req.MimeType) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Тип файла не поддерживается"})
		return
	}

	// Создание записи о файле в БД
	fileRecord := File{
		ApplicationID: req.ApplicationID,
		FileName:      req.FileName,
		OriginalName:  req.FileName,
		FilePath:      "", // Будет заполнено после загрузки
		FileSize:      req.FileSize,
		MimeType:      req.MimeType,
		UploadedAt:    time.Now(),
		IsDeleted:     false,
	}

	if err := db.Create(&fileRecord).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка создания записи о файле"})
		return
	}

	// TODO: Реализовать presigned URL для S3
	// Пока возвращаем заглушку
	uploadURL := fmt.Sprintf("/api/files/upload/%d", fileRecord.ID)
	expiresAt := time.Now().Add(1 * time.Hour)

	response := PresignedUploadResponse{
		UploadURL: uploadURL,
		FileID:    fileRecord.ID,
		ExpiresAt: expiresAt,
	}

	c.JSON(http.StatusOK, response)
}

// GetFiles возвращает список файлов заявки
func GetFiles(c *gin.Context) {
	applicationID := c.Param("applicationId")

	var files []File
	query := db.Where("application_id = ? AND is_deleted = ?", applicationID, false)

	if err := query.Find(&files).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка получения файлов"})
		return
	}

	response := FileListResponse{
		Files: files,
		Total: int64(len(files)),
	}

	c.JSON(http.StatusOK, response)
}

// DeleteFile удаляет файл
func DeleteFile(c *gin.Context) {
	fileID := c.Param("fileId")

	var file File
	if err := db.First(&file, fileID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Файл не найден"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка получения файла"})
		return
	}

	// Мягкое удаление
	file.IsDeleted = true
	if err := db.Save(&file).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка удаления файла"})
		return
	}

	// Удаление физического файла
	if file.FilePath != "" {
		os.Remove(file.FilePath)
	}

	c.JSON(http.StatusOK, gin.H{"message": "Файл успешно удален"})
}

// DownloadFile скачивает файл
func DownloadFile(c *gin.Context) {
	fileID := c.Param("fileId")

	var file File
	if err := db.First(&file, fileID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Файл не найден"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка получения файла"})
		return
	}

	if file.IsDeleted {
		c.JSON(http.StatusNotFound, gin.H{"error": "Файл удален"})
		return
	}

	// Проверка существования файла
	if _, err := os.Stat(file.FilePath); os.IsNotExist(err) {
		c.JSON(http.StatusNotFound, gin.H{"error": "Файл не найден на диске"})
		return
	}

	// Отправка файла
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", file.OriginalName))
	c.Header("Content-Type", file.MimeType)
	c.File(file.FilePath)
}

// validateFile валидирует загружаемый файл
func validateFile(header *multipart.FileHeader) error {
	// Проверка размера файла (максимум 10MB)
	if header.Size > 10*1024*1024 {
		return fmt.Errorf("размер файла превышает 10MB")
	}

	// Проверка типа файла
	if !isAllowedFileType(header.Header.Get("Content-Type")) {
		return fmt.Errorf("тип файла не поддерживается")
	}

	// Проверка расширения файла
	ext := strings.ToLower(filepath.Ext(header.Filename))
	allowedExts := []string{".pdf", ".jpg", ".jpeg", ".png", ".doc", ".docx", ".xls", ".xlsx"}

	allowed := false
	for _, allowedExt := range allowedExts {
		if ext == allowedExt {
			allowed = true
			break
		}
	}

	if !allowed {
		return fmt.Errorf("расширение файла не поддерживается")
	}

	return nil
}

// isAllowedFileType проверяет, разрешен ли тип файла
func isAllowedFileType(mimeType string) bool {
	allowedTypes := []string{
		"application/pdf",
		"image/jpeg",
		"image/png",
		"application/msword",
		"application/vnd.openxmlformats-officedocument.wordprocessingml.document",
		"application/vnd.ms-excel",
		"application/vnd.openxmlformats-officedocument.spreadsheetml.sheet",
	}

	for _, allowedType := range allowedTypes {
		if mimeType == allowedType {
			return true
		}
	}

	return false
}

// calculateFileHash вычисляет хеш файла
func calculateFileHash(filePath string) (string, error) {
	// TODO: Реализовать вычисление SHA-256 хеша
	// Пока возвращаем заглушку
	return fmt.Sprintf("hash_%d", time.Now().Unix()), nil
}

// generateRandomString генерирует случайную строку
func generateRandomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[time.Now().UnixNano()%int64(len(charset))]
	}
	return string(b)
}
