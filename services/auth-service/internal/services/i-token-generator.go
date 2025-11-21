package services

import "context"

type ITokenGeneratorService interface {
	// length - количество симолов
	GenerateToken(ctx context.Context, length int) (string, error)
}
