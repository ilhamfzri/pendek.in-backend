package domain

import "time"

type SocialMediaType struct {
	ID        uint   `gorm:"primaryKey"`
	Name      string `gorm:"unique;index;<-:create"`
	Example   string
	IconUrl   string
	CreatedAt time.Time
	UpdatedAt time.Time
}
