package models

import (
	"gorm.io/gorm"
)

type Response struct {
	gorm.Model
	StatusCode int
	Method     string
	URL        string
	Duration   int64
}
