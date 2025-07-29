package db

import (
	"context"
	"dozenChairs/pkg/config"
	"dozenChairs/pkg/logger"
	"go.uber.org/zap"
	"time"
)

import "github.com/jackc/pgx/v5/pgxpool"

func Connect(ctx context.Context, dsn string) (*pgxpool.Pool, error) {
	return pgxpool.New(ctx, dsn)
}

func MustConnectDB(cfg *config.Config, log logger.Logger) *pgxpool.Pool {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	conn, err := Connect(ctx, cfg.DatabaseDSN)
	if err != nil {
		log.Fatal("failed to connect to database", zap.Error(err))
	}
	return conn
}
