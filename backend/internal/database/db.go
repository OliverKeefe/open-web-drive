package database

import (
	"context"
	"fmt"
	"log/slog"
	"os"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Pool interface {
	Begin(context.Context) (pgx.Tx, error)
	Exec(context.Context, string, ...any) (pgconn.CommandTag, error)
	Query(context.Context, string, ...any) (pgx.Rows, error)
	//QueryRow(context.Context, string, ...any) (pgx.Rows, error)
	QueryRow(context.Context, string, ...any) pgx.Row
}

type MetadataDatabase struct {
	Pool Pool
}

func New(ctx context.Context, envVar string) (*MetadataDatabase, error) {
	slog.Info("starting database.")

	dsn := os.Getenv(envVar)
	if dsn == "" {
		return nil, fmt.Errorf("environment variable %s is not set", envVar)
	}

	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		return nil, err
	}

	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, err
	}

	return &MetadataDatabase{
		Pool: pool,
	}, nil
}
