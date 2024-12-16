package application

import (
	"auth_rest/internal/api"
	"auth_rest/internal/config"
	"auth_rest/internal/lib/tokens"
	"auth_rest/internal/services"
	"auth_rest/internal/storage"
	"auth_rest/internal/storage/sql"
	"auth_rest/internal/utils"
	"context"
	"go.uber.org/zap"
)

type App struct {
	ctx          *utils.AppContext
	srv          *api.Server
	cfg          config.Config
	postgres     *sql.AuthPostgres
	tokenManager *tokens.Manager
}

func NewApp(context context.Context, logger *zap.SugaredLogger, cfg config.Config) *App {
	ctx := utils.NewAppContext(context, logger)

	postgres := sql.NewAuthPostgres(storage.NewPostgresDB(cfg.Postgres, ctx), *ctx)

	tokenManager, err := tokens.NewManager(cfg.SigningKey)
	if err != nil {
		ctx.Logger.Error(err)
	}

	return &App{
		ctx:          ctx,
		cfg:          cfg,
		postgres:     postgres,
		tokenManager: tokenManager,
	}
}

func (a *App) InitService() {
	newService := services.NewAuthService(a.postgres, a.tokenManager, a.ctx, a.cfg.AccessTTL, a.cfg.RefreshTTL)
	a.srv = api.NewServer()
	a.srv.HandleAuth(newService, a.ctx)
}

func (a *App) Run() error {
	go func() {
		err := a.srv.Run()
		if err != nil {
			a.ctx.Logger.Fatalf("error running http server: %s", err.Error())
		}
	}()

	a.ctx.Logger.Info("server running")
	return nil
}

func (a *App) Shutdown() error {
	err := a.srv.Shutdown(a.ctx.Ctx)
	if err != nil {
		a.ctx.Logger.Errorf("failed to disconnect from server %v", err)
		return err
	}

	err = a.postgres.Close()
	if err != nil {
		a.ctx.Logger.Errorf("failed to disconnect from bd %v", err)
	}

	a.ctx.Logger.Info("server shutdown successfully")
	return nil
}
