package models

import (
	"gorm.io/gorm"
)

type EmailLog struct {
	gorm.Model
	Subject string `gorm:"type:varchar(255);index;not null"`
	From    string `gorm:"type:varchar(255);index;not null"`
	To      string `gorm:"type:varchar(255);index;not null"`
	Body    string `gorm:"type:text"`
}
