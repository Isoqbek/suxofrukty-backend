package handler

import (
	"net/http"

	"github.com/Isoqbek/suxofrukty-backend/internal/dto"
	"github.com/Isoqbek/suxofrukty-backend/internal/repository"
	"github.com/gin-gonic/gin"
)

type CategoryHandler struct {
	repo *repository.CategoryRepository
}

func NewCategoryHandler(repo *repository.CategoryRepository) *CategoryHandler {
	return &CategoryHandler{repo: repo}
}

func (h *CategoryHandler) List(c *gin.Context) {
	cats, err := h.repo.List()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	out := make([]dto.CategoryResponse, len(cats))
	for i, cat := range cats {
		out[i] = dto.CategoryResponse{ID: cat.ID, Name: cat.Name, Slug: cat.Slug}
	}
	c.JSON(http.StatusOK, out)
}
