package database

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Options struct {
	Path              string
	MaxConns          int
	MinConns          int
	MaxConnLifetime   time.Duration
	MaxConnIdleTime   time.Duration
	HealthCheckPeriod time.Duration
	ConnectTimeout    time.Duration
}

type DB struct {
	Pool *pgxpool.Pool
}

func New(options Options) (*DB, error) {
	dbConfig, err := pgxpool.ParseConfig(options.Path)
	if err != nil {
		return nil, fmt.Errorf("ошибка создания конфига pgxpool: %s", err.Error())
	}

	dbConfig.MaxConns = int32(options.MaxConns)
	dbConfig.MinConns = int32(options.MinConns)
	dbConfig.MaxConnLifetime = options.MaxConnLifetime
	dbConfig.MaxConnIdleTime = options.MaxConnIdleTime
	dbConfig.HealthCheckPeriod = options.HealthCheckPeriod
	dbConfig.ConnConfig.ConnectTimeout = options.ConnectTimeout

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
