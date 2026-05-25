package handler

import (
	"net/http"
	"strconv"

	"github.com/Isoqbek/suxofrukty-backend/internal/dto"
	"github.com/Isoqbek/suxofrukty-backend/internal/model"
	"github.com/Isoqbek/suxofrukty-backend/internal/repository"
	"github.com/gin-gonic/gin"
)

type ProductHandler struct {
	repo *repository.ProductRepository
}

func NewProductHandler(repo *repository.ProductRepository) *ProductHandler {
	return &ProductHandler{repo: repo}
}

func (h *ProductHandler) List(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	f := repository.ProductFilter{
		Category: c.Query("category"),
		Search:   c.Query("search"),
		Sale:     c.Query("sale") == "true",
		Page:     page,
		PageSize: pageSize,
	}

	products, total, err := h.repo.List(f)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, dto.PaginatedProductsResponse{
		Count:   total,
		Results: mapProducts(products),
	})
}

func (h *ProductHandler) Get(c *gin.Context) {
	slug := c.Param("slug")
	p, err := h.repo.GetBySlug(slug)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		return
	}
	c.JSON(http.StatusOK, mapProduct(*p))
}

func mapProducts(ps []model.Product) []dto.ProductResponse {
	out := make([]dto.ProductResponse, len(ps))
	for i, p := range ps {
		out[i] = mapProduct(p)
	}
	return out
}

func mapProduct(p model.Product) dto.ProductResponse {
	images := make([]dto.ProductImageResponse, len(p.Images))
	for i, img := range p.Images {
		images[i] = dto.ProductImageResponse{ID: img.ID, URL: img.URL, IsMain: img.IsMain}
	}

	variants := make([]dto.ProductVariantResponse, len(p.Variants))
	for i, v := range p.Variants {
		variants[i] = dto.ProductVariantResponse{
			ID: v.ID, WeightGrams: v.WeightGrams, Price: v.Price, Stock: v.Stock,
		}
	}

	return dto.ProductResponse{
		ID:            p.ID,
		Name:          p.Name,
		Slug:          p.Slug,
		Description:   p.Description,
		Price:         p.Price,
		DiscountPrice:  p.DiscountPrice,
		Stock:         p.Stock,
		Category:      dto.CategoryResponse{ID: p.Category.ID, Name: p.Category.Name, Slug: p.Category.Slug},
		Images:        images,
		Variants:      variants,
	}
}
