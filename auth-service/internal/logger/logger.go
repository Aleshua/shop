package logger

import (
	"fmt"
	"os"

	c "auth/internal/config"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

type Logger struct {
	logger *zap.Logger
}

func New(env string, params c.Logger) *Logger {
	var level zapcore.Level

	switch params.Level {
	case -1:
		level = zapcore.DebugLevel
	case 0:
		level = zapcore.InfoLevel
	case 1:
		level = zapcore.WarnLevel
	case 2:
		level = zap.ErrorLevel
	default:
		level = zapcore.InfoLevel
	}

	var core zapcore.Core

	switch env {
	case "local":
		core = consoleCore(level)
	case "test":
		core = consoleCore(level)
		// core = testCore(level, params)
	default:
		core = fileCore(level, params)
	}

	logger := zap.New(
		core,
		zap.AddCaller(),
		zap.AddCallerSkip(1),
	)

	return &Logger{logger: logger}
}

func consoleCore(level zapcore.Level) zapcore.Core {
	ws := zapcore.AddSync(os.Stdout)

	encoderConfig := zapcore.EncoderConfig{
		NameKey: "zap",

		MessageKey:  "msg",
		LevelKey:    "level",
		TimeKey:     "timestamp",
		CallerKey:   "caller",
		FunctionKey: "function",

		SkipLineEnding: false,
		EncodeCaller:   zapcore.ShortCallerEncoder,
		EncodeLevel:    zapcore.CapitalColorLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
	}

	encoder := zapcore.NewConsoleEncoder(encoderConfig)

	return zapcore.NewCore(encoder, ws, level)
}

func fileCore(level zapcore.Level, params c.Logger) zapcore.Core {
	ws := &zapcore.BufferedWriteSyncer{
		WS: zapcore.AddSync(
			&lumberjack.Logger{
				Filename:   "./logs/app.log",
				MaxSize:    params.FileMaxMegabytes,
				MaxBackups: params.MaxBackups,
				MaxAge:     params.MaxAgeDays,
				Compress:   false,
				LocalTime:  true,
			},
		),
		Size:          params.BufferSize,
		FlushInterval: params.FlushInterval,
	}

	encoderConfig := zapcore.EncoderConfig{
		NameKey: "zap",

		MessageKey:  "msg",
		LevelKey:    "level",
		TimeKey:     "timestamp",
		CallerKey:   "caller",
		FunctionKey: "function",

		SkipLineEnding: false,
		EncodeCaller:   zapcore.ShortCallerEncoder,
		EncodeLevel:    zapcore.CapitalLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
	}

	encoder := zapcore.NewJSONEncoder(encoderConfig)

	return zapcore.NewCore(encoder, ws, level)
}

func testCore(level zapcore.Level, params c.Logger) zapcore.Core {
	return zapcore.NewNopCore()
}

func (l *Logger) Sync() {
	l.logger.Sync()
}

func (l *Logger) Debug(msg ...interface{}) {
	if ce := l.logger.Check(zap.DebugLevel, fmt.Sprint(msg...)); ce != nil {
		ce.Write()
	}
}

func (l *Logger) Debugf(msg string, args ...interface{}) {
	if ce := l.logger.Check(zap.DebugLevel, fmt.Sprintf(msg, args...)); ce != nil {
		ce.Write()
	}
}

func (l *Logger) Info(msg ...interface{}) {
	if ce := l.logger.Check(zap.InfoLevel, fmt.Sprint(msg...)); ce != nil {
		ce.Write()
	}
}

func (l *Logger) Infof(msg string, args ...interface{}) {
	if ce := l.logger.Check(zap.InfoLevel, fmt.Sprintf(msg, args...)); ce != nil {
		ce.Write()
	}
}

func (l *Logger) Warn(msg ...interface{}) {
	if ce := l.logger.Check(zap.WarnLevel, fmt.Sprint(msg...)); ce != nil {
		ce.Write()
	}
}

func (l *Logger) Warnf(msg string, args ...interface{}) {
	if ce := l.logger.Check(zap.WarnLevel, fmt.Sprintf(msg, args...)); ce != nil {
		ce.Write()
	}
}

func (l *Logger) Error(msg ...interface{}) {
	if ce := l.logger.Check(zap.ErrorLevel, fmt.Sprint(msg...)); ce != nil {
		ce.Write()
	}
}

func (l *Logger) Errorf(msg string, args ...interface{}) {
	if ce := l.logger.Check(zap.ErrorLevel, fmt.Sprintf(msg, args...)); ce != nil {
		ce.Write()
	}
}

func (l *Logger) Panic(msg ...interface{}) {
	if ce := l.logger.Check(zap.PanicLevel, fmt.Sprint(msg...)); ce != nil {
		ce.Write()
	}
}

func (l *Logger) Panicf(msg string, args ...interface{}) {
	if ce := l.logger.Check(zap.PanicLevel, fmt.Sprintf(msg, args...)); ce != nil {
		ce.Write()
	}
}

func (l *Logger) Fatal(msg ...interface{}) {
	if ce := l.logger.Check(zap.FatalLevel, fmt.Sprint(msg...)); ce != nil {
		ce.Write()
	}
}

func (l *Logger) Fatalf(msg string, args ...interface{}) {
	if ce := l.logger.Check(zap.FatalLevel, fmt.Sprintf(msg, args...)); ce != nil {
		ce.Write()
	}
}

func (l *Logger) With(args ...Field) ILogger {
	return &Logger{logger: l.logger.With(fieldToAny(args)...)}
}

func fieldToAny(fields []Field) []zapcore.Field {
	var result []zapcore.Field

	for _, f := range fields {
		result = append(result, zap.Any(f.Key, f.Value))
	}

	return result
}
