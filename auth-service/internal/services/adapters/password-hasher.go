package adapters

import (
	"golang.org/x/crypto/bcrypt"
)

type PasswordHasherService struct{}

func NewPasswordHasherService() *PasswordHasherService {
	return &PasswordHasherService{}
}

func (s *PasswordHasherService) Hash(value string) (string, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(value), 10)
	return string(hashed), err
}

func (s *PasswordHasherService) Compare(value, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(value))
	return err == nil
}
