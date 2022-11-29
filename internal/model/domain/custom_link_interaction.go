package domain

import (
	"gorm.io/gorm"
)

type CustomLinkInteraction struct {
	gorm.Model
	CustomLinkID    string
	ClientIP        string
	Username        string
	SocialMediaName string
}
