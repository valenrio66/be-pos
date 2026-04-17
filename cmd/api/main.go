package main

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/valenrio66/be-pos/config"
	"github.com/valenrio66/be-pos/internal/delivery/http"
	"github.com/valenrio66/be-pos/internal/delivery/http/handler"
	"github.com/valenrio66/be-pos/internal/domain"
	"github.com/valenrio66/be-pos/internal/infrastructure/database"
	"github.com/valenrio66/be-pos/internal/infrastructure/logger"
	"github.com/valenrio66/be-pos/internal/repository/postgres"
	"github.com/valenrio66/be-pos/internal/usecase"
	"go.uber.org/fx"
)

func main() {
	fx.New(
		fx.Provide(
			config.LoadConfig,
			logger.NewLogger,
			database.NewPostgresConn,

			fx.Annotate(
				postgres.NewUserRepository,
				fx.As(new(domain.UserRepository)),
			),
			usecase.NewUserUsecase,
			handler.NewUserHandler,

			fx.Annotate(
				postgres.NewProductRepository,
				fx.As(new(domain.ProductRepository)),
			),
			usecase.NewProductUsecase,
			handler.NewProductHandler,

			fx.Annotate(
				postgres.NewTransactionRepository,
				fx.As(new(domain.TransactionRepository)),
			),
			usecase.NewTransactionUsecase,
			handler.NewTransactionHandler,

			http.NewGinServer,
		),
		fx.Invoke(func(r *gin.Engine) {
			r.Use(cors.New(cors.Config{
				AllowOrigins:     []string{"*"},
				AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
				AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
				ExposeHeaders:    []string{"Content-Length"},
				AllowCredentials: true,
			}))
		}),
	).Run()
}
