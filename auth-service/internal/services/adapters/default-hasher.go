package adapters

import (
	"crypto/sha256"
	"encoding/hex"
)

type DefaultHasherService struct{}

func NewDefaultHasherService() *DefaultHasherService {
	return &DefaultHasherService{}
}

func (s *DefaultHasherService) Hash(value string) (string, error) {
	sum := sha256.Sum256([]byte(value))
	return hex.EncodeToString(sum[:]), nil
}

func (s *DefaultHasherService) Compare(value, hash string) bool {
	valueHash, _ := s.Hash(value)
	return valueHash == hash
}
