package decorators

import (
	"context"
	"shared/logger"
	"time"

	d "auth/internal/domain"
	reps "auth/internal/repositories"

	ctxHelper "shared/ctxhelper"
)

type IUserRepositoryWithLogging struct {
	base   reps.IUserRepository
	logger logger.ILogger
}

func NewIUserRepositoryWithLogging(
	base reps.IUserRepository,
	logger logger.ILogger,
) *IUserRepositoryWithLogging {
	return &IUserRepositoryWithLogging{
		base:   base,
		logger: logger,
	}
}

func (d *IUserRepositoryWithLogging) Find(ctx context.Context, userId int32) (_ d.User, err error) {
	start := time.Now()

	defer func() {
		l := d.logger.With(
			logger.NewField("request_id", ctxHelper.GetRequestIdFromContext(ctx)),
			logger.NewField("method", "Find"),
			logger.NewField("duration", time.Since(start)),
		)

		if err != nil {
			l.Errorf("вызвал ошибку: %v", err)
			return
		}

		l.Info("вызов прошёл успешно")
	}()

	return d.base.Find(ctx, userId)
}

func (d *IUserRepositoryWithLogging) FindByEmail(ctx context.Context, email string) (_ d.User, err error) {
	start := time.Now()

	defer func() {
		l := d.logger.With(
			logger.NewField("request_id", ctxHelper.GetRequestIdFromContext(ctx)),
			logger.NewField("method", "FindByEmail"),
			logger.NewField("duration", time.Since(start)),
		)

		if err != nil {
			l.Errorf("вызвал ошибку: %v", err)
			return
		}

		l.Info("вызов прошёл успешно")
	}()

	return d.base.FindByEmail(ctx, email)
}

func (d *IUserRepositoryWithLogging) Save(ctx context.Context, user d.User) (_ int32, err error) {
	start := time.Now()

	defer func() {
		l := d.logger.With(
			logger.NewField("request_id", ctxHelper.GetRequestIdFromContext(ctx)),
			logger.NewField("method", "Save"),
			logger.NewField("duration", time.Since(start)),
		)

		if err != nil {
			l.Errorf("вызвал ошибку: %v", err)
			return
		}

		l.Info("вызов прошёл успешно")
	}()

	return d.base.Save(ctx, user)
}

func (d *IUserRepositoryWithLogging) SaveOrUpdateUnverified(ctx context.Context, user d.User) (_ int32, err error) {
	start := time.Now()

	defer func() {
		l := d.logger.With(
			logger.NewField("request_id", ctxHelper.GetRequestIdFromContext(ctx)),
			logger.NewField("method", "SaveOrUpdateUnverified"),
			logger.NewField("duration", time.Since(start)),
		)

		if err != nil {
			l.Errorf("вызвал ошибку: %v", err)
			return
		}

		l.Info("вызов прошёл успешно")
	}()

	return d.base.SaveOrUpdateUnverified(ctx, user)
}

func (d *IUserRepositoryWithLogging) MarkEmailVerified(ctx context.Context, userId int32) (err error) {
	start := time.Now()

	defer func() {
		l := d.logger.With(
			logger.NewField("request_id", ctxHelper.GetRequestIdFromContext(ctx)),
			logger.NewField("method", "MarkEmailVerified"),
			logger.NewField("duration", time.Since(start)),
		)

		if err != nil {
			l.Errorf("вызвал ошибку: %v", err)
			return
		}

		l.Info("вызов прошёл успешно")
	}()

	return d.base.MarkEmailVerified(ctx, userId)
}

func (d *IUserRepositoryWithLogging) SaveOrUpdateConfirmCode(ctx context.Context, code d.ConfirmCode) (err error) {
	start := time.Now()

	defer func() {
		l := d.logger.With(
			logger.NewField("request_id", ctxHelper.GetRequestIdFromContext(ctx)),
			logger.NewField("method", "SaveOrUpdateConfirmCode"),
			logger.NewField("duration", time.Since(start)),
		)

		if err != nil {
			l.Errorf("вызвал ошибку: %v", err)
			return
		}

		l.Info("вызов прошёл успешно")
	}()

	return d.base.SaveOrUpdateConfirmCode(ctx, code)
}

func (d *IUserRepositoryWithLogging) FindConfirmCode(ctx context.Context, userId int32) (_ d.ConfirmCode, err error) {
	start := time.Now()

	defer func() {
		l := d.logger.With(
			logger.NewField("request_id", ctxHelper.GetRequestIdFromContext(ctx)),
			logger.NewField("method", "FindConfirmCode"),
			logger.NewField("duration", time.Since(start)),
		)

		if err != nil {
			l.Errorf("вызвал ошибку: %v", err)
			return
		}

		l.Info("вызов прошёл успешно")
	}()

	return d.base.FindConfirmCode(ctx, userId)
}

func (d *IUserRepositoryWithLogging) UpdateConfirmCode(ctx context.Context, code d.ConfirmCode) (err error) {
	start := time.Now()

	defer func() {
		l := d.logger.With(
			logger.NewField("request_id", ctxHelper.GetRequestIdFromContext(ctx)),
			logger.NewField("method", "UpdateConfirmCode"),
			logger.NewField("duration", time.Since(start)),
		)

		if err != nil {
			l.Errorf("вызвал ошибку: %v", err)
			return
		}

		l.Info("вызов прошёл успешно")
	}()

	return d.base.UpdateConfirmCode(ctx, code)
}

func (d *IUserRepositoryWithLogging) DeleteConfirmCode(ctx context.Context, userId int32) (err error) {
	start := time.Now()

	defer func() {
		l := d.logger.With(
			logger.NewField("request_id", ctxHelper.GetRequestIdFromContext(ctx)),
			logger.NewField("method", "DeleteConfirmCode"),
			logger.NewField("duration", time.Since(start)),
		)

		if err != nil {
			l.Errorf("вызвал ошибку: %v", err)
			return
		}

		l.Info("вызов прошёл успешно")
	}()

	return d.base.DeleteConfirmCode(ctx, userId)
}
