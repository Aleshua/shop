package decorators

import (
	"context"
	ctxHelper "shared/ctxhelper"
	"shared/logger"
	"time"

	reps "auth/internal/services"
)

type IEmailSenderServiceWithLogging struct {
	base   reps.IEmailSenderService
	logger logger.ILogger
}

func NewIEmailSenderServiceWithLogging(
	base reps.IEmailSenderService,
	logger logger.ILogger,
) *IEmailSenderServiceWithLogging {
	return &IEmailSenderServiceWithLogging{
		base:   base,
		logger: logger,
	}
}

func (d *IEmailSenderServiceWithLogging) Send(ctx context.Context, to string, body string) (err error) {
	start := time.Now()

	defer func() {
		l := d.logger.With(
			logger.NewField("request_id", ctxHelper.GetRequestIdFromContext(ctx)),
			logger.NewField("method", "Send"),
			logger.NewField("duration", time.Since(start)),
		)

		if err != nil {
			l.Errorf("вызвал ошибку: %v", err)
			return
		}

		l.Info("вызов прошёл успешно")
	}()

	return d.base.Send(ctx, to, body)
}
