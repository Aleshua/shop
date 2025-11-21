package requests

import (
	"context"

	shUtils "shared/utils"
)

type ResendCodeRequest struct {
	UserId int32 `json:"user_id" validation:"required,min=1"`
}

func (r ResendCodeRequest) Validate(ctx context.Context) error {
	return shUtils.ValidateStruct(ctx, r)
}
