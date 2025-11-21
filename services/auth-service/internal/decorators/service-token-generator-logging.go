package decorators

import (
	"context"
	"shared/logger"
	"time"

	ctxHelper "shared/ctxhelper"

	reps "auth/internal/services"
)

type ITokenGeneratorServiceWithLogging struct {
	base   reps.ITokenGeneratorService
	logger logger.ILogger
}

func NewITokenGeneratorServiceWithLogging(
	base reps.ITokenGeneratorService,
	logger logger.ILogger,
) *ITokenGeneratorServiceWithLogging {
	return &ITokenGeneratorServiceWithLogging{
		base:   base,
		logger: logger,
	}
}

func (d *ITokenGeneratorServiceWithLogging) GenerateToken(ctx context.Context, length int) (_ string, err error) {
	start := time.Now()

	defer func() {
		l := d.logger.With(
			logger.NewField("request_id", ctxHelper.GetRequestIdFromContext(ctx)),
			logger.NewField("method", "GenerateToken"),
			logger.NewField("duration", time.Since(start)),
		)

		if err != nil {
			l.Errorf("вызвал ошибку: %v", err)
			return
		}

		l.Info("вызов прошёл успешно")
	}()

	return d.base.GenerateToken(ctx, length)
}
