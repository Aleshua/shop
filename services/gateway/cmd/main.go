package main

import (
	"log"
	"shared/logger"

	"gateway/internal/config"
	"gateway/internal/http/rest/router"
)

func main() {
	config := mustConfig()
	logger := mustLogger(config)

	router := router.New(config.Services, logger)

	router.Run(config.Http.Port)
}

func mustConfig() *config.Config {
	config, err := config.New()
	if err != nil {
		log.Fatalf("не удалось создать конфиг: %s", err.Error())
	}
	return config
}

func mustLogger(config *config.Config) logger.ILogger {
	logger := logger.New(config.Env, logger.Options{
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
