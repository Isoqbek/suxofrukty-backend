package handler

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/Isoqbek/suxofrukty-backend/internal/dto"
	"github.com/Isoqbek/suxofrukty-backend/internal/model"
	"github.com/Isoqbek/suxofrukty-backend/internal/repository"
	"github.com/shopspring/decimal"

	"github.com/gin-gonic/gin"
)

type OrderHandler struct {
	repo        *repository.OrderRepository
	productRepo *repository.ProductRepository
}

func NewOrderHandler(repo *repository.OrderRepository, productRepo *repository.ProductRepository) *OrderHandler {
	return &OrderHandler{repo: repo, productRepo: productRepo}
}

func (h *OrderHandler) Create(c *gin.Context) {
	var req dto.CreateOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var total decimal.Decimal
	items := make([]model.OrderItem, 0, len(req.Items))

	for _, it := range req.Items {
		p, err := h.productRepo.GetByID(it.ProductID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "product not found"})
			return
		}

		unitPrice, err := resolvePrice(p, it.VariantID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		items = append(items, model.OrderItem{
			ProductID: it.ProductID,
			VariantID: it.VariantID,
			Quantity:  it.Quantity,
			UnitPrice: unitPrice,
		})
		total = total.Add(unitPrice.Mul(decimal.NewFromInt(int64(it.Quantity))))
	}

	order := &model.Order{
		ContactName:    req.Contact.Name,
		ContactPhone:   req.Contact.Phone,
		ContactEmail:   req.Contact.Email,
		DeliveryCity:   req.Delivery.City,
		DeliveryBranch: req.Delivery.Branch,
		TotalPrice:     total,
		Status:         model.StatusPending,
		Items:          items,
	}

	if err := h.repo.Create(order); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, mapOrder(*order))
}

func (h *OrderHandler) Get(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	order, err := h.repo.GetByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		return
	}

	c.JSON(http.StatusOK, mapOrder(*order))
}

func resolvePrice(p *model.Product, variantID *uint) (decimal.Decimal, error) {
	if variantID != nil {
		for _, v := range p.Variants {
			if v.ID == *variantID {
				return v.Price, nil
			}
		}
		return decimal.Zero, errors.New("variant not found")
	}
	if p.DiscountPrice != nil {
		return *p.DiscountPrice, nil
	}
	return p.Price, nil
}

func mapOrder(o model.Order) dto.OrderResponse {
	items := make([]dto.OrderItemResponse, len(o.Items))
	for i, it := range o.Items {
		items[i] = dto.OrderItemResponse{
			ID:        it.ID,
			ProductID: it.ProductID,
			VariantID: it.VariantID,
			Quantity:  it.Quantity,
			UnitPrice: it.UnitPrice,
		}
	}
	return dto.OrderResponse{
		ID:             o.ID,
		Status:         string(o.Status),
		ContactName:    o.ContactName,
		ContactPhone:   o.ContactPhone,
		ContactEmail:   o.ContactEmail,
		DeliveryCity:   o.DeliveryCity,
		DeliveryBranch: o.DeliveryBranch,
		TotalPrice:     o.TotalPrice,
		Discount:       o.Discount,
		Items:          items,
	}
}
