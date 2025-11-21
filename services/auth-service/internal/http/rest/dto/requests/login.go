package requests

import (
	"context"

	shUtils "shared/utils"
)

type LoginRequest struct {
	Password string `json:"password" validate:"required"`
	Email    string `json:"email" validate:"required,email"`
}

func (r LoginRequest) Validate(ctx context.Context) error {
	return shUtils.ValidateStruct(ctx, r)
}
