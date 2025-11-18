package utils

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type pgxpoolQuerier interface {
	Exec(ctx context.Context, sql string, arguments ...any) (pgconn.CommandTag, error)
	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
}

func RunExec(
	pool *pgxpool.Pool,
	ctx context.Context,
	sql string,
	args ...any,
) (pgconn.CommandTag, error) {
	var querier pgxpoolQuerier

	tx := GetTx(ctx)
	if tx != nil {
		querier = tx
	} else {
		conn, err := pool.Acquire(ctx)
		if err != nil {
			return pgconn.CommandTag{}, err
		}
		defer conn.Release()

		querier = conn
	}

	return querier.Exec(ctx, sql, args...)
}

func RunQueryRow(
	pool *pgxpool.Pool,
	scan func(row pgx.Row) error,
	ctx context.Context,
	sql string,
	args ...any,
) error {
	var querier pgxpoolQuerier

	tx := GetTx(ctx)
	if tx != nil {
		querier = tx
	} else {
		conn, err := pool.Acquire(ctx)
		if err != nil {
			return err
		}
		defer conn.Release()

		querier = conn
	}

	return scan(querier.QueryRow(ctx, sql, args...))
}

func RunQuery(
	pool *pgxpool.Pool,
	scan func(rows pgx.Rows) error,
	ctx context.Context,
	sql string,
	args ...any,
) error {
	var querier pgxpoolQuerier

	tx := GetTx(ctx)
	if tx != nil {
		querier = tx
	} else {
		conn, err := pool.Acquire(ctx)
		if err != nil {
			return err
		}
		defer conn.Release()

		querier = conn
	}

	rows, err := querier.Query(ctx, sql, args...)
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		if err := scan(rows); err != nil {
			return err
		}
	}

	return rows.Err()
}

func RunTx(pool *pgxpool.Pool, ctx context.Context, fn func(tx pgx.Tx) error) error {
	conn, err := pool.Acquire(ctx)
	if err != nil {
		return err
	}
	defer conn.Release()

	tx, err := conn.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	if err := fn(tx); err != nil {
		return err
	}

	return tx.Commit(ctx)
}
