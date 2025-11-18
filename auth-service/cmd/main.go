package main

import (
	"crypto/rsa"
	"log"
	"os"

	c "auth/internal/config"
	db "auth/internal/database"
	r "auth/internal/http/rest/router"
	l "auth/internal/logger"
	repoAdapters "auth/internal/repositories/adapters"
	srvAdapters "auth/internal/services/adapters"
	srvMocks "auth/internal/services/mocks"
	ucAdapters "auth/internal/usecases/adapters"
	"auth/internal/utils"

	"github.com/golang-jwt/jwt/v5"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
)

func main() {
	config := mustConfig()
	database := mustDatabase(config)
	runMigrations(database)
	logger := l.New(config.Env, config.Logger)

	authUC := buildUsecases(config, database, logger)
	router := r.New(config.HTTP, logger, authUC)

	router.Run()
}

func mustConfig() *c.Config {
	config, err := c.New()
	if err != nil {
		log.Fatalf("не удалось создать конфиг: %s", err.Error())
	}
	return config
}

func mustDatabase(config *c.Config) *db.DB {
	database, err := db.New(config.Database)
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

func buildUsecases(config *c.Config, database *db.DB, logger l.ILogger) *ucAdapters.AuthUseCases {
	switch config.Env {
	case "test":
		return buildTestUsecases(config, database, logger)
	default:
		return buildProdUsecases(config, database, logger)
	}
}

func buildProdUsecases(config *c.Config, database *db.DB, logger l.ILogger) *ucAdapters.AuthUseCases {
	privateKey, publicKey := mustLoadKeys()

	// Repos
	sessionRepo := repoAdapters.NewSessionRepository(database.Pool)
	userRepo := repoAdapters.NewUserRepository(database.Pool)

	// Services
	jwtService := srvAdapters.NewJWT(config.JWT, privateKey, publicKey)
	tokenGenerator := srvAdapters.NewTokenGeneratorService()
	defaultHasher := srvAdapters.NewDefaultHasherService()
	passwordHasher := srvAdapters.NewPasswordHasherService()
	emailSender := srvAdapters.NewConfirmEmailCodeSenderService(config.Email)

	// Transaction manager
	tx := utils.NewTxManager(database.Pool)

	return ucAdapters.NewAuthUseCases(
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
	)
}

func buildTestUsecases(config *c.Config, database *db.DB, logger l.ILogger) *ucAdapters.AuthUseCases {
	privateKey, publicKey := mustLoadKeys()

	// Repos
	sessionRepo := repoAdapters.NewSessionRepository(database.Pool)
	userRepo := repoAdapters.NewUserRepository(database.Pool)

	// Services
	jwtService := srvAdapters.NewJWT(config.JWT, privateKey, publicKey)
	defaultHasher := srvAdapters.NewDefaultHasherService()
	passwordHasher := srvAdapters.NewPasswordHasherService()

	// В тестах — моковый генератор токенов
	tokenGenerator := srvMocks.NewPredictableTokenGeneratorService()

	// В тестах — моковый email sender
	emailSender := srvMocks.NewMockEmailSender()

	// Transaction manager
	tx := utils.NewTxManager(database.Pool)

	return ucAdapters.NewAuthUseCases(
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
	)
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
