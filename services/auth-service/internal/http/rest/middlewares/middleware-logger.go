package middlewares

import (
	ctxHelper "shared/ctxhelper"
	l "shared/logger"
	"time"

	"github.com/gin-contrib/requestid"
	"github.com/gin-gonic/gin"
)

func LoggerMiddleware(logger l.ILogger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		reqId := requestid.Get(c)
		ctxHelper.AddRequestIdToGinContext(c, reqId)

		c.Next()

		status := c.Writer.Status()

		fields := map[string]interface{}{
			"method":     c.Request.Method,
			"path":       c.Request.URL.Path,
			"status":     status,
			"duration":   time.Since(start),
			"request_id": reqId,
			"ip":         c.ClientIP(),
		}

		lw := logger.With(l.NewField("fields", fields))

		if status >= 500 {
			lw.Error("http запрос завершился с ошибкой")
		} else {
			lw.Info("http запрос")
		}
	}
}
