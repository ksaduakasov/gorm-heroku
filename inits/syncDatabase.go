package inits

import "kettkal/models"

func SyncDB() {
	DB.AutoMigrate(&models.User{})
}
