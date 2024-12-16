package storage

import (
	"auth_rest/internal/config"
	"auth_rest/internal/utils"
	"fmt"
	"github.com/jackc/pgx/v5"
	"os"
)

func NewPostgresDB(cfg config.PostgresConfig, ctx *utils.AppContext) *pgx.Conn {
	cfgStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s", cfg.Username, cfg.Password, cfg.Host, cfg.Port, cfg.DBName, cfg.SSLMode)
	conn, err := pgx.Connect(ctx.Ctx, cfgStr)
	if err != nil {
		ctx.Logger.Error(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}

	return conn
}
