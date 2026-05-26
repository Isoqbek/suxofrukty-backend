package handler

import (
	"net/http"
	"time"

	"github.com/Isoqbek/suxofrukty-backend/internal/dto"
	"github.com/Isoqbek/suxofrukty-backend/internal/model"
	"github.com/Isoqbek/suxofrukty-backend/internal/repository"
	"github.com/Isoqbek/suxofrukty-backend/pkg/jwtutil"
	"github.com/gin-gonic/gin"
	"github.com/shopspring/decimal"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type AdminHandler struct {
	repo        *repository.AdminRepository
	productRepo *repository.ProductRepository
	orderRepo   *repository.OrderRepository
	db          *gorm.DB
	jwtSecret   string
}

func NewAdminHandler(
	repo *repository.AdminRepository,
	productRepo *repository.ProductRepository,
	orderRepo *repository.OrderRepository,
	db *gorm.DB,
	jwtSecret string,
) *AdminHandler {
	return &AdminHandler{repo: repo, productRepo: productRepo, orderRepo: orderRepo, db: db, jwtSecret: jwtSecret}
}

func (h *AdminHandler) Login(c *gin.Context) {
	var req dto.AdminLoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	admin, err := h.repo.GetByEmail(req.Email)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(admin.PasswordHash), []byte(req.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}

	token, err := jwtutil.Generate(admin.ID, h.jwtSecret)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "token generation failed"})
		return
	}

	c.JSON(http.StatusOK, dto.AdminLoginResponse{Token: token})
}

func (h *AdminHandler) Stats(c *gin.Context) {
	var ordersToday int64
	today := time.Now().Truncate(24 * time.Hour)
	h.db.Model(&model.Order{}).Where("created_at >= ?", today).Count(&ordersToday)

	var revenueRow struct{ Total decimal.Decimal }
	h.db.Model(&model.Order{}).
		Where("status != ?", model.StatusCancelled).
		Select("COALESCE(SUM(total_price), 0) as total").
		Scan(&revenueRow)

	var productsCount int64
	h.db.Model(&model.Product{}).Count(&productsCount)

	var lowStock int64
	h.db.Model(&model.Product{}).Where("stock > 0 AND stock <= 5").Count(&lowStock)

	total, _ := revenueRow.Total.Float64()
	c.JSON(http.StatusOK, dto.AdminStatsResponse{
		OrdersToday:   ordersToday,
		RevenueTotal:  total,
		ProductsCount: productsCount,
		LowStock:      lowStock,
	})
}
