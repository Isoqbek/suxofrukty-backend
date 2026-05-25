package main

import (
	"fmt"

	"github.com/Isoqbek/suxofrukty-backend/internal/config"
	"github.com/Isoqbek/suxofrukty-backend/internal/handler"
	"github.com/Isoqbek/suxofrukty-backend/internal/middleware"
	"github.com/Isoqbek/suxofrukty-backend/internal/model"
	"github.com/Isoqbek/suxofrukty-backend/internal/repository"
	"github.com/Isoqbek/suxofrukty-backend/pkg/database"
	"github.com/Isoqbek/suxofrukty-backend/pkg/logger"
	"github.com/gin-gonic/gin"
)

func main() {
	log := logger.New()

	cfg, err := config.Load()
	if err != nil {
		log.Fatal().Err(err).Msg("failed to load config")
	}

	db, err := database.New(cfg.Database.DSN)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to connect to database")
	}

	if err := db.AutoMigrate(
		&model.Category{},
		&model.Product{},
		&model.ProductImage{},
		&model.ProductVariant{},
		&model.Order{},
		&model.OrderItem{},
		&model.Coupon{},
	); err != nil {
		log.Fatal().Err(err).Msg("failed to migrate database")
	}

	productRepo := repository.NewProductRepository(db)
	categoryRepo := repository.NewCategoryRepository(db)
	orderRepo := repository.NewOrderRepository(db)
	couponRepo := repository.NewCouponRepository(db)

	productH := handler.NewProductHandler(productRepo)
	categoryH := handler.NewCategoryHandler(categoryRepo)
	orderH := handler.NewOrderHandler(orderRepo, productRepo)
	couponH := handler.NewCouponHandler(couponRepo)

	r := gin.Default()
	r.Use(middleware.CORS())

	v1 := r.Group("/api/v1")
	{
		v1.GET("/products/", productH.List)
		v1.GET("/products/:slug/", productH.Get)
		v1.GET("/categories/", categoryH.List)
		v1.POST("/orders/", orderH.Create)
		v1.GET("/orders/:id/", orderH.Get)
		v1.POST("/coupons/validate/", couponH.Validate)
	}

	addr := fmt.Sprintf(":%s", cfg.Server.Port)
	log.Info().Str("addr", addr).Msg("server starting")
	if err := r.Run(addr); err != nil {
		log.Fatal().Err(err).Msg("server error")
	}
}
