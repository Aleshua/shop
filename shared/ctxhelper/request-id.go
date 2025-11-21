package ctxhelper

import (
	"context"

	"github.com/gin-gonic/gin"
)

type ctxKey string

const RequestIDKey ctxKey = "request_id"

func AddRequestIdToGinContext(c *gin.Context, reqID string) {
	ctx := context.WithValue(c.Request.Context(), RequestIDKey, reqID)
	c.Request = c.Request.WithContext(ctx)
}

func GetRequestIdFromContext(ctx context.Context) string {
	if reqID, ok := ctx.Value(RequestIDKey).(string); ok {
		return reqID
	}
	return ""
}
