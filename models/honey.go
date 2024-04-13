package models

import (
	"time"
)

type Honey struct{
	ID 			uint 		`gorm:"primary_key;autoIncrement:true;" json:"id"`
	Style		string		`gorm:"type:varchar(50);" json:"style"`
	Prompt 		string		`gorm:"type:text;" json:"prompt"`
	Commit 		string		`gorm:"type:text;" json:"commit"`
	CreatedAt	time.Time	`gorm:"default:CURRENT_TIMESTAMP()"  json:"created_at"`
}