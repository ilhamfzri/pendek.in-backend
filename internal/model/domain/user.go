package domain

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID                string `gorm:"type:uuid;default:gen_random_uuid()"`
	Username          string `gorm:"unique;index"`
	FullName          string
	Bio               string
	Email             string `gorm:"unique;index;<-:create"`
	Password          string
	Verified          bool
	ResetPasswordCode string
	VerificationCode  string
	ProfilePic        string
	SocialMediaLinks  []SocialMediaLink `gorm:"foreignKey:UserID"`
	LastLogin         time.Time
	CreatedAt         time.Time
	UpdatedAt         time.Time
	DeletedAt         gorm.DeletedAt `gorm:"index"`
}
