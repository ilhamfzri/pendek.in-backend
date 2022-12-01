package domain

import "gorm.io/gorm"

type CustomThumbnail struct {
	gorm.Model
	UserID  string
	ImageID string
}
