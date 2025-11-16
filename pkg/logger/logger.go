package logger

import (
	"fmt"
	"strings"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const consoleFormat = "console"

// Interface -.
type Interface interface {
	Debug(message interface{}, args ...interface{})
	Info(message string, args ...interface{})
	Warn(message string, args ...interface{})
	Error(message interface{}, args ...interface{})
	Fatal(message interface{}, args ...interface{})
}

// Logger -.
type Logger struct {
	logger *zap.SugaredLogger
}

var _ Interface = (*Logger)(nil)

// New creates a logger with the specified level and format.
// format can be "console" for pretty printing or "json" for production.
func New(level, format string) *Logger {
	var zapLevel zapcore.Level

	switch strings.ToLower(level) {
	case "error":
		zapLevel = zapcore.ErrorLevel
	case "warn":
		zapLevel = zapcore.WarnLevel
	case "info":
		zapLevel = zapcore.InfoLevel
	case "debug":
		zapLevel = zapcore.DebugLevel
	default:
		zapLevel = zapcore.InfoLevel
	}

	config := zap.NewProductionConfig()
	config.Level = zap.NewAtomicLevelAt(zapLevel)
	config.EncoderConfig.TimeKey = "timestamp"
	config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	config.EncoderConfig.CallerKey = "caller"
	config.EncoderConfig.EncodeCaller = zapcore.ShortCallerEncoder
	if strings.ToLower(strings.TrimSpace(format)) == consoleFormat {
		config.Encoding = consoleFormat
		config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	}

	logger, _ := config.Build(zap.AddCallerSkip(1))
	sugar := logger.Sugar()

	return &Logger{
		logger: sugar,
	}
}

// Debug -.
func (l *Logger) Debug(message interface{}, args ...interface{}) {
	l.msg(l.logger.Debugw, message, args...)
}

// Info -.
func (l *Logger) Info(message string, args ...interface{}) {
	if len(args) == 0 {
		l.logger.Info(message)
	} else {
		l.logger.Infow(message, args...)
	}
}

// Warn -.
func (l *Logger) Warn(message string, args ...interface{}) {
	if len(args) == 0 {
		l.logger.Warn(message)
	} else {
		l.logger.Warnw(message, args...)
	}
}

// Error -.
func (l *Logger) Error(message interface{}, args ...interface{}) {
	l.msg(l.logger.Errorw, message, args...)
}

// Fatal -.
func (l *Logger) Fatal(message interface{}, args ...interface{}) {
	l.msg(l.logger.Fatalw, message, args...)
}

func (l *Logger) msg(logFunc func(string, ...interface{}), message interface{}, args ...interface{}) {
	switch msg := message.(type) {
	case error:
		logFunc(msg.Error(), args...)
	case string:
		logFunc(msg, args...)
	default:
		logFunc(fmt.Sprintf("%v", message), args...)
	}
}
