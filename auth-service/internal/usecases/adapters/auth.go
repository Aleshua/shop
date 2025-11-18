package adapters

import (
	"context"
	"time"

	"auth/internal/config"
	d "auth/internal/domain"
	l "auth/internal/logger"
	portsRepo "auth/internal/repositories/ports"
	portsSrv "auth/internal/services/ports"
	"auth/internal/usecases"
	"auth/internal/utils"
)

type AuthUseCases struct {
	logger l.ILogger

	emailConfirmCodeLength  int
	emailConfirmMaxAttempts int
	emailConfirmTTLSeconds  int

	jwtRefreshTokenTTLSeconds int64
	jwtRefreshTokenLength     int

	jwtSrv                portsSrv.IJWTService
	defaultHasherSrv      portsSrv.IHasherService
	passwordHasherSrv     portsSrv.IHasherService
	tokenGeneratorSrv     portsSrv.ITokenGeneratorService
	emailConfirmSenderSrv portsSrv.IEmailSenderService

	userRepo    portsRepo.IUserRepository
	sessionRepo portsRepo.ISessionRepository

	tx utils.Transaction
}

func NewAuthUseCases(
	logger l.ILogger,
	emailVerifyParams config.EmailVerify,
	jwtParams config.JWT,
	jwtSrv portsSrv.IJWTService,
	defaultHasherSrv portsSrv.IHasherService,
	passwordHasherSrv portsSrv.IHasherService,
	tokenGeneratorSrv portsSrv.ITokenGeneratorService,
	emailConfirmSenderSrv portsSrv.IEmailSenderService,
	userRepo portsRepo.IUserRepository,
	sessionRepo portsRepo.ISessionRepository,
	tx utils.Transaction,
) *AuthUseCases {
	return &AuthUseCases{
		logger: logger,

		emailConfirmCodeLength:  emailVerifyParams.CodeLength,
		emailConfirmMaxAttempts: emailVerifyParams.MaxAttempts,
		emailConfirmTTLSeconds:  emailVerifyParams.TTLSeconds,

		jwtRefreshTokenTTLSeconds: jwtParams.RefreshTokenTTLSeconds,
		jwtRefreshTokenLength:     jwtParams.RefreshTokenLength,

		jwtSrv:                jwtSrv,
		defaultHasherSrv:      defaultHasherSrv,
		passwordHasherSrv:     passwordHasherSrv,
		tokenGeneratorSrv:     tokenGeneratorSrv,
		emailConfirmSenderSrv: emailConfirmSenderSrv,

		userRepo:    userRepo,
		sessionRepo: sessionRepo,

		tx: tx,
	}
}

func (uc *AuthUseCases) Register(ctx context.Context, password, email string) (int32, error) {
	var userId int32
	err := uc.tx.WithinReadCommittedTx(ctx, func(ctx context.Context) error {
		passwordHash, err := uc.passwordHasherSrv.Hash(password)
		if err != nil {
			return err
		}

		userId, err = uc.userRepo.SaveOrUpdateUnverified(ctx, d.NewUser(email, passwordHash))
		if err != nil {
			return err
		}

		code, err := uc.tokenGeneratorSrv.GenerateToken(uc.emailConfirmCodeLength)
		if err != nil {
			return err
		}

		err = uc.userRepo.SaveOrUpdateConfirmCode(ctx, d.NewConfirmCode(userId, code))
		if err != nil {
			return err
		}

		go func() {
			err := uc.emailConfirmSenderSrv.Send(email, code)
			if err != nil {
				uc.logger.Errorf("не удалось отправить данные о коде подтвеждения пользователю с id: %d", userId)
			}
		}()

		return nil
	})

	return userId, err
}

