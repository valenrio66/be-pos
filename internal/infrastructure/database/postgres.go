package database

import (
	"database/sql"

	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
	"github.com/valenrio66/be-pos/config"
	"go.uber.org/zap"
)

func NewPostgresConn(cfg *config.Config, l *zap.Logger) *bun.DB {
	sqldb := sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(cfg.DatabaseURL)))
	db := bun.NewDB(sqldb, pgdialect.New())

	if err := db.Ping(); err != nil {
		l.Fatal("Failed to connect to the database", zap.Error(err))
	}

	l.Info("Database connection initialized successfully")
	return db
}
