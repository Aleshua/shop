package adapters

import (
	"crypto/rand"
	"encoding/hex"

	"auth/internal/services"
)

type TokenGeneratorService struct{}

func NewTokenGeneratorService() *TokenGeneratorService {
	return &TokenGeneratorService{}
}

func (g TokenGeneratorService) GenerateToken(length int) (string, error) {
	if length <= 0 {
		return "", services.ErrTokenGeneratorInvalidLength
	}

	buf := make([]byte, (length+1)/2)
	_, err := rand.Read(buf)
	if err != nil {
		return "", err
	}

	token := hex.EncodeToString(buf)
	token = token[:length]
	return token, nil
}
