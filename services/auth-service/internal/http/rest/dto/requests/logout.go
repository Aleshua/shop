package requests

import (
	"context"

	shUtils "shared/utils"
)

type LogoutRequest struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
}

func (r LogoutRequest) Validate(ctx context.Context) error {
	return shUtils.ValidateStruct(ctx, r)
}
