// Package database реализует создание соединения с базой данных
package database

import (
	"github.com/jackc/pgx/v5/pgxpool"

	"context"
)

// Функция InitDb реализует создание соединения с базой данных 
func InitDb(ctx context.Context, dsn string) (*pgxpool.Pool, error) {
	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		return nil, err
	}
	return pool, nil
}
