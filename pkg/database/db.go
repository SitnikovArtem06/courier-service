package database

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"os"
	"time"
)

const (
	TimeOut       = 5
	RetryAttempts = 5
	RetryDelay    = time.Second
)

func InitDb(ctx context.Context) (*pgxpool.Pool, error) {

	port := os.Getenv("POSTGRES_PORT")
	dbName := os.Getenv("POSTGRES_DB")
	dbUser := os.Getenv("POSTGRES_USER")
	dbPassword := os.Getenv("POSTGRES_PASSWORD")
	host := os.Getenv("POSTGRES_HOST")

	connStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s", dbUser, dbPassword, host, port, dbName)

	var lastErr error

	for i := 1; i <= RetryAttempts; i++ {

		tctx, cancel := context.WithTimeout(ctx, TimeOut*time.Second)

		dbpool, err := pgxpool.New(tctx, connStr)
		if err != nil {
			lastErr = fmt.Errorf("new pool fail: %w", err)

		} else {
			if err = dbpool.Ping(tctx); err != nil {
				dbpool.Close()
				lastErr = fmt.Errorf("db ping: %w", err)
			} else {
				cancel()
				return dbpool, nil
			}
		}

		cancel()

		if i < RetryAttempts {
			select {
			case <-ctx.Done():
				return nil, ctx.Err()
			case <-time.After(RetryDelay):

			}
		}

	}

	return nil, fmt.Errorf("db connect failed after %d attempts: %w", RetryAttempts, lastErr)
}
