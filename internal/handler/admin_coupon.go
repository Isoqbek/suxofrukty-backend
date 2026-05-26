package handler

import (
	"net/http"
	"strconv"

	"github.com/Isoqbek/suxofrukty-backend/internal/model"
	"github.com/gin-gonic/gin"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

type AdminCouponHandler struct {
	db *gorm.DB
}

func NewAdminCouponHandler(db *gorm.DB) *AdminCouponHandler {
	return &AdminCouponHandler{db: db}
}

type couponRequest struct {
	Code     string          `json:"code"     binding:"required"`
	Type     string          `json:"type"     binding:"required"`
	Discount decimal.Decimal `json:"discount"`
	MaxUses  int             `json:"max_uses"`
}

func (h *AdminCouponHandler) List(c *gin.Context) {
	var coupons []model.Coupon
	h.db.Find(&coupons)
	c.JSON(http.StatusOK, coupons)
}

func (h *AdminCouponHandler) Create(c *gin.Context) {
	var req couponRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	coupon := model.Coupon{
		Code:     req.Code,
		Type:     model.CouponType(req.Type),
		Discount: req.Discount,
		MaxUses:  req.MaxUses,
		Active:   true,
	}
	if err := h.db.Create(&coupon).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, coupon)
}

func (h *AdminCouponHandler) Toggle(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	var coupon model.Coupon
	if err := h.db.First(&coupon, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		return
	}
	h.db.Model(&coupon).Update("active", !coupon.Active)
	c.JSON(http.StatusOK, gin.H{"active": !coupon.Active})
}
