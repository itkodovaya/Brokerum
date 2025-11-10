package main

import (
	"log"
	"tenderhelp/internal/database"
	"tenderhelp/internal/models"
	"time"
)

func main() {
	db := database.InitDB()

	// Создаем тестовых пользователей
	users := []models.User{
		{Email: "agent@tenderhelp.ru", Password: "password", Name: "Максим Искин Олегович", Role: "agent"},
		{Email: "manager@tenderhelp.ru", Password: "password", Name: "Сутулина Елизавета", Role: "manager"},
	}

	for _, user := range users {
		db.Create(&user)
	}

	// Создаем тестовых клиентов
	clients := []models.Client{
		{
			Name:             "ООО \"КОНТУР\"",
			INN:              "7743287540",
			OGRN:             "5187746029466",
			Phone:            "+7 (499) 213-24-40",
			Email:            "souz8@bk.ru",
			City:             "Москва",
			CreditManager:    "Кондратенко Марина",
			BGManager:        "Сутулина Елизавета",
			RegistrationDate: time.Date(2025, 10, 7, 0, 0, 0, 0, time.UTC),
		},
		{
			Name:             "ООО \"АЛЬЯНС\"",
			INN:              "7743287540",
			OGRN:             "5187746029466",
			Phone:            "+7 (499) 213-24-40",
			Email:            "souz8@bk.ru",
			City:             "Москва",
			CreditManager:    "Кондратенко Марина",
			BGManager:        "Сутулина Елизавета",
			RegistrationDate: time.Date(2025, 10, 7, 0, 0, 0, 0, time.UTC),
		},
		{
			Name:             "ООО \"ДИАН ИНВЕСТ\"",
			INN:              "2635257242",
			OGRN:             "1232600005483",
			Phone:            "+7 (918) 877-88-87",
			Email:            "dian.24@bk.ru",
			City:             "Ставрополь",
			CreditManager:    "Кондратенко Марина",
			BGManager:        "Сутулина Елизавета",
			RegistrationDate: time.Date(2024, 8, 9, 0, 0, 0, 0, time.UTC),
		},
	}

	for _, client := range clients {
		db.Create(&client)
	}

	// Создаем тестовые заявки
	var agent models.User
	db.Where("role = ?", "agent").First(&agent)

	requests := []models.Request{
		{
			Number:        "БГ (И) К 158532",
			ClientID:      1,
			AgentID:       agent.ID,
			BGManager:     "Сутулина Елизавета",
			CreditManager: "Кондратенко Марина",
			Term:          462,
			Amount:        36797117.00,
			Bank:          "",
			NoticeNumber:  "32515306322",
			LawNumber:     "223-ФЗ",
			Status:        "Черновик",
			Signature:     "",
			StartDate:     time.Date(2025, 10, 27, 0, 0, 0, 0, time.UTC),
			EndDate:       time.Date(2027, 2, 1, 0, 0, 0, 0, time.UTC),
		},
		{
			Number:        "БГ (И) 158411",
			ClientID:      2,
			AgentID:       agent.ID,
			BGManager:     "Сутулина Елизавета",
			CreditManager: "Кондратенко Марина",
			Term:          60,
			Amount:        9292749.00,
			Bank:          "",
			NoticeNumber:  "",
			LawNumber:     "",
			Status:        "Черновик",
			Signature:     "",
			StartDate:     time.Date(2025, 10, 23, 0, 0, 0, 0, time.UTC),
			EndDate:       time.Date(2025, 12, 22, 0, 0, 0, 0, time.UTC),
		},
	}

	for _, request := range requests {
		db.Create(&request)
	}

	// Создаем тестовые новости
	news := []models.News{
		{
			Title:    "Уважаемый Партнер!",
			Content:  "Для улучшения качества обслуживания наших пользователей, мы вводим единый номер WhatsApp 8-967-552-30-75, просим добавить в контакты данный номер. Теперь сообщение от Вас всегда попадет к Вашему менеджеру.",
			Category: "Новости компании",
			Date:     time.Date(2021, 10, 5, 0, 0, 0, 0, time.UTC),
		},
		{
			Title:    "Создание заявки на Кредит",
			Content:  "Инструкция по созданию заявки на кредит в системе TenderHelp.",
			Category: "Новости компании",
			Date:     time.Date(2021, 10, 5, 0, 0, 0, 0, time.UTC),
		},
	}

	for _, newsItem := range news {
		db.Create(&newsItem)
	}

	// Создаем категории помощи
	categories := []models.HelpCategory{
		{Name: "Часто задаваемые вопросы", Description: "Популярные вопросы и ответы"},
		{Name: "Условия банков", Description: "Условия работы с банками"},
		{Name: "Инструкции", Description: "Пошаговые инструкции"},
		{Name: "Регламенты", Description: "Официальные регламенты"},
	}

	for _, category := range categories {
		db.Create(&category)
	}

	// Создаем статьи помощи
	var faqCategory models.HelpCategory
	db.Where("name = ?", "Часто задаваемые вопросы").First(&faqCategory)

	articles := []models.HelpArticle{
		{
			Title:      "Где найти Акты",
			Content:    "Акты можно найти в разделе 'Документы' в личном кабинете.",
			CategoryID: faqCategory.ID,
		},
		{
			Title:      "Как сделать превышение/снижение по заявке?",
			Content:    "Для изменения суммы заявки обратитесь к вашему менеджеру.",
			CategoryID: faqCategory.ID,
		},
	}

	for _, article := range articles {
		db.Create(&article)
	}

	log.Println("Тестовые данные успешно созданы!")
}
