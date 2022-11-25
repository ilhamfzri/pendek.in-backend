package domain

import (
	"gorm.io/gorm"
)

type SocialMediaLink struct {
	gorm.Model
	TypeID          uint
	SocialMediaType SocialMediaType `gorm:"foreignKey:TypeID"`
	UserID          string
	LinkOrUsername  string
	Activate        bool
}
