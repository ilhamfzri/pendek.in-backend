package domain

import "gorm.io/gorm"

type CustomLink struct {
	gorm.Model
	UserID                string `gorm:"index"`
	Title                 string
	ShortLinkCode         string `gorm:"unique"`
	LongLink              string
	ShowOnProfile         bool
	Activate              bool
	CustomThumbnailID     *uint
	CustomThumbnail       CustomThumbnail `gorm:"foreignKey:CustomThumbnailID"`
	ThumbnailID           *uint
	Thumbnail             Thumbnail `gorm:"foreignKey:ThumbnailID"`
	CustomLinkAnalytic    []CustomLinkAnalytic
	CustomLinkInteraction []CustomLinkInteraction
}
