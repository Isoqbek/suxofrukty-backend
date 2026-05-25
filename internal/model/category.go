package model

import "gorm.io/gorm"

type Category struct {
	gorm.Model
	Name string `gorm:"not null;uniqueIndex"`
	Slug string `gorm:"not null;uniqueIndex"`
}
