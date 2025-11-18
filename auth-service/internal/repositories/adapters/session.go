package adapters

import (
	"context"

	d "auth/internal/domain"
	"auth/internal/utils"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type SessionRepository struct {
	pool *pgxpool.Pool
}

func NewSessionRepository(pool *pgxpool.Pool) *SessionRepository {
	return &SessionRepository{pool: pool}
}

func (r *SessionRepository) Find(ctx context.Context, refreshTokenHash string) (d.Session, error) {
	sql := `
		SELECT id, user_id, token_hash, expires_at, revoked, created_at
		FROM refresh_tokens
		WHERE token_hash = $1
	`

	var session d.Session
	err := utils.RunQueryRow(
		r.pool,
		func(row pgx.Row) error {
			return row.Scan(
				&session.Id,
				&session.UserId,
				&session.TokenHash,
				&session.ExpiresAt,
				&session.Revoked,
				&session.CreatedAt,
			)
		},
		ctx,
		sql,
		refreshTokenHash,
	)

	return session, err
}

func (r *SessionRepository) Update(ctx context.Context, session d.Session) error {
	sql := `
		UPDATE refresh_tokens
		SET user_id = $1, token_hash = $2, expires_at = $3, revoked = $4
		WHERE id = $5
	`

	_, err := utils.RunExec(
		r.pool,
		ctx,
		sql,
		session.UserId,
		session.TokenHash,
		session.ExpiresAt,
		session.Revoked,
		session.Id,
	)
	return err
}

func (r *SessionRepository) Save(ctx context.Context, session d.Session) error {
	sql := `
		INSERT INTO refresh_tokens (user_id, token_hash, expires_at, revoked, created_at)
		VALUES ($1, $2, $3, $4, $5)
	`

	_, err := utils.RunExec(
		r.pool,
		ctx,
		sql,
		session.UserId,
		session.TokenHash,
		session.ExpiresAt,
		session.Revoked,
		session.CreatedAt,
	)
	return err
}
