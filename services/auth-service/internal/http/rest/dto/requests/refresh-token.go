package requests

import (
	"context"

	shUtils "shared/utils"
)

type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
}

func (r RefreshTokenRequest) Validate(ctx context.Context) error {
	return shUtils.ValidateStruct(ctx, r)
}
