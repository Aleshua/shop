package usecases

import (
	"context"
	"time"

	"auth/internal/config"
	d "auth/internal/domain"
	reps "auth/internal/repositories"
	srvs "auth/internal/services"
	pgh "shared/database/pghelper"
	l "shared/logger"
)

type AuthUseCases struct {
	logger l.ILogger

	emailConfirmCodeLength  int
	emailConfirmMaxAttempts int
	emailConfirmTTLSeconds  int

	jwtRefreshTokenTTLSeconds int64
	jwtRefreshTokenLength     int

	jwtSrv                srvs.IJWTService
	defaultHasherSrv      srvs.IHasherService
	passwordHasherSrv     srvs.IHasherService
	tokenGeneratorSrv     srvs.ITokenGeneratorService
	emailConfirmSenderSrv srvs.IEmailSenderService

	userRepo    reps.IUserRepository
	sessionRepo reps.ISessionRepository

	tx pgh.Transaction
}

func NewAuthUseCases(
	logger l.ILogger,
	emailVerifyParams config.EmailVerify,
	jwtParams config.JWT,
	jwtSrv srvs.IJWTService,
	defaultHasherSrv srvs.IHasherService,
	passwordHasherSrv srvs.IHasherService,
	tokenGeneratorSrv srvs.ITokenGeneratorService,
	emailConfirmSenderSrv srvs.IEmailSenderService,
	userRepo reps.IUserRepository,
	sessionRepo reps.ISessionRepository,
	tx pgh.Transaction,
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
		passwordHash, err := uc.passwordHasherSrv.Hash(ctx, password)
		if err != nil {
			return err
		}

		userId, err = uc.userRepo.SaveOrUpdateUnverified(ctx, d.NewUser(email, passwordHash))
		if err != nil {
			return err
		}

		code, err := uc.tokenGeneratorSrv.GenerateToken(ctx, uc.emailConfirmCodeLength)
		if err != nil {
			return err
		}

		err = uc.userRepo.SaveOrUpdateConfirmCode(ctx, d.NewConfirmCode(userId, code))
		if err != nil {
			return err
		}

		go uc.emailConfirmSenderSrv.Send(ctx, email, code)

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
		return d.ErrVerificationAttemptsExceeded
	}

	if confirmCode.CreatedAt.Unix()+int64(uc.emailConfirmTTLSeconds) < time.Now().Unix() {
		return d.ErrVerificationCodeExpired
	}

	if confirmCode.Code != code {

		confirmCode.Attempts += 1

		uc.userRepo.UpdateConfirmCode(ctx, confirmCode)

		return d.ErrVerificationInvalidCode
	}

	err = uc.userRepo.MarkEmailVerified(ctx, userId)
	if err != nil {
		return err
	}

	go uc.userRepo.DeleteConfirmCode(context.Background(), userId)

	return nil
}

func (uc *AuthUseCases) ResendCode(ctx context.Context, userId int32) error {
	user, err := uc.userRepo.Find(ctx, userId)
	if err != nil {
		return nil
	}

	code, err := uc.tokenGeneratorSrv.GenerateToken(ctx, uc.emailConfirmCodeLength)
	if err != nil {
		return err
	}

	err = uc.userRepo.UpdateConfirmCode(ctx, d.NewConfirmCode(userId, code))
	if err != nil {
		return err
	}

	go uc.emailConfirmSenderSrv.Send(ctx, user.Email, code)

	return nil
}

func (uc *AuthUseCases) RefreshToken(ctx context.Context, refreshToken string) (string, error) {
	refreshTokenHash, err := uc.defaultHasherSrv.Hash(ctx, refreshToken)
	if err != nil {
		return "", err
	}

	session, err := uc.sessionRepo.Find(ctx, refreshTokenHash)
	if err != nil {
		return "", err
	}

	if session.Revoked {
		return "", d.ErrRefreshTokenRevoked
	}

	if session.ExpiresAt < time.Now().Unix() {
		return "", d.ErrRefreshTokenExpired
	}

	go func() {
		session.ExpiresAt = time.Now().Unix() + uc.jwtRefreshTokenTTLSeconds
		uc.sessionRepo.Update(ctx, session)
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

		if !uc.passwordHasherSrv.Compare(ctx, password, user.PasswordHash) {
			return d.ErrInvalidPassword
		}

		if user.EmailVerifiedAt == nil {
			return d.ErrEmailNotVerified
		}

		refreshToken, err = uc.tokenGeneratorSrv.GenerateToken(ctx, uc.jwtRefreshTokenLength)
		if err != nil {
			return err
		}

		refreshTokenHash, err := uc.defaultHasherSrv.Hash(ctx, refreshToken)
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
	refreshTokenHash, err := uc.defaultHasherSrv.Hash(ctx, refreshToken)
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
