package models

import (
	"time"
	// "github.com/google/uuid"
)

type Promotion struct {
	ID			string 	    `gorm:"primary_key;" json:"id"`
	Publisher	string		`gorm:"index;not null;type:varchar(50);" json:"punlisher"`
	Consumer 	string 		`gorm:"default:null;type:text;" json:"consumer"` 
	CreatedAt	time.Time	`gorm:"default:CURRENT_TIMESTAMP()" json:"created_at"`
	UpdatedAt	time.Time	`gorm:"default:CURRENT_TIMESTAMP()" json:"updated_at"`		
}