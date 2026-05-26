package handler

import (
	"net/http"
	"strconv"

	"github.com/Isoqbek/suxofrukty-backend/internal/dto"
	"github.com/Isoqbek/suxofrukty-backend/internal/model"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type AdminCategoryHandler struct {
	db *gorm.DB
}

func NewAdminCategoryHandler(db *gorm.DB) *AdminCategoryHandler {
	return &AdminCategoryHandler{db: db}
}

type categoryRequest struct {
	Name string `json:"name" binding:"required"`
	Slug string `json:"slug"`
}

func (h *AdminCategoryHandler) List(c *gin.Context) {
	var cats []model.Category
	h.db.Find(&cats)
	out := make([]dto.CategoryResponse, len(cats))
	for i, cat := range cats {
		out[i] = dto.CategoryResponse{ID: cat.ID, Name: cat.Name, Slug: cat.Slug}
	}
	c.JSON(http.StatusOK, out)
}

func (h *AdminCategoryHandler) Create(c *gin.Context) {
	var req categoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	slug := req.Slug
	if slug == "" {
		slug = slugify(req.Name)
	}
	cat := model.Category{Name: req.Name, Slug: slug}
	if err := h.db.Create(&cat).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, dto.CategoryResponse{ID: cat.ID, Name: cat.Name, Slug: cat.Slug})
}

func (h *AdminCategoryHandler) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	var req categoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	slug := req.Slug
	if slug == "" {
		slug = slugify(req.Name)
	}
	h.db.Model(&model.Category{}).Where("id = ?", id).Updates(map[string]any{"name": req.Name, "slug": slug})
	c.JSON(http.StatusOK, gin.H{"ok": true})
}

func (h *AdminCategoryHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	h.db.Delete(&model.Category{}, id)
	c.JSON(http.StatusOK, gin.H{"ok": true})
}
