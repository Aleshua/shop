package usecases

import "errors"

var (
	ErrVerificationAttemptsExceeded = errors.New("превышено количество попыток ввода кода подтверждения")
	ErrVerificationCodeExpired      = errors.New("срок действия кода подтверждения истёк")
	ErrVerificationInvalidCode      = errors.New("неверный код")

	ErrRefreshTokenExpired = errors.New("срок действия refresh токена истёк")
	ErrRefreshTokenRevoked = errors.New("refresh токен добавлен в blacklist")

	ErrInvalidPassword  = errors.New("введён не правильный пароль")
	ErrEmailNotVerified = errors.New("почта не подтверждена")
)
