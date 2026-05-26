package handler

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/Isoqbek/suxofrukty-backend/internal/dto"
	"github.com/Isoqbek/suxofrukty-backend/internal/model"
	"github.com/gin-gonic/gin"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

type AdminProductHandler struct {
	db *gorm.DB
}

func NewAdminProductHandler(db *gorm.DB) *AdminProductHandler {
	return &AdminProductHandler{db: db}
}

type adminProductRequest struct {
	Name          string                   `json:"name"           binding:"required"`
	Slug          string                   `json:"slug"`
	Description   string                   `json:"description"`
	Price         decimal.Decimal          `json:"price"          binding:"required"`
	DiscountPrice  *decimal.Decimal         `json:"discount_price"`
	Stock         int                      `json:"stock"`
	CategoryID    uint                     `json:"category_id"    binding:"required"`
	Images        []dto.ProductImageResponse `json:"images"`
	Variants      []adminVariantRequest     `json:"variants"`
}

type adminVariantRequest struct {
	ID          uint            `json:"id"`
	WeightGrams int             `json:"weight_grams" binding:"required"`
	Price       decimal.Decimal `json:"price"        binding:"required"`
	Stock       int             `json:"stock"`
}

func (h *AdminProductHandler) List(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	search := c.Query("search")

	q := h.db.Model(&model.Product{}).Preload("Images").Preload("Variants").Preload("Category")
	if search != "" {
		q = q.Where("name ILIKE ?", "%"+search+"%")
	}

	var total int64
	q.Count(&total)

	var products []model.Product
	q.Offset((page - 1) * pageSize).Limit(pageSize).Find(&products)

	c.JSON(http.StatusOK, dto.PaginatedProductsResponse{
		Count:   total,
		Results: mapProducts(products),
	})
}

func (h *AdminProductHandler) Create(c *gin.Context) {
	var req adminProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	slug := req.Slug
	if slug == "" {
		slug = slugify(req.Name)
	}

	p := model.Product{
		Name:         req.Name,
		Slug:         slug,
		Description:  req.Description,
		Price:        req.Price,
		DiscountPrice: req.DiscountPrice,
		Stock:        req.Stock,
		CategoryID:   req.CategoryID,
	}

	if err := h.db.Create(&p).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	h.syncImages(p.ID, req.Images)
	h.syncVariants(p.ID, req.Variants)

	h.db.Preload("Images").Preload("Variants").Preload("Category").First(&p, p.ID)
	c.JSON(http.StatusCreated, mapProduct(p))
}

func (h *AdminProductHandler) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	var p model.Product
	if err := h.db.First(&p, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		return
	}

	var req adminProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	slug := req.Slug
	if slug == "" {
		slug = slugify(req.Name)
	}

	p.Name = req.Name
	p.Slug = slug
	p.Description = req.Description
	p.Price = req.Price
	p.DiscountPrice = req.DiscountPrice
	p.Stock = req.Stock
	p.CategoryID = req.CategoryID

	h.db.Save(&p)
	h.syncImages(p.ID, req.Images)
	h.syncVariants(p.ID, req.Variants)

	h.db.Preload("Images").Preload("Variants").Preload("Category").First(&p, p.ID)
	c.JSON(http.StatusOK, mapProduct(p))
}

func (h *AdminProductHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	h.db.Delete(&model.Product{}, id)
	c.JSON(http.StatusOK, gin.H{"ok": true})
}

func (h *AdminProductHandler) syncImages(productID uint, images []dto.ProductImageResponse) {
	h.db.Where("product_id = ?", productID).Delete(&model.ProductImage{})
	for _, img := range images {
		h.db.Create(&model.ProductImage{
			ProductID: productID,
			URL:       img.URL,
			IsMain:    img.IsMain,
		})
	}
}

func (h *AdminProductHandler) syncVariants(productID uint, variants []adminVariantRequest) {
	existingIDs := map[uint]bool{}
	for _, v := range variants {
		if v.ID != 0 {
			existingIDs[v.ID] = true
			h.db.Model(&model.ProductVariant{}).Where("id = ?", v.ID).Updates(map[string]any{
				"weight_grams": v.WeightGrams,
				"price":        v.Price,
				"stock":        v.Stock,
			})
		} else {
			h.db.Create(&model.ProductVariant{
				ProductID:   productID,
				WeightGrams: v.WeightGrams,
				Price:       v.Price,
				Stock:       v.Stock,
			})
		}
	}
}

func slugify(s string) string {
	s = strings.ToLower(s)
	s = strings.ReplaceAll(s, " ", "-")
	return s
}
