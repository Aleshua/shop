package services

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
)

type DefaultHasherService struct{}

func NewDefaultHasherService() *DefaultHasherService {
	return &DefaultHasherService{}
}

func (s *DefaultHasherService) Hash(ctx context.Context, value string) (string, error) {
	sum := sha256.Sum256([]byte(value))
	return hex.EncodeToString(sum[:]), nil
}

func (s *DefaultHasherService) Compare(ctx context.Context, value, hash string) bool {
	valueHash, _ := s.Hash(ctx, value)
	return valueHash == hash
}
