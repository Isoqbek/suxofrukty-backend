package model

import (
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

type CouponType string

const (
	CouponPercent     CouponType = "percent"
	CouponFixed       CouponType = "fixed"
	CouponFreeDelivery CouponType = "free_delivery"
)

type Coupon struct {
	gorm.Model
	Code      string          `gorm:"not null;uniqueIndex"`
	Type      CouponType      `gorm:"type:varchar(20);not null"`
	Discount  decimal.Decimal `gorm:"type:numeric(10,2);not null;default:0"`
	UsedCount int             `gorm:"not null;default:0"`
	MaxUses   int             `gorm:"not null;default:0"`
	Active    bool            `gorm:"not null;default:true"`
}
