package models

import (
	"time"
)
type Model struct {
	ID        uint			`gorm:"primary_key;autoIncrement:true;" json:"id"`
	CreatedAt time.Time		`gorm:"default:CURRENT_TIMESTAMP()"  json:"created_at"`
	UpdatedAt time.Time		`gorm:"default:CURRENT_TIMESTAMP()"  json:"updated_at"`
}

type User struct {
	Ticket			int64	`gorm:"default:5000" json:"ticket"`
	Account			string  `gorm:"type:varchar(50)" json:"account"`
	Name			string	`gorm:"type:varchar(50)" json:"name"`
	Email			string	`gorm:"type:varchar(50)" json:"email"`
	Password        string	`gorm:"type:varchar(50)" json:"password"`
	Wxid			string	`gorm:"index;type:varchar(50);" json:"wxid"`
	Avatar 			string 	`gorm:"type:varchar(200);json:"avatar"`
	IsActive	    bool	`gorm:"default:false" json:"is_active"`
	IsVerified		bool	`gorm:"default:false" json:"is_verified"`
	Model
}


type WxUser struct{
	Account			string  `json:"account"`
	Name 			string 	`json:"name"`
	Wxid 			string 	`json:"wxid"`
	Mobile			string 	`json:"mobile"`
	Avatar 			string 	`json:"headImage"`
}

