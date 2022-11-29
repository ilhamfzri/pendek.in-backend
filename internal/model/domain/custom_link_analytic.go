package domain

import (
	"time"

	"gorm.io/gorm"
)

type CustomLinkAnalytic struct {
	gorm.Model
	TotalClick       int
	CustomLinkID     uint
	DeviceAnalyticID uint
	DeviceAnalytic   DeviceAnalytic `gorm:"foreignKey:DeviceAnalyticID"`
	Date             time.Time
}
