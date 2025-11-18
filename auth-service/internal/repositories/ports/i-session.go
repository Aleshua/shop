package ports

import (
	"context"

	d "auth/internal/domain"
)

type ISessionRepository interface {
	// Находит сессию по хешу токена.
	// Если не находит выбрасывает ошибку
	Find(ctx context.Context, refreshTokenHash string) (d.Session, error)

	// По id обновляет все параметры
	Update(ctx context.Context, session d.Session) error

	Save(ctx context.Context, session d.Session) error
}
