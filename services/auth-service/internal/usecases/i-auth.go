package usecases

import (
	"context"
)

type IAuthUseCases interface {
	// Возвращает id созданного пользователя
	Register(ctx context.Context, password, email string) (int32, error)

	ConfirmEmail(ctx context.Context, userId int32, code string) error

	ResendCode(ctx context.Context, userId int32) error

	// Возвращает AccessToken
	RefreshToken(ctx context.Context, refreshToken string) (string, error)

	// Возвращает AccessToken и RefreshToken
	Login(ctx context.Context, password, email string) (string, string, error)

	// Добавляет токен в blacklist
	Logout(ctx context.Context, refreshToken string) error
}
