package database

import (
	"context"
	"fmt"

	c "auth/internal/config"

	"github.com/jackc/pgx/v5/pgxpool"
)

type DB struct {
	Pool *pgxpool.Pool
}

func New(params c.Database) (*DB, error) {
	dbConfig, err := pgxpool.ParseConfig(params.Path)
	if err != nil {
		return nil, fmt.Errorf("ошибка создания конфига pgxpool: %s", err.Error())
	}

	dbConfig.MaxConns = params.MaxConns
	dbConfig.MinConns = params.MinConns
	dbConfig.MaxConnLifetime = params.MaxConnLifetime
	dbConfig.MaxConnIdleTime = params.MaxConnIdleTime
	dbConfig.HealthCheckPeriod = params.HealthCheckPeriod
	dbConfig.ConnConfig.ConnectTimeout = params.ConnectTimeout

	dbpool, err := pgxpool.NewWithConfig(context.Background(), dbConfig)
	if err != nil {
		return nil, fmt.Errorf("не удается создать пул подключений: %s", err.Error())
	}

	err = dbpool.Ping(context.Background())
	if err != nil {
		return nil, fmt.Errorf("не удается проверить пул подключений: %s", err.Error())
	}

	return &DB{Pool: dbpool}, nil
}
