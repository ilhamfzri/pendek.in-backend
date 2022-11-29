package domain

import (
	"time"

	"gorm.io/gorm"
)

type SocialMediaAnalytic struct {
	gorm.Model
	ClickCount        int
	ViewCount         int
	SocialMediaLinkID uint
	DeviceAnalyticID  uint
	DeviceAnalytic    DeviceAnalytic `gorm:"foreignKey:DeviceAnalyticID"`
	Date              time.Time
}
