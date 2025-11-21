package services

import (
	"context"

	"golang.org/x/crypto/bcrypt"
)

type PasswordHasherService struct{}

func NewPasswordHasherService() *PasswordHasherService {
	return &PasswordHasherService{}
}

func (s *PasswordHasherService) Hash(ctx context.Context, value string) (string, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(value), 10)
	return string(hashed), err
}

func (s *PasswordHasherService) Compare(ctx context.Context, value, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(value))
	return err == nil
}
