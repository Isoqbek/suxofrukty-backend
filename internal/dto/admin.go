package dto

type AdminLoginRequest struct {
	Email    string `json:"email"    binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type AdminLoginResponse struct {
	Token string `json:"token"`
}

type AdminStatsResponse struct {
	OrdersToday   int64   `json:"orders_today"`
	RevenueTotal  float64 `json:"revenue_total"`
	ProductsCount int64   `json:"products_count"`
	LowStock      int64   `json:"low_stock"`
}
