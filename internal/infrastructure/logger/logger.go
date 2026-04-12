package logger

import "go.uber.org/zap"

func NewLogger() *zap.Logger {
	config := zap.NewProductionConfig()
	l, _ := config.Build()
	return l
}
