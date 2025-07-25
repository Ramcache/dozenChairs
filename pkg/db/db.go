package db

import (
	"context"

	"github.com/jackc/pgx/v5"
)

func Connect(ctx context.Context, dsn string) (*pgx.Conn, error) {
	return pgx.Connect(ctx, dsn)
}
