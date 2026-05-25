package dto

import "github.com/shopspring/decimal"

type CreateOrderRequest struct {
	Contact  ContactDTO  `json:"contact"  binding:"required"`
	Delivery DeliveryDTO `json:"delivery" binding:"required"`
	Items    []OrderItemDTO `json:"items" binding:"required,min=1"`
}

type ContactDTO struct {
	Name  string `json:"name"  binding:"required"`
	Phone string `json:"phone" binding:"required"`
	Email string `json:"email" binding:"required,email"`
}

type DeliveryDTO struct {
	City   string `json:"city"   binding:"required"`
	Branch string `json:"branch" binding:"required"`
}

type OrderItemDTO struct {
	ProductID uint  `json:"product_id" binding:"required"`
	VariantID *uint `json:"variant_id"`
	Quantity  int   `json:"quantity"   binding:"required,min=1"`
}

type OrderItemResponse struct {
	ID        uint            `json:"id"`
	ProductID uint            `json:"product_id"`
	VariantID *uint           `json:"variant_id"`
	Quantity  int             `json:"quantity"`
	UnitPrice decimal.Decimal `json:"unit_price"`
}

type OrderResponse struct {
	ID             uint               `json:"id"`
	Status         string             `json:"status"`
	ContactName    string             `json:"contact_name"`
	ContactPhone   string             `json:"contact_phone"`
	ContactEmail   string             `json:"contact_email"`
	DeliveryCity   string             `json:"delivery_city"`
	DeliveryBranch string             `json:"delivery_branch"`
	TotalPrice     decimal.Decimal    `json:"total_price"`
	Discount       decimal.Decimal    `json:"discount"`
	Items          []OrderItemResponse `json:"items"`
}
