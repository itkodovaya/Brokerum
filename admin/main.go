package main

import (
	"log"
	"net/http"
	"tenderhelp/internal/database"
	"tenderhelp/internal/handlers"
	"tenderhelp/internal/models"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	// Инициализация базы данных
	db := database.InitDB()

	// Автомиграция моделей
	db.AutoMigrate(
		&models.User{},
		&models.Client{},
		&models.Request{},
		&models.News{},
		&models.HelpCategory{},
		&models.HelpArticle{},
	)

	// Настройка Gin
	r := gin.Default()

	// CORS настройки для админки
	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"http://localhost:3001", "http://localhost:3000"}
	config.AllowMethods = []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}
	config.AllowHeaders = []string{"Origin", "Content-Type", "Accept", "Authorization"}
	r.Use(cors.New(config))

	// Статические файлы для админки
	r.Static("/static", "./static")
	r.LoadHTMLGlob("templates/*")

	// Админские API маршруты
	admin := r.Group("/api/admin")
	{
		// Банки
		admin.GET("/banks", handlers.GetBanks)
		admin.POST("/banks", handlers.CreateBank)
		admin.PUT("/banks/:id", handlers.UpdateBank)
		admin.DELETE("/banks/:id", handlers.DeleteBank)
		admin.POST("/upload-banks", handlers.UploadBanks)

		// Пользователи
		admin.GET("/users", handlers.GetUsers)
		admin.POST("/users", handlers.CreateUser)
		admin.PUT("/users/:id", handlers.UpdateUser)
		admin.DELETE("/users/:id", handlers.DeleteUser)

		// Аналитика
		admin.GET("/analytics", handlers.GetAdminAnalytics)
		admin.GET("/stats", handlers.GetAdminStats)

		// Заявки
		admin.GET("/requests", handlers.GetAdminRequests)
		admin.PUT("/requests/:id", handlers.UpdateRequestStatus)
		admin.DELETE("/requests/:id", handlers.DeleteRequest)

		// Клиенты
		admin.GET("/clients", handlers.GetAdminClients)
		admin.PUT("/clients/:id", handlers.UpdateClient)
		admin.DELETE("/clients/:id", handlers.DeleteClient)

		// Статистика
		admin.GET("/stats", handlers.GetAdminStats)
	}

	// Главная страница админки
	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "admin.html", gin.H{})
	})

	// Страница входа в админку
	r.GET("/login", func(c *gin.Context) {
		c.HTML(http.StatusOK, "login.html", gin.H{})
	})

	log.Println("Админка запущена на порту :8081")
	r.Run(":8081")
}
