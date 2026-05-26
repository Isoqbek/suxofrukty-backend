package main

import (
	"fmt"
	"os"

	"github.com/Isoqbek/suxofrukty-backend/internal/config"
	"github.com/Isoqbek/suxofrukty-backend/internal/handler"
	"github.com/Isoqbek/suxofrukty-backend/internal/middleware"
	"github.com/Isoqbek/suxofrukty-backend/internal/model"
	"github.com/Isoqbek/suxofrukty-backend/internal/repository"
	"github.com/Isoqbek/suxofrukty-backend/pkg/database"
	"github.com/Isoqbek/suxofrukty-backend/pkg/logger"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
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
		&model.Admin{},
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

	adminRepo := repository.NewAdminRepository(db)
	seedAdmin(adminRepo, cfg)

	productRepo := repository.NewProductRepository(db)
	categoryRepo := repository.NewCategoryRepository(db)
	orderRepo := repository.NewOrderRepository(db)
	couponRepo := repository.NewCouponRepository(db)

	productH := handler.NewProductHandler(productRepo)
	categoryH := handler.NewCategoryHandler(categoryRepo)
	orderH := handler.NewOrderHandler(orderRepo, productRepo)
	couponH := handler.NewCouponHandler(couponRepo)
	adminH := handler.NewAdminHandler(adminRepo, productRepo, orderRepo, db, cfg.JWT.Secret)
	adminProductH := handler.NewAdminProductHandler(db)
	adminOrderH := handler.NewAdminOrderHandler(db)
	adminCategoryH := handler.NewAdminCategoryHandler(db)
	adminCouponH := handler.NewAdminCouponHandler(db)

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

		admin := v1.Group("/admin")
		admin.POST("/auth/login", adminH.Login)

		protected := admin.Group("")
		protected.Use(middleware.AdminOnly(cfg.JWT.Secret))
		{
			protected.GET("/stats", adminH.Stats)

			protected.GET("/products/", adminProductH.List)
			protected.POST("/products/", adminProductH.Create)
			protected.PUT("/products/:id", adminProductH.Update)
			protected.DELETE("/products/:id", adminProductH.Delete)

			protected.GET("/orders/", adminOrderH.List)
			protected.GET("/orders/:id", adminOrderH.Get)
			protected.PATCH("/orders/:id/status", adminOrderH.UpdateStatus)

			protected.GET("/categories/", adminCategoryH.List)
			protected.POST("/categories/", adminCategoryH.Create)
			protected.PUT("/categories/:id", adminCategoryH.Update)
			protected.DELETE("/categories/:id", adminCategoryH.Delete)

			protected.GET("/coupons/", adminCouponH.List)
			protected.POST("/coupons/", adminCouponH.Create)
			protected.PATCH("/coupons/:id/toggle", adminCouponH.Toggle)
		}
	}

	addr := fmt.Sprintf(":%s", cfg.Server.Port)
	log.Info().Str("addr", addr).Msg("server starting")
	if err := r.Run(addr); err != nil {
		log.Fatal().Err(err).Msg("server error")
	}
}

func seedAdmin(repo *repository.AdminRepository, cfg *config.Config) {
	exists, _ := repo.Exists()
	if exists {
		return
	}

	email := os.Getenv("ADMIN_EMAIL")
	password := os.Getenv("ADMIN_PASSWORD")
	if email == "" || password == "" {
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return
	}
	repo.Create(email, string(hash))
}
