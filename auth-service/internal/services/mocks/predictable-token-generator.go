package mocks

import (
	"encoding/hex"

	"auth/internal/services"
)

type PredictableTokenGeneratorService struct{}

func NewPredictableTokenGeneratorService() *PredictableTokenGeneratorService {
	return &PredictableTokenGeneratorService{}
}

func (g PredictableTokenGeneratorService) GenerateToken(length int) (string, error) {
	if length <= 0 {
		return "", services.ErrTokenGeneratorInvalidLength
	}

	buf := make([]byte, (length+1)/2)
	token := hex.EncodeToString(buf)
	token = token[:length]
	return token, nil
}
