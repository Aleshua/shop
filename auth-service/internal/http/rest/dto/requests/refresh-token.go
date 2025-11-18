package requests

import (
	"context"

	"auth/internal/utils"
)

type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
}

func (r RefreshTokenRequest) Validate(ctx context.Context) error {
	return utils.ValidateStruct(ctx, r)
}
