package domain

import (
	"time"

	"gorm.io/gorm"
)

type CustomLinkAnalytic struct {
	gorm.Model
	ClickCount       int
	ViewCount        int
	CustomLinkID     uint
	DeviceAnalyticID uint
	DeviceAnalytic   DeviceAnalytic `gorm:"foreignKey:DeviceAnalyticID"`
	Date             time.Time
}
