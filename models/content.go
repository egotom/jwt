package models
import (
	"time"
)



type Context struct{
	ID 			uint64 		`gorm:"primary_key;autoIncrement:true;" json:"id"`
	Vendor		uint8 		`gorm:"" json:"vendor"`
	Wxid		string 		`gorm:"index;type:varchar(50);" json:"wxid"`
	Role 		string		`gorm:"varchar(50);" json:"role"`
	Msg 		string		`gorm:"text" json:"msg"`
	CreatedAt	time.Time	`gorm:"default:CURRENT_TIMESTAMP()" json:"created_at"`
}