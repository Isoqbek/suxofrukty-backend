package repository

import (
	"github.com/Isoqbek/suxofrukty-backend/internal/model"
	"gorm.io/gorm"
)

type ProductFilter struct {
	Category string
	Search   string
	Sale     bool
	Page     int
	PageSize int
}

type ProductRepository struct {
	db *gorm.DB
}

func NewProductRepository(db *gorm.DB) *ProductRepository {
	return &ProductRepository{db: db}
}

func (r *ProductRepository) List(f ProductFilter) ([]model.Product, int64, error) {
	q := r.db.Model(&model.Product{}).
		Preload("Images").
		Preload("Variants").
		Preload("Category")

	if f.Category != "" {
		q = q.Joins("JOIN categories ON categories.id = products.category_id").
			Where("categories.slug = ?", f.Category)
	}
	if f.Search != "" {
		q = q.Where("products.name ILIKE ?", "%"+f.Search+"%")
	}
	if f.Sale {
		q = q.Where("products.discount_price IS NOT NULL")
	}

	var total int64
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if f.PageSize <= 0 {
		f.PageSize = 20
	}
	if f.Page <= 0 {
		f.Page = 1
	}

	var products []model.Product
	err := q.Offset((f.Page - 1) * f.PageSize).Limit(f.PageSize).Find(&products).Error
	return products, total, err
}

func (r *ProductRepository) GetBySlug(slug string) (*model.Product, error) {
	var p model.Product
	err := r.db.Preload("Images").Preload("Variants").Preload("Category").
		Where("slug = ?", slug).First(&p).Error
	if err != nil {
		return nil, err
	}
	return &p, nil
}

func (r *ProductRepository) GetByID(id uint) (*model.Product, error) {
	var p model.Product
	err := r.db.Preload("Variants").First(&p, id).Error
	if err != nil {
		return nil, err
	}
	return &p, nil
}
