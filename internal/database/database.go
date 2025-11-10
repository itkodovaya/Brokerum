package database

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func InitDB() *gorm.DB {
	db, err := gorm.Open(sqlite.Open("tenderhelp.db"), &gorm.Config{})
	if err != nil {
		panic("Не удалось подключиться к базе данных")
	}
	return db
}
