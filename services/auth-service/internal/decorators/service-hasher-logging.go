package decorators

import (
	"context"
	"shared/logger"
	"time"

	ctxHelper "shared/ctxhelper"

	reps "auth/internal/services"
)

type IHasherServiceWithLogging struct {
	base   reps.IHasherService
	logger logger.ILogger
}

func NewIHasherServiceWithLogging(
	base reps.IHasherService,
	logger logger.ILogger,
) *IHasherServiceWithLogging {
	return &IHasherServiceWithLogging{
		base:   base,
		logger: logger,
	}
}

func (d *IHasherServiceWithLogging) Hash(ctx context.Context, value string) (_ string, err error) {
	start := time.Now()

	defer func() {
		l := d.logger.With(
			logger.NewField("request_id", ctxHelper.GetRequestIdFromContext(ctx)),
			logger.NewField("method", "Hash"),
			logger.NewField("duration", time.Since(start)),
		)

		if err != nil {
			l.Errorf("вызвал ошибку: %v", err)
			return
		}

		l.Info("вызов прошёл успешно")
	}()

	return d.base.Hash(ctx, value)
}

func (d *IHasherServiceWithLogging) Compare(ctx context.Context, value, hash string) (_ bool) {
	start := time.Now()

	defer func() {
		l := d.logger.With(
			logger.NewField("request_id", ctxHelper.GetRequestIdFromContext(ctx)),
			logger.NewField("method", "Compare"),
			logger.NewField("duration", time.Since(start)),
		)

		l.Info("вызов прошёл успешно")
	}()

	return d.base.Compare(ctx, value, hash)
}
