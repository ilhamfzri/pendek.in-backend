package domain

import "gorm.io/gorm"

type Thumbnail struct {
	gorm.Model
	Name    string
	IconUrl string
}
