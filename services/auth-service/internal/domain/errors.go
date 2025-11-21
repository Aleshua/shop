package domain

import "errors"

var (
	// Общие ошибки для репозиториев
	ErrEntityNotFound = errors.New("сущность не найдена")

	// Ошибки из user репозитория
	ErrUserEmailConflict = errors.New("email занят другим пользователем")

	// Ошибки для token generator сервиса
	ErrTokenGeneratorInvalidLength = errors.New("минимальная длина должна быть 1")

	// Ошибки для jwt сервиса
	ErrJWTTokenExpired           = errors.New("срок действия токена прошёл")
	ErrJWTInvalidToken           = errors.New("токен не валидный")
	ErrJWTInvalidMethodSignature = errors.New("неожиданный метод подписи")

	// Ошибки для auth usecases
	ErrVerificationAttemptsExceeded = errors.New("закончились попытки для подтверждения почты")
	ErrVerificationInvalidCode      = errors.New("неверный код")
	ErrVerificationCodeExpired      = errors.New("время жизни токена прошло")
	ErrRefreshTokenRevoked          = errors.New("refresh токен добавлен в чёрный список")
	ErrRefreshTokenExpired          = errors.New("время жизни refresh токена прошло")
	ErrInvalidPassword              = errors.New("неверный пароль")
	ErrEmailNotVerified             = errors.New("почта не подтверждена")
)
