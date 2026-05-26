package repository

import (
	"github.com/Isoqbek/suxofrukty-backend/internal/model"
	"gorm.io/gorm"
)

type AdminRepository struct {
	db *gorm.DB
}

func NewAdminRepository(db *gorm.DB) *AdminRepository {
	return &AdminRepository{db: db}
}

func (r *AdminRepository) GetByEmail(email string) (*model.Admin, error) {
	var a model.Admin
	err := r.db.Where("email = ?", email).First(&a).Error
	if err != nil {
		return nil, err
	}
	return &a, nil
}

func (r *AdminRepository) Create(email, passwordHash string) (*model.Admin, error) {
	a := &model.Admin{Email: email, PasswordHash: passwordHash}
	err := r.db.Create(a).Error
	return a, err
}

func (r *AdminRepository) Exists() (bool, error) {
	var count int64
	err := r.db.Model(&model.Admin{}).Count(&count).Error
	return count > 0, err
}
