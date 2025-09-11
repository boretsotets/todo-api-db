// Package database реализует создание соединения с базой данных
package database

import (
	"github.com/jackc/pgx/v5/pgxpool"

	"context"
)

// InitBb устанавливает соединение с базой данных Postgres
// и возвращает пул соединений. Используется для работы
// всех репозиториев
func InitDb(ctx context.Context, dsn string) (*pgxpool.Pool, error) {
	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		return nil, err
	}
	return pool, nil
}
