package services

import "errors"

var (
	ErrTokenGeneratorInvalidLength = errors.New("длина должна быть > 0")

	ErrJWTInvalidMethodSignature = errors.New("неожиданный способ подписи")
	ErrJWTInvalidInvalidToken    = errors.New("недействительный токен")
	ErrJWTTokenExpired           = errors.New("токен истёк")
)
