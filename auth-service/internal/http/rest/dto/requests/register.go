package requests

import (
	"context"

	"auth/internal/utils"
)

type RegisterRequest struct {
	Password string `json:"password" validate:"required,min=8,max=32"`
	Email    string `json:"email" validate:"required,email,max=255"`
}

func (r RegisterRequest) Validate(ctx context.Context) error {
	return utils.ValidateStruct(ctx, r)
}
