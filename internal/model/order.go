package model

import (
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

type OrderStatus string

const (
	StatusPending    OrderStatus = "pending"
	StatusPaid       OrderStatus = "paid"
	StatusProcessing OrderStatus = "processing"
	StatusShipped    OrderStatus = "shipped"
	StatusDelivered  OrderStatus = "delivered"
	StatusCancelled  OrderStatus = "cancelled"
)

type Order struct {
	gorm.Model
	ContactName    string          `gorm:"not null"`
	ContactPhone   string          `gorm:"not null"`
	ContactEmail   string          `gorm:"not null"`
	DeliveryCity   string          `gorm:"not null"`
	DeliveryBranch string          `gorm:"not null"`
	TotalPrice     decimal.Decimal `gorm:"type:numeric(10,2);not null"`
	Status         OrderStatus     `gorm:"type:varchar(20);not null;default:'pending'"`
	CouponCode     string
	Discount       decimal.Decimal `gorm:"type:numeric(10,2);default:0"`
	Items          []OrderItem     `gorm:"foreignKey:OrderID"`
}

type OrderItem struct {
	gorm.Model
	OrderID   uint            `gorm:"not null;index"`
	ProductID uint            `gorm:"not null"`
	VariantID *uint
	Quantity  int             `gorm:"not null"`
	UnitPrice decimal.Decimal `gorm:"type:numeric(10,2);not null"`
}
