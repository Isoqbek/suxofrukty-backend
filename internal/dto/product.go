package dto

import "github.com/shopspring/decimal"

type ProductImageResponse struct {
	ID     uint   `json:"id"`
	URL    string `json:"url"`
	IsMain bool   `json:"is_main"`
}

type ProductVariantResponse struct {
	ID          uint            `json:"id"`
	WeightGrams int             `json:"weight_grams"`
	Price       decimal.Decimal `json:"price"`
	Stock       int             `json:"stock"`
}

type ProductResponse struct {
	ID            uint                     `json:"id"`
	Name          string                   `json:"name"`
	Slug          string                   `json:"slug"`
	Description   string                   `json:"description"`
	Price         decimal.Decimal          `json:"price"`
	DiscountPrice  *decimal.Decimal         `json:"discount_price"`
	Stock         int                      `json:"stock"`
	Category      CategoryResponse         `json:"category"`
	Images        []ProductImageResponse   `json:"images"`
	Variants      []ProductVariantResponse `json:"variants"`
}

type PaginatedProductsResponse struct {
	Count   int64             `json:"count"`
	Results []ProductResponse `json:"results"`
}
