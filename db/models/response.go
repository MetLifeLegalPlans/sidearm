package models

import (
	"gorm.io/gorm"
)

type Response struct {
	gorm.Model
	StatusCode int    `gorm:"index"`
	Method     string `gorm:"index"`
	URL        string `gorm:"index"`
	Duration   int64  `gorm:"index"`
}
