package model

import (
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

type Product struct {
	gorm.Model
	Name          string          `gorm:"not null"`
	Slug          string          `gorm:"not null;uniqueIndex"`
	Description   string
	Price         decimal.Decimal `gorm:"type:numeric(10,2);not null"`
	DiscountPrice  *decimal.Decimal `gorm:"type:numeric(10,2)"`
	Stock         int             `gorm:"not null;default:0"`
	CategoryID    uint
	Category      Category        `gorm:"foreignKey:CategoryID"`
	Images        []ProductImage  `gorm:"foreignKey:ProductID"`
	Variants      []ProductVariant `gorm:"foreignKey:ProductID"`
}

type ProductImage struct {
	gorm.Model
	ProductID uint   `gorm:"not null;index"`
	URL       string `gorm:"not null"`
	IsMain    bool   `gorm:"not null;default:false"`
}

type ProductVariant struct {
	gorm.Model
	ProductID   uint            `gorm:"not null;index"`
	WeightGrams int             `gorm:"not null"`
	Price       decimal.Decimal `gorm:"type:numeric(10,2);not null"`
	Stock       int             `gorm:"not null;default:0"`
}
