package utils

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Transaction interface {
	WithinReadCommittedTx(ctx context.Context, fn func(ctx context.Context) error) error
}

type pgTxManager struct {
	db *pgxpool.Pool
}

func NewTxManager(db *pgxpool.Pool) Transaction {
	return &pgTxManager{db: db}
}

func (m *pgTxManager) WithinReadCommittedTx(ctx context.Context, fn func(ctx context.Context) error) error {
	tx, err := m.db.BeginTx(ctx, pgx.TxOptions{
		IsoLevel:   pgx.ReadCommitted,
		AccessMode: pgx.ReadWrite,
	})
	if err != nil {
		return err
	}

	ctx = context.WithValue(ctx, txKey{}, tx)

	if err := fn(ctx); err != nil {
		_ = tx.Rollback(ctx)
		return err
	}

	return tx.Commit(ctx)
}

type txKey struct{}

func GetTx(ctx context.Context) pgx.Tx {
	val := ctx.Value(txKey{})
	if val == nil {
		return nil
	}
	return val.(pgx.Tx)
}
