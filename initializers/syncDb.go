package initializers

import "jwt/models"

func SyncDb() {
	DB.AutoMigrate(&models.User{})
	DB.AutoMigrate(&models.Ticket{})
	DB.AutoMigrate(&models.Promotion{})
	DB.AutoMigrate(&models.Context{})
	DB.AutoMigrate(&models.ChatRoomUser{})
	DB.AutoMigrate(&models.Honey{})
}