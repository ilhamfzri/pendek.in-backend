package domain

import "gorm.io/gorm"

type DeviceAnalytic struct {
	gorm.Model
	Mobile  int
	Tablet  int
	Desktop int
	Other   int
}
