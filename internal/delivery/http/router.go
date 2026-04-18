package http

import (
	"github.com/gin-gonic/gin"
	"github.com/valenrio66/be-pos/config"
	"github.com/valenrio66/be-pos/internal/delivery/http/handler"
	"github.com/valenrio66/be-pos/internal/delivery/http/middleware"
)

func SetupRoutes(r *gin.Engine, cfg *config.Config, userHandler *handler.UserHandler, productHandler *handler.ProductHandler, transHandler *handler.TransactionHandler) {
	v1 := r.Group("/api/v1")
	{
		public := v1.Group("/users")
		{
			public.POST("/", userHandler.CreateUser)
			public.POST("/login", userHandler.Login)
			public.POST("/refresh", userHandler.RefreshToken)
		}

		protected := v1.Group("/users")
		protected.Use(middleware.AuthMiddleware(cfg.JWTSecret))
		{
			protected.GET("/me", userHandler.GetProfile)
		}

		staff := v1.Group("/products")
		staff.Use(middleware.AuthMiddleware(cfg.JWTSecret))
		staff.Use(middleware.RequireRole("admin", "cashier"))
		{
			staff.GET("/", productHandler.GetAllProducts)
			staff.GET("/:id", productHandler.GetProductByID)
		}

		adminProducts := v1.Group("/products")
		adminProducts.Use(middleware.AuthMiddleware(cfg.JWTSecret))
		adminProducts.Use(middleware.RequireRole("admin"))
		{
			adminProducts.POST("/", productHandler.CreateProduct)
			adminProducts.PUT("/:id", productHandler.UpdateProduct)
			adminProducts.DELETE("/:id", productHandler.DeleteProduct)
		}

		pos := v1.Group("/transactions")
		pos.Use(middleware.AuthMiddleware(cfg.JWTSecret))
		pos.Use(middleware.RequireRole("admin", "cashier"))
		{
			pos.POST("/inquiry", transHandler.Inquiry)
			pos.POST("/checkout", transHandler.Checkout)
		}

		reports := v1.Group("/reports")
		reports.Use(middleware.AuthMiddleware(cfg.JWTSecret))
		reports.Use(middleware.RequireRole("admin"))
		{
			reports.GET("/daily-summary", transHandler.GetTodayDashboardSummary)
		}
	}
}
