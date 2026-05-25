package repository

import (
	"github.com/Isoqbek/suxofrukty-backend/internal/model"
	"gorm.io/gorm"
)

type CategoryRepository struct {
	db *gorm.DB
}

func NewCategoryRepository(db *gorm.DB) *CategoryRepository {
	return &CategoryRepository{db: db}
}

func (r *CategoryRepository) List() ([]model.Category, error) {
	var cats []model.Category
	err := r.db.Find(&cats).Error
	return cats, err
}
