package logger

type Field struct {
	Key   string
	Value interface{}
}

type ILogger interface {
	// Flushes buffer, if any
	Sync()

	// Message at the Debug level
	Debug(msg ...interface{})

	// Debugf formats according to the format specifier
	// and logs the message at the Debug level
	Debugf(msg string, args ...interface{})

	// Message at the Info level
	Info(msg ...interface{})

	// Infof formats according to the format specifier
	// and logs the message at the Info level
	Infof(msg string, args ...interface{})

	// Message at the Warn level
	Warn(msg ...interface{})

	// Warnf formats according to the format specifier
	// and logs the message at the Warn level
	Warnf(msg string, args ...interface{})

	// Message at the Error level
	Error(msg ...interface{})

	// Errorf formats according to the format specifier
	// and logs the message at the Error level
	Errorf(msg string, args ...interface{})

	// Message at the Panic level and panics
	Panic(msg ...interface{})

	// Panicf formats according to the format specifier
	// and logs the message at the Panic level and panics
	Panicf(msg string, args ...interface{})

	// Message at the Fatal level  and exists the application
	Fatal(msg ...interface{})

	// Fatalf formats according to the format specifier
	// and logs the message at the Fatal level  and exists the application
	Fatalf(msg string, args ...interface{})

	// With additional context (key-value pairs) to the logger
	With(args ...Field) ILogger
}

func NewField(key string, value interface{}) Field {
	return Field{Key: key, Value: value}
}
