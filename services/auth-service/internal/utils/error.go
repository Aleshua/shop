package utils

import (
	"errors"
	"net/http"

	d "auth/internal/domain"
)

func TranslateErrorToHTTP(err error) (string, int) {
	switch {
	case errors.Is(err, d.ErrEntityNotFound):
		return "данные не найдены", http.StatusNotFound
	case errors.Is(err, d.ErrUserEmailConflict):
		return "аккаунт с таким email уже зарегистрирован", http.StatusConflict
	}

	switch {
	case errors.Is(err, d.ErrJWTTokenExpired):
		return "срок действия токена прошёл", http.StatusUnauthorized
	case errors.Is(err, d.ErrJWTInvalidToken):
		return "токен не валидный", http.StatusUnauthorized
	case errors.Is(err, d.ErrJWTInvalidMethodSignature):
		return "неожиданный метод подписи", http.StatusUnauthorized
	}

	switch {
	case errors.Is(err, d.ErrVerificationAttemptsExceeded):
		return "превышено количество попыток ввода кода подтверждения", http.StatusForbidden
	case errors.Is(err, d.ErrVerificationCodeExpired):
		return "срок действия кода подтверждения истёк", http.StatusGone
	case errors.Is(err, d.ErrVerificationInvalidCode):
		return "неверный код", http.StatusBadRequest
	case errors.Is(err, d.ErrRefreshTokenExpired):
		return "срок действия refresh токена истёк", http.StatusUnauthorized
	case errors.Is(err, d.ErrRefreshTokenRevoked):
		return "refresh токен добавлен в blacklist", http.StatusUnauthorized
	case errors.Is(err, d.ErrInvalidPassword):
		return "введён не правильный пароль", http.StatusUnauthorized
	case errors.Is(err, d.ErrEmailNotVerified):
		return "почта не подтверждена", http.StatusForbidden
	}

	return "внутренняя ошибка сервера", http.StatusInternalServerError
}
