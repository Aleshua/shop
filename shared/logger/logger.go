package logger

import (
	"fmt"
	"os"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

type Options struct {
	Level            int
	LogFilePath      string
	BufferSize       int
	FlushInterval    time.Duration
	FileMaxMegabytes int
	MaxBackups       int
	MaxAgeDays       int
}

type Logger struct {
	logger *zap.Logger
}

// env ~ (local, test, dev, prod) .
// Если env ~ (local, test) то все параметры кроме Level игнорируются
func New(env string, options Options) *Logger {
	var zapcoreLevel zapcore.Level

	switch options.Level {
	case -1:
		zapcoreLevel = zapcore.DebugLevel
	case 0:
		zapcoreLevel = zapcore.InfoLevel
	case 1:
		zapcoreLevel = zapcore.WarnLevel
	case 2:
		zapcoreLevel = zap.ErrorLevel
	default:
		zapcoreLevel = zapcore.InfoLevel
	}

	var core zapcore.Core

	if env == "local" || env == "test" {
		core = consoleCore(zapcoreLevel)
	} else {
		core = fileCore(zapcoreLevel, options)
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

func fileCore(level zapcore.Level, options Options) zapcore.Core {
	ws := &zapcore.BufferedWriteSyncer{
		WS: zapcore.AddSync(
			&lumberjack.Logger{
				Filename:   "./logs/app.log",
				MaxSize:    options.FileMaxMegabytes,
				MaxBackups: options.MaxBackups,
				MaxAge:     options.MaxAgeDays,
				Compress:   false,
				LocalTime:  true,
			},
		),
		Size:          options.BufferSize,
		FlushInterval: options.FlushInterval,
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
