package main

import (
	"log"
	"net/http"
	"tenderhelp/internal/database"
	"tenderhelp/internal/handlers"
	"tenderhelp/internal/models"
	"tenderhelp/internal/scoring"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func main() {
	// Инициализация базы данных
	db := database.InitDB()
	handlers.InitDB() // Инициализация db в handlers

	// Автомиграция моделей
	db.AutoMigrate(
		&models.User{},
		&models.Client{},
		&models.Request{},
		&models.News{},
		&models.HelpCategory{},
		&models.HelpArticle{},
		// Новые модели для системы брокериджа
		&handlers.Application{},
		&handlers.StatusHistory{},
		&handlers.File{},
		&handlers.POSApplication{},
		&handlers.GuaranteeApplication{},
	)

	// Инициализация системы скоринга
	scoringEngine := scoring.NewScoringEngine()
	handlers.SetScoringEngine(scoringEngine)

	// Инициализация системы интеграций
	handlers.InitIntegrations()

	// Настройка Gin
	r := gin.Default()

	// CORS настройки
	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"http://localhost:3000"}
	config.AllowMethods = []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}
	config.AllowHeaders = []string{"Origin", "Content-Type", "Accept", "Authorization"}
	r.Use(cors.New(config))

	// Статические файлы
	r.Static("/static", "./static")
	r.LoadHTMLGlob("templates/*")

	// API маршруты
	api := r.Group("/api")
	{
		// Заявки (старая система) - только для пользователей и выше
		api.GET("/requests", handlers.RequireAuth(), handlers.RequirePermission("view_applications"), handlers.GetRequests)
		api.POST("/requests", handlers.RequireAuth(), handlers.RequirePermission("create_applications"), handlers.CreateRequest)
		api.PUT("/requests/:id", handlers.RequireAuth(), handlers.RequirePermission("manage_applications"), handlers.UpdateRequest)
		api.DELETE("/requests/:id", handlers.RequireAuth(), handlers.RequirePermission("manage_applications"), handlers.DeleteRequest)

		// Новые заявки (система брокериджа)
		api.GET("/applications", handlers.RequireAuth(), handlers.RequirePermission("view_applications"), handlers.GetApplications)
		api.POST("/applications", handlers.RequireAuth(), handlers.RequirePermission("create_applications"), handlers.CreateApplication)
		api.GET("/applications/:id", handlers.RequireAuth(), handlers.RequirePermission("view_applications"), handlers.GetApplication)
		api.PUT("/applications/:id", handlers.RequireAuth(), handlers.RequirePermission("manage_applications"), handlers.UpdateApplication)
		api.POST("/applications/:id/submit", handlers.RequireAuth(), handlers.RequirePermission("manage_applications"), handlers.SubmitApplication)

		// Файлы
		api.POST("/files/upload", handlers.UploadFile)
		api.GET("/files/presigned", handlers.GetPresignedUploadURL)
		api.GET("/applications/:id/files", handlers.GetFiles)
		api.DELETE("/files/:fileId", handlers.DeleteFile)
		api.GET("/files/:fileId", handlers.DownloadFile)

		// Скоринг
		api.POST("/scoring/run/:id", handlers.RunScoring)
		api.GET("/scoring/result/:id", handlers.GetScoringResult)

		// Интеграции с банками
		api.POST("/applications/:id/send", handlers.SendApplicationToBanks)
		api.GET("/banks/summary", handlers.GetBankSummary)
		api.GET("/banks/supported", handlers.GetSupportedBanks)
		api.GET("/banks/:bankId/availability", handlers.CheckBankAvailability)
		api.GET("/applications/:id/banks/:bankId/status", handlers.GetApplicationStatusFromBank)
		api.GET("/applications/:id/responses", handlers.GetBankResponses)

		// Очереди и задачи
		api.GET("/queue/stats", handlers.GetQueueStats)
		api.GET("/queue/tasks/:taskId", handlers.GetTaskStatus)
		api.POST("/queue/tasks", handlers.CreateTask)

		// ПОС и банковские гарантии
		api.POST("/applications/:id/pos", handlers.CreatePOSApplication)
		api.GET("/applications/:id/pos", handlers.GetPOSApplication)
		api.POST("/applications/:id/guarantee", handlers.CreateGuaranteeApplication)
		api.GET("/applications/:id/guarantee", handlers.GetGuaranteeApplication)
		api.POST("/applications/:id/pos/document", handlers.GeneratePOSDocument)
		api.POST("/applications/:id/guarantee/document", handlers.GenerateGuaranteeDocument)

		// Клиенты - только для пользователей и выше
		api.GET("/clients", handlers.RequireAuth(), handlers.RequirePermission("view_clients"), handlers.GetClients)
		api.POST("/clients", handlers.RequireAuth(), handlers.RequirePermission("manage_clients"), handlers.CreateClient)
		api.PUT("/clients/:id", handlers.RequireAuth(), handlers.RequirePermission("manage_clients"), handlers.UpdateClient)
		api.DELETE("/clients/:id", handlers.RequireAuth(), handlers.RequirePermission("manage_clients"), handlers.DeleteClient)

		// Аналитика - только для директора и выше
		api.GET("/analytics", handlers.RequireAuth(), handlers.RequirePermission("view_analytics"), handlers.GetAnalytics)

		// Новости
		api.GET("/news", handlers.GetNews)
		api.POST("/news", handlers.CreateNews)

		// Помощь
		api.GET("/help/categories", handlers.GetHelpCategories)
		api.GET("/help/articles", handlers.GetHelpArticles)
		api.GET("/help/search", handlers.SearchHelp)

		// Аутентификация
		api.POST("/login", handlers.Login)
		api.POST("/register", handlers.Register)
		api.POST("/logout", handlers.Logout)
		api.GET("/profile", handlers.RequireAuth(), handlers.GetProfile)

		// Пользователи / Регистрация
		api.POST("/register-agent", handlers.RegisterAgent)
		api.POST("/register-client", handlers.RegisterClient)
	}

	// Swagger документация
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Админские маршруты - только для админов
	admin := r.Group("/api/admin")
	admin.Use(handlers.RequireAuth(), handlers.RequireRole("admin"))
	{
		admin.GET("/banks", handlers.GetBanks)
		admin.POST("/upload-banks", handlers.UploadBanks)
		admin.GET("/analytics", handlers.GetAdminAnalytics)

		// Управление скорингом
		admin.GET("/scoring/rulesets", handlers.GetScoringRuleSets)
		admin.POST("/scoring/rulesets", handlers.CreateScoringRuleSet)
		admin.PUT("/scoring/rulesets/:id", handlers.UpdateScoringRuleSet)
	}

	// Главная страница
	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", gin.H{})
	})

	log.Println("Сервер запущен на порту :8080")
	r.Run(":8080")
}
