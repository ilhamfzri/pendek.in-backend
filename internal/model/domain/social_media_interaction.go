package domain

import (
	"gorm.io/gorm"
)

type SocialMediaInteraction struct {
	gorm.Model
	SocialMediaLinkID uint `gorm:"index"`
	ClientIP          string
	UserAgent         string
}
