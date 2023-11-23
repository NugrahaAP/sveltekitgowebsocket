package db

import (
	"backend/models"
	"log"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func Database() (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open("./db/dev.db"), &gorm.Config{})

	if err != nil {
		panic(err)
	}

	if err = db.Debug().AutoMigrate(&models.User{}, &models.ChatRoom{}, &models.Message{}, &models.GroupChatRoom{}); err != nil {
		log.Println(err)
	}

	return db, err
}
