package handler

import (
	"net/http"
	"strconv"

	"github.com/Isoqbek/suxofrukty-backend/internal/dto"
	"github.com/Isoqbek/suxofrukty-backend/internal/model"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type AdminOrderHandler struct {
	db *gorm.DB
}

func NewAdminOrderHandler(db *gorm.DB) *AdminOrderHandler {
	return &AdminOrderHandler{db: db}
}

func (h *AdminOrderHandler) List(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	status := c.Query("status")

	q := h.db.Model(&model.Order{}).Preload("Items")
	if status != "" {
		q = q.Where("status = ?", status)
	}

	var total int64
	q.Count(&total)

	var orders []model.Order
	q.Order("created_at DESC").Offset((page - 1) * pageSize).Limit(pageSize).Find(&orders)

	results := make([]dto.OrderResponse, len(orders))
	for i, o := range orders {
		results[i] = mapOrder(o)
	}

	c.JSON(http.StatusOK, gin.H{"count": total, "results": results})
}

func (h *AdminOrderHandler) Get(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	var o model.Order
	if err := h.db.Preload("Items").First(&o, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		return
	}

	c.JSON(http.StatusOK, mapOrder(o))
}

func (h *AdminOrderHandler) UpdateStatus(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	var body struct {
		Status string `json:"status" binding:"required"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	validStatuses := map[string]bool{
		"pending": true, "paid": true, "processing": true,
		"shipped": true, "delivered": true, "cancelled": true,
	}
	if !validStatuses[body.Status] {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid status"})
		return
	}

	h.db.Model(&model.Order{}).Where("id = ?", id).Update("status", body.Status)
	c.JSON(http.StatusOK, gin.H{"ok": true})
}
