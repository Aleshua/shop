package decorators

import (
	"context"
	"shared/logger"
	"time"

	d "auth/internal/domain"
	reps "auth/internal/repositories"

	ctxHelper "shared/ctxhelper"
)

type ISessionRepositoryWithLogging struct {
	base   reps.ISessionRepository
	logger logger.ILogger
}

func NewISessionRepositoryWithLogging(
	base reps.ISessionRepository,
	logger logger.ILogger,
) *ISessionRepositoryWithLogging {
	return &ISessionRepositoryWithLogging{
		base:   base,
		logger: logger,
	}
}

func (d *ISessionRepositoryWithLogging) Find(ctx context.Context, refreshTokenHash string) (_ d.Session, err error) {
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

	return d.base.Find(ctx, refreshTokenHash)
}

func (d *ISessionRepositoryWithLogging) Update(ctx context.Context, session d.Session) (err error) {
	start := time.Now()

	defer func() {
		l := d.logger.With(
			logger.NewField("request_id", ctxHelper.GetRequestIdFromContext(ctx)),
			logger.NewField("method", "Update"),
			logger.NewField("duration", time.Since(start)),
		)

		if err != nil {
			l.Errorf("вызвал ошибку: %v", err)
			return
		}

		l.Info("вызов прошёл успешно")
	}()

	return d.base.Update(ctx, session)
}

func (d *ISessionRepositoryWithLogging) Save(ctx context.Context, session d.Session) (err error) {
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

	return d.base.Save(ctx, session)
}