func (uc *AuthUseCases) ConfirmEmail(ctx context.Context, userId int32, code string) error {
	confirmCode, err := uc.userRepo.FindConfirmCode(ctx, userId)
	if err != nil {
		return err
	}

	if confirmCode.Attempts+1 > uc.emailConfirmMaxAttempts {
		return usecases.ErrVerificationAttemptsExceeded
	}

	if confirmCode.CreatedAt.Unix()+int64(uc.emailConfirmTTLSeconds) < time.Now().Unix() {
		return usecases.ErrVerificationCodeExpired
	}

	if confirmCode.Code != code {

		confirmCode.Attempts += 1

		err = uc.userRepo.UpdateConfirmCode(ctx, confirmCode)
		if err != nil {
			uc.logger.Errorf("не удалось повысить счётчик поптыок подтверждения почты: %s", err.Error())
		}

		return usecases.ErrVerificationInvalidCode
	}

	err = uc.userRepo.MarkEmailVerified(ctx, userId)
	if err != nil {
		return err
	}

	go func() {
		err := uc.userRepo.DeleteConfirmCode(context.Background(), userId)
		if err != nil {
			uc.logger.Errorf(
				"не удалось удалить данные о коде подтвеждения пользователя с id: %d после его подтвеждения",
				userId,
			)
		}
	}()

	return nil
}

func (uc *AuthUseCases) ResendCode(ctx context.Context, userId int32) error {
	user, err := uc.userRepo.Find(ctx, userId)
	if err != nil {
		return nil
	}

	code, err := uc.tokenGeneratorSrv.GenerateToken(uc.emailConfirmCodeLength)
	if err != nil {
		return err
	}

	err = uc.userRepo.UpdateConfirmCode(ctx, d.NewConfirmCode(userId, code))
	if err != nil {
		return err
	}

	go func() {
		err := uc.emailConfirmSenderSrv.Send(user.Email, code)
		if err != nil {
			uc.logger.Errorf("не удалось отправить данные о коде подтвеждения пользователю с id: %d", userId)
		}
	}()

	return nil
}

func (uc *AuthUseCases) RefreshToken(ctx context.Context, refreshToken string) (string, error) {
	refreshTokenHash, err := uc.defaultHasherSrv.Hash(refreshToken)
	if err != nil {
		return "", err
	}

	session, err := uc.sessionRepo.Find(ctx, refreshTokenHash)
	if err != nil {
		return "", err
	}

	if session.Revoked {
		return "", usecases.ErrRefreshTokenRevoked
	}

	if session.ExpiresAt < time.Now().Unix() {
		return "", usecases.ErrRefreshTokenExpired
	}

	go func() {
		session.ExpiresAt = time.Now().Unix() + uc.jwtRefreshTokenTTLSeconds
		err = uc.sessionRepo.Update(ctx, session)
		if err != nil {
			uc.logger.Errorf("не удалось обновить время refresh токена: %s", err.Error())
		}
	}()

	accessToken, err := uc.jwtSrv.GenerateToken(ctx, session.UserId)
	if err != nil {
		return "", err
	}

	return accessToken, nil
}

func (uc *AuthUseCases) Login(ctx context.Context, password, email string) (string, string, error) {
	var (
		accessToken  string
		refreshToken string
	)

	err := uc.tx.WithinReadCommittedTx(ctx, func(ctx context.Context) error {
		user, err := uc.userRepo.FindByEmail(ctx, email)
		if err != nil {
			return err
		}

		if !uc.passwordHasherSrv.Compare(password, user.PasswordHash) {
			return usecases.ErrInvalidPassword
		}

		if user.EmailVerifiedAt == nil {
			return usecases.ErrEmailNotVerified
		}

		refreshToken, err = uc.tokenGeneratorSrv.GenerateToken(uc.jwtRefreshTokenLength)
		if err != nil {
			return err
		}

		refreshTokenHash, err := uc.defaultHasherSrv.Hash(refreshToken)
		if err != nil {
			return err
		}

		err = uc.sessionRepo.Save(ctx, d.NewSession(user.Id, refreshTokenHash, uc.jwtRefreshTokenTTLSeconds))
		if err != nil {
			return err
		}

		accessToken, err = uc.jwtSrv.GenerateToken(ctx, user.Id)
		if err != nil {
			return err
		}

		return nil
	})

	return accessToken, refreshToken, err
}

func (uc *AuthUseCases) Logout(ctx context.Context, refreshToken string) error {
	refreshTokenHash, err := uc.defaultHasherSrv.Hash(refreshToken)
	if err != nil {
		return err
	}

	session, err := uc.sessionRepo.Find(ctx, refreshTokenHash)
	if err != nil {
		return err
	}

	session.Revoked = true

	err = uc.sessionRepo.Update(ctx, session)
	if err != nil {
		return err
	}

	return nil
}
