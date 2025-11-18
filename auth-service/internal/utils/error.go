package utils

import (
	"errors"
	"net/http"

	"auth/internal/repositories"
	"auth/internal/usecases"

	"github.com/golang-jwt/jwt/v5"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

func TranslateErrorToHTTP(err error) (string, int) {
	// JWT ошибки
	switch {
	case errors.Is(err, jwt.ErrTokenExpired):
		return "токен истёк", http.StatusUnauthorized
	case errors.Is(err, jwt.ErrTokenNotValidYet):
		return "токен ещё не действителен", http.StatusUnauthorized
	case errors.Is(err, jwt.ErrTokenMalformed):
		return "некорректный токен", http.StatusBadRequest
	case errors.Is(err, jwt.ErrTokenSignatureInvalid):
		return "неверная подпись токена", http.StatusUnauthorized
	}

	// Ошибки pgxpool / pgx
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		switch pgErr.Code {
		case "23505":
			return "запись уже существует", http.StatusConflict
		default:
			return "внутренняя ошибка сервера", http.StatusInternalServerError
		}
	}
	if errors.Is(err, pgx.ErrNoRows) {
		return "данные не найдены", http.StatusNotFound
	}

	// Мои ошибки
	switch {
	case errors.Is(err, usecases.ErrVerificationAttemptsExceeded):
		return "превышено количество попыток ввода кода подтверждения", http.StatusForbidden
	case errors.Is(err, usecases.ErrVerificationCodeExpired):
		return "срок действия кода подтверждения истёк", http.StatusGone
	case errors.Is(err, usecases.ErrVerificationInvalidCode):
		return "неверный код", http.StatusBadRequest
	case errors.Is(err, usecases.ErrRefreshTokenExpired):
		return "срок действия refresh токена истёк", http.StatusUnauthorized
	case errors.Is(err, usecases.ErrRefreshTokenRevoked):
		return "refresh токен добавлен в blacklist", http.StatusUnauthorized
	case errors.Is(err, usecases.ErrInvalidPassword):
		return "введён не правильный пароль", http.StatusUnauthorized
	case errors.Is(err, repositories.ErrEmailConflict):
		return "аккаунт с таким email уже зарегистрирован", http.StatusConflict
	}

	return "внутренняя ошибка сервера", http.StatusInternalServerError
}
