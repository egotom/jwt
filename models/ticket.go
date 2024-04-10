package models
import (
	"time"
)

type Ticket struct{
	ID 			uint 		`gorm:"primary_key;autoIncrement:true;" json:"id"`
	Amount		int64
	Wxid		string 		`gorm:"index;type:varchar(50);" json:"wxid"`
	Commit 		string		`gorm:"type:text;" json:"commit"`
	CreatedAt	time.Time	`gorm:"default:CURRENT_TIMESTAMP()"  json:"created_at"`
}