package repository

import (
	"github.com/Isoqbek/suxofrukty-backend/internal/model"
	"gorm.io/gorm"
)

type CouponRepository struct {
	db *gorm.DB
}

func NewCouponRepository(db *gorm.DB) *CouponRepository {
	return &CouponRepository{db: db}
}

func (r *CouponRepository) GetByCode(code string) (*model.Coupon, error) {
	var c model.Coupon
	err := r.db.Where("code = ? AND active = true", code).First(&c).Error
	if err != nil {
		return nil, err
	}
	return &c, nil
}
