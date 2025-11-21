package services

import (
	"context"
	"crypto/rand"
	"encoding/hex"

	d "auth/internal/domain"
)

type TokenGeneratorService struct{}

func NewTokenGeneratorService() *TokenGeneratorService {
	return &TokenGeneratorService{}
}

func (g TokenGeneratorService) GenerateToken(ctx context.Context, length int) (string, error) {
	if length <= 0 {
		return "", d.ErrTokenGeneratorInvalidLength
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
