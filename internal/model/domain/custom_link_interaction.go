package domain

import (
	"gorm.io/gorm"
)

type CustomLinkInteraction struct {
	gorm.Model
	CustomLinkID uint `gorm:"index"`
	ClientIP     string
	UserAgent    string
}
