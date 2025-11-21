package main

import (
	"crypto/rsa"
	"log"
	"os"

	c "auth/internal/config"
	decs "auth/internal/decorators"
	r "auth/internal/http/rest/router"
	reps "auth/internal/repositories"
	srvs "auth/internal/services"
	ucs "auth/internal/usecases"
	db "shared/database"
	pgh "shared/database/pghelper"
	l "shared/logger"

	"github.com/golang-jwt/jwt/v5"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
)

func main() {
	config := mustConfig()
	logger := mustLogger(config)
	database := mustDatabase(config)
	runMigrations(database)

	authUC := buildUsecases(config, database, logger)
	router := r.New(logger, authUC)

	router.Run(config.HTTP.Port)
}

func mustConfig() *c.Config {
	config, err := c.New()
	if err != nil {
		log.Fatalf("не удалось создать конфиг: %s", err.Error())
	}
	return config
}

func mustLogger(config *c.Config) l.ILogger {
	logger := l.New(config.Env, l.Options{
		Level:            config.Logger.Level,
		LogFilePath:      config.Logger.LogFilePath,
		BufferSize:       config.Logger.BufferSize,
		FlushInterval:    config.Logger.FlushInterval,
		FileMaxMegabytes: config.Logger.FileMaxMegabytes,
		MaxBackups:       config.Logger.MaxBackups,
		MaxAgeDays:       config.Logger.MaxAgeDays,
	})
	return logger
}

func mustDatabase(config *c.Config) *db.DB {
	database, err := db.New(db.Options{
		Path:              config.Database.Path,
		MaxConns:          config.Database.MaxConns,
		MinConns:          config.Database.MinConns,
		MaxConnLifetime:   config.Database.MaxConnLifetime,
		MaxConnIdleTime:   config.Database.MaxConnIdleTime,
		HealthCheckPeriod: config.Database.HealthCheckPeriod,
		ConnectTimeout:    config.Database.ConnectTimeout,
	})
	if err != nil {
		log.Fatalf("не удалось подключиться к базе данных: %s", err.Error())
	}
	return database
}

func runMigrations(database *db.DB) {
	if err := goose.Up(stdlib.OpenDBFromPool(database.Pool), "./migrations"); err != nil {
		log.Fatalf("миграции не применены: %s", err)
	}
}

func buildUsecases(config *c.Config, database *db.DB, logger l.ILogger) ucs.IAuthUseCases {
	privateKey, publicKey := mustLoadKeys()

	// Repos
	sessionRepo := decs.NewISessionRepositoryWithLogging(reps.NewSessionRepository(database.Pool), logger)
	userRepo := decs.NewIUserRepositoryWithLogging(reps.NewUserRepository(database.Pool), logger)

	// Services
	jwtService := decs.NewIJWTServiceWithLogging(srvs.NewJWT(config.JWT, privateKey, publicKey), logger)
	tokenGenerator := decs.NewITokenGeneratorServiceWithLogging(srvs.NewTokenGeneratorService(), logger)
	defaultHasher := decs.NewIHasherServiceWithLogging(srvs.NewDefaultHasherService(), logger)
	passwordHasher := decs.NewIHasherServiceWithLogging(srvs.NewPasswordHasherService(), logger)
	emailSender := decs.NewIEmailSenderServiceWithLogging(srvs.NewConfirmEmailCodeSenderService(config.Email), logger)

	// Transaction manager
	tx := pgh.NewTxManager(database.Pool)

	return decs.NewIAuthUseCasesWithLogging(ucs.NewAuthUseCases(
		logger,
		config.EmailVerify,
		config.JWT,
		jwtService,
		defaultHasher,
		passwordHasher,
		tokenGenerator,
		emailSender,
		userRepo,
		sessionRepo,
		tx,
	), logger)
}

func mustLoadKeys() (*rsa.PrivateKey, *rsa.PublicKey) {
	privateBytes, err := os.ReadFile("keys/private.pem")
	if err != nil {
		log.Fatal("не удалось прочитать private.pem:", err)
	}

	privateKey, err := jwt.ParseRSAPrivateKeyFromPEM(privateBytes)
	if err != nil {
		log.Fatal("ошибка парсинга приватного ключа:", err)
	}

	publicBytes, err := os.ReadFile("keys/public.pem")
	if err != nil {
		log.Fatal("не удалось прочитать public.pem:", err)
	}

	publicKey, err := jwt.ParseRSAPublicKeyFromPEM(publicBytes)
	if err != nil {
		log.Fatal("ошибка парсинга публичного ключа:", err)
	}

	return privateKey, publicKey
}
