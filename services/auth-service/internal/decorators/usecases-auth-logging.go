package decorators

import (
	"context"
	"shared/logger"
	"time"

	ctxHelper "shared/ctxhelper"

	reps "auth/internal/usecases"
)

type IAuthUseCasesWithLogging struct {
	base   reps.IAuthUseCases
	logger logger.ILogger
}

func NewIAuthUseCasesWithLogging(
	base reps.IAuthUseCases,
	logger logger.ILogger,
) *IAuthUseCasesWithLogging {
	return &IAuthUseCasesWithLogging{
		base:   base,
		logger: logger,
	}
}

func (d *IAuthUseCasesWithLogging) Register(ctx context.Context, password, email string) (_ int32, err error) {
	start := time.Now()

	defer func() {
		l := d.logger.With(
			logger.NewField("request_id", ctxHelper.GetRequestIdFromContext(ctx)),
			logger.NewField("method", "Register"),
			logger.NewField("duration", time.Since(start)),
		)

		if err != nil {
			l.Errorf("вызвал ошибку: %v", err)
			return
		}

		l.Info("вызов прошёл успешно")
	}()

	return d.base.Register(ctx, password, email)
}

func (d *IAuthUseCasesWithLogging) ConfirmEmail(ctx context.Context, userId int32, code string) (err error) {
	start := time.Now()

	defer func() {
		l := d.logger.With(
			logger.NewField("request_id", ctxHelper.GetRequestIdFromContext(ctx)),
			logger.NewField("method", "ConfirmEmail"),
			logger.NewField("duration", time.Since(start)),
		)

		if err != nil {
			l.Errorf("вызвал ошибку: %v", err)
			return
		}

		l.Info("вызов прошёл успешно")
	}()

	return d.base.ConfirmEmail(ctx, userId, code)
}

func (d *IAuthUseCasesWithLogging) ResendCode(ctx context.Context, userId int32) (err error) {
	start := time.Now()

	defer func() {
		l := d.logger.With(
			logger.NewField("request_id", ctxHelper.GetRequestIdFromContext(ctx)),
			logger.NewField("method", "ResendCode"),
			logger.NewField("duration", time.Since(start)),
		)

		if err != nil {
			l.Errorf("вызвал ошибку: %v", err)
			return
		}

		l.Info("вызов прошёл успешно")
	}()

	return d.base.ResendCode(ctx, userId)
}

func (d *IAuthUseCasesWithLogging) RefreshToken(ctx context.Context, refreshToken string) (_ string, err error) {
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

	return d.base.RefreshToken(ctx, refreshToken)
}

func (d *IAuthUseCasesWithLogging) Login(ctx context.Context, password, email string) (_, _ string, err error) {
	start := time.Now()

	defer func() {
		l := d.logger.With(
			logger.NewField("request_id", ctxHelper.GetRequestIdFromContext(ctx)),
			logger.NewField("method", "Login"),
			logger.NewField("duration", time.Since(start)),
		)

		if err != nil {
			l.Errorf("вызвал ошибку: %v", err)
			return
		}

		l.Info("вызов прошёл успешно")
	}()

	return d.base.Login(ctx, password, email)
}

func (d *IAuthUseCasesWithLogging) Logout(ctx context.Context, refreshToken string) (err error) {
	start := time.Now()

	defer func() {
		l := d.logger.With(
			logger.NewField("request_id", ctxHelper.GetRequestIdFromContext(ctx)),
			logger.NewField("method", "Logout"),
			logger.NewField("duration", time.Since(start)),
		)

		if err != nil {
			l.Errorf("вызвал ошибку: %v", err)
			return
		}

		l.Info("вызов прошёл успешно")
	}()

	return d.base.Logout(ctx, refreshToken)
}
