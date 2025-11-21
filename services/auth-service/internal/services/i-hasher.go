package services

import "context"

type IHasherService interface {
	Hash(ctx context.Context, value string) (string, error)
	Compare(ctx context.Context, value, hash string) bool
}
