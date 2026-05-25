package dto

import "github.com/shopspring/decimal"

type ValidateCouponRequest struct {
	Code string `json:"code" binding:"required"`
}

type ValidateCouponResponse struct {
	Discount decimal.Decimal `json:"discount"`
	Type     string          `json:"type"`
}
