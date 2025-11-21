package repositories

import (
	"context"
	"time"

	d "auth/internal/domain"
	pgh "shared/database/pghelper"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepository struct {
	pool *pgxpool.Pool
}

func NewUserRepository(pool *pgxpool.Pool) *UserRepository {
	return &UserRepository{pool: pool}
}

func (r *UserRepository) Find(ctx context.Context, userId int32) (d.User, error) {
	sql := `
		SELECT id, email, password_hash, email_verified_at, created_at, updated_at
		FROM users WHERE id = $1
	`
	var user d.User

	err := pgh.RunQueryRow(
		r.pool,
		func(row pgx.Row) error {
			return row.Scan(
				&user.Id,
				&user.Email,
				&user.PasswordHash,
				&user.EmailVerifiedAt,
				&user.CreatedAt,
				&user.UpdatedAt,
			)
		},
		ctx,
		sql,
		userId,
	)

	if pgh.IsErrNoRows(err) {
		return user, d.ErrEntityNotFound
	}

	return user, err
}

func (r *UserRepository) FindByEmail(ctx context.Context, email string) (d.User, error) {
	sql := `
		SELECT id, email, password_hash, email_verified_at, created_at, updated_at
		FROM users WHERE email = $1
	`
	var user d.User

	err := pgh.RunQueryRow(
		r.pool,
		func(row pgx.Row) error {
			return row.Scan(
				&user.Id,
				&user.Email,
				&user.PasswordHash,
				&user.EmailVerifiedAt,
				&user.CreatedAt,
				&user.UpdatedAt,
			)
		},
		ctx,
		sql,
		email,
	)

	if pgh.IsErrNoRows(err) {
		return user, d.ErrEntityNotFound
	}

	return user, err
}

func (r *UserRepository) Save(ctx context.Context, user d.User) (int32, error) {
	sql := `
		INSERT INTO users (email, password_hash, email_verified_at, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id
	`
	var id int32
	err := pgh.RunQueryRow(
		r.pool,
		func(row pgx.Row) error {
			return row.Scan(
				&id,
			)
		},
		ctx,
		sql,
		user.Email,
		user.PasswordHash,
		user.EmailVerifiedAt,
		user.CreatedAt,
		user.UpdatedAt,
	)

	if pgh.IsUniqueConstraint(err, "name") {
		return id, d.ErrUserEmailConflict
	}

	return id, err
}

func (r *UserRepository) SaveOrUpdateUnverified(ctx context.Context, user d.User) (int32, error) {
	sql := `
		INSERT INTO users (email, password_hash, email_verified_at, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5)
		ON CONFLICT (email)
		DO UPDATE
		SET password_hash = EXCLUDED.password_hash,
			updated_at = EXCLUDED.updated_at
		WHERE users.email_verified_at IS NULL
		RETURNING id;
	`

	var id int32
	err := pgh.RunQueryRow(
		r.pool,
		func(row pgx.Row) error {
			return row.Scan(
				&id,
			)
		},
		ctx,
		sql,
		user.Email,
		user.PasswordHash,
		user.EmailVerifiedAt,
		user.CreatedAt,
		user.UpdatedAt,
	)

	if pgh.IsErrNoRows(err) {
		return 0, d.ErrUserEmailConflict
	}

	return id, err
}

func (r *UserRepository) MarkEmailVerified(ctx context.Context, userId int32) error {
	sql := `
		UPDATE users
		SET email_verified_at = $1, updated_at = $1
		WHERE id = $2
	`

	tag, err := pgh.RunExec(r.pool, ctx, sql, time.Now(), userId)

	if pgh.CheckRowsAffected(tag, err) {
		return d.ErrEntityNotFound
	}

	return err
}

func (r *UserRepository) SaveOrUpdateConfirmCode(ctx context.Context, code d.ConfirmCode) error {
	sql := `
		INSERT INTO confirm_cods (user_id, code, attempts, created_at)
		VALUES ($1, $2, $3, $4)
		ON CONFLICT (user_id)
		DO UPDATE
		SET code = EXCLUDED.code,
			attempts = EXCLUDED.attempts,
			created_at = EXCLUDED.created_at;
	`
	_, err := pgh.RunExec(
		r.pool,
		ctx,
		sql,
		code.UserId,
		code.Code,
		code.Attempts,
		code.CreatedAt,
	)
	return err
}

func (r *UserRepository) FindConfirmCode(ctx context.Context, userId int32) (d.ConfirmCode, error) {
	sql := `
		SELECT id, user_id, code, attempts, created_at
		FROM confirm_cods
		WHERE user_id = $1
	`
	var code d.ConfirmCode
	err := pgh.RunQueryRow(
		r.pool,
		func(row pgx.Row) error {
			return row.Scan(
				&code.Id,
				&code.UserId,
				&code.Code,
				&code.Attempts,
				&code.CreatedAt,
			)
		},
		ctx,
		sql,
		userId,
	)

	if pgh.IsErrNoRows(err) {
		return code, d.ErrEntityNotFound
	}

	return code, err
}

func (r *UserRepository) UpdateConfirmCode(ctx context.Context, code d.ConfirmCode) error {
	sql := `
		UPDATE confirm_cods
		SET code = $1, attempts = $2, created_at = $3
		WHERE user_id = $4
	`
	tag, err := pgh.RunExec(
		r.pool,
		ctx,
		sql,
		code.Code,
		code.Attempts,
		code.CreatedAt,
		code.UserId,
	)

	if pgh.CheckRowsAffected(tag, err) {
		return d.ErrEntityNotFound
	}

	return err
}

func (r *UserRepository) DeleteConfirmCode(ctx context.Context, userId int32) error {
	sql := `DELETE FROM confirm_cods WHERE user_id = $1`

	tag, err := pgh.RunExec(
		r.pool,
		ctx,
		sql,
		userId,
	)

	if pgh.CheckRowsAffected(tag, err) {
		return d.ErrEntityNotFound
	}

	return err
}
