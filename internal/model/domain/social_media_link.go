package domain

import (
	"gorm.io/gorm"
)

type SocialMediaLink struct {
	gorm.Model
	TypeID                 uint
	SocialMediaType        SocialMediaType `gorm:"foreignKey:TypeID"`
	SocialMediaAnalytic    []SocialMediaAnalytic
	SocialMediaInteraction []SocialMediaInteraction
	UserID                 string
	LinkOrUsername         string
	Activate               bool
}
