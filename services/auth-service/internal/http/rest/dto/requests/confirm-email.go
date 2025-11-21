package requests

import (
	"context"
	shUtils "shared/utils"
)

type ConfirmEmailRequest struct {
	Code   string `json:"code" validate:"required"`
	UserId int32  `json:"user_id" validate:"required,min=1"`
}

func (r ConfirmEmailRequest) Validate(ctx context.Context) error {
	return shUtils.ValidateStruct(ctx, r)
}
