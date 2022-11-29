package domain

import "gorm.io/gorm"

type CustomLink struct {
	gorm.Model
	UserID        string `gorm:"index"`
	Title         string
	ShortLinkID   string `gorm:"unique;index;<-:create"`
	LongLink      string
	ShowOnProfile bool
	Activate      bool
}
