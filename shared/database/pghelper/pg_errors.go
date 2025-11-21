package pghelper

import (
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

// IsUniqueConstraint проверяет, является ли ошибка нарушением уникальности
// для конкретного ограничения (constraintName).
func IsUniqueConstraint(err error, constraintName string) bool {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		// 23505 = unique_violation
		if pgErr.Code == "23505" {
			return pgErr.ConstraintName == constraintName
		}
	}
	return false
}

// IsErrNoRows проверяет, является ли ошибка pgx.ErrNoRows.
func IsErrNoRows(err error) bool {
	if errors.Is(err, pgx.ErrNoRows) {
		return true
	}
	return false
}

// CheckRowsAffected проверяет, были ли затронуты строки при UPDATE/DELETE.
// Если ошибок нет, но затронуто 0 строк -> возвращает true.
func CheckRowsAffected(tag pgconn.CommandTag, err error) bool {
	if err != nil {
		return false
	}
	if tag.RowsAffected() == 0 {
		return true
	}
	return false
}
