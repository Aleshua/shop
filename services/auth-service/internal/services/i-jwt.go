package services

import "context"

type IJWTService interface {
	// Возвращает access token
	GenerateToken(ctx context.Context, userId int32) (string, error)

	// Возвращает userId если токен валидный
	ValidateToken(ctx context.Context, token string) (int32, error)

	// Возвращает userId без проверки
	ExtractClaims(ctx context.Context, token string) (int32, error)

	// Возвращает токен с обновлённым временем
	RefreshToken(ctx context.Context, oldToken string) (string, error)
}
