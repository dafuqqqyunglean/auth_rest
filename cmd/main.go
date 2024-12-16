package main

import (
	"auth_rest/cmd/application"
	"auth_rest/internal/config"
	"context"
	"go.uber.org/zap"
	"os"
	"os/signal"
	"syscall"
)

// @title Auth API
// @version 1.0
// @description Authorization API Server

// @host localhost:8080
// @BasePath /

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	prvLogger, _ := zap.NewProduction()
	defer prvLogger.Sync()
	logger := prvLogger.Sugar()

	cfg := config.NewConfig()

	app := application.NewApp(ctx, logger, cfg)
	app.InitService()
	if err := app.Run(); err != nil {
		logger.Errorf(err.Error())
		return
	}

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGTERM, syscall.SIGINT)
	<-sigChan

	if err := app.Shutdown(); err != nil {
		logger.Errorf(err.Error())
		return
	}
}
