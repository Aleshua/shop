package decorators

import (
	"context"
	"shared/logger"
	"time"

	ctxHelper "shared/ctxhelper"

	reps "auth/internal/services"
)

type IJWTServiceWithLogging struct {
	base   reps.IJWTService
	logger logger.ILogger
}

func NewIJWTServiceWithLogging(
	base reps.IJWTService,
	logger logger.ILogger,
) *IJWTServiceWithLogging {
	return &IJWTServiceWithLogging{
		base:   base,
		logger: logger,
	}
}

func (d *IJWTServiceWithLogging) GenerateToken(ctx context.Context, userId int32) (_ string, err error) {
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

	return d.base.GenerateToken(ctx, userId)
}

func (d *IJWTServiceWithLogging) ValidateToken(ctx context.Context, token string) (_ int32, err error) {
	start := time.Now()

	defer func() {
		l := d.logger.With(
			logger.NewField("request_id", ctxHelper.GetRequestIdFromContext(ctx)),
			logger.NewField("method", "ValidateToken"),
			logger.NewField("duration", time.Since(start)),
		)

		if err != nil {
			l.Errorf("вызвал ошибку: %v", err)
			return
		}

		l.Info("вызов прошёл успешно")
	}()

	return d.base.ValidateToken(ctx, token)
}

func (d *IJWTServiceWithLogging) ExtractClaims(ctx context.Context, token string) (_ int32, err error) {
	start := time.Now()

	defer func() {
		l := d.logger.With(
			logger.NewField("request_id", ctxHelper.GetRequestIdFromContext(ctx)),
			logger.NewField("method", "ExtractClaims"),
			logger.NewField("duration", time.Since(start)),
		)

		if err != nil {
			l.Errorf("вызвал ошибку: %v", err)
			return
		}

		l.Info("вызов прошёл успешно")
	}()

	return d.base.ExtractClaims(ctx, token)
}

func (d *IJWTServiceWithLogging) RefreshToken(ctx context.Context, oldToken string) (_ string, err error) {
	start := time.Now()

	defer func() {
		l := d.logger.With(
			logger.NewField("request_id", ctxHelper.GetRequestIdFromContext(ctx)),
			logger.NewField("method", "RefreshToken"),
			logger.NewField("duration", time.Since(start)),
		)

		if err != nil {
			l.Errorf("вызвал ошибку: %v", err)
			return
		}

		l.Info("вызов прошёл успешно")
	}()

	return d.base.RefreshToken(ctx, oldToken)
}
