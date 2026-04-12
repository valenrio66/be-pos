package http

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/valenrio66/be-pos/config"
	"github.com/valenrio66/be-pos/internal/delivery/http/handler"
	"github.com/valenrio66/be-pos/internal/delivery/http/middleware"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

func NewGinServer(lc fx.Lifecycle, cfg *config.Config, log *zap.Logger, userHandler *handler.UserHandler) *gin.Engine {
	r := gin.Default()

	r.Use(middleware.ZapLogger(log))
	SetupRoutes(r, userHandler)

	addr := fmt.Sprintf(":%s", cfg.AppPort)

	srv := &http.Server{
		Addr:    addr,
		Handler: r,
	}

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			log.Info("Starting HTTP server", zap.String("port", cfg.AppPort))

			go func() {
				if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
					log.Fatal("HTTP server gagal berjalan", zap.Error(err))
				}
			}()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			log.Info("Menerima sinyal mati. Menghentikan HTTP server secara Graceful...")
			return srv.Shutdown(ctx)
		},
	})

	return r
}
