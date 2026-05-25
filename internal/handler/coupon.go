package handler

import (
	"net/http"

	"github.com/Isoqbek/suxofrukty-backend/internal/dto"
	"github.com/Isoqbek/suxofrukty-backend/internal/repository"
	"github.com/gin-gonic/gin"
)

type CouponHandler struct {
	repo *repository.CouponRepository
}

func NewCouponHandler(repo *repository.CouponRepository) *CouponHandler {
	return &CouponHandler{repo: repo}
}

func (h *CouponHandler) Validate(c *gin.Context) {
	var req dto.ValidateCouponRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	coupon, err := h.repo.GetByCode(req.Code)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "coupon not found or inactive"})
		return
	}

	if coupon.MaxUses > 0 && coupon.UsedCount >= coupon.MaxUses {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": "coupon usage limit reached"})
		return
	}

	c.JSON(http.StatusOK, dto.ValidateCouponResponse{
		Discount: coupon.Discount,
		Type:     string(coupon.Type),
	})
}
