package slogx

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"strings"
	"time"

	"github.com/lmittmann/tint"
)

type contextKey struct{}

var loggerKey = contextKey{}

const (
	// FormatText is the text output format for logs
	FormatText = "text"

	// FormatJSON is the JSON output format for logs
	FormatJSON = "json"

	// FormatPretty is a prettier output format for logs
	FormatPretty = "pretty"

	// DefaultLogLevel is the default logging level
	DefaultLogLevel = slog.LevelInfo
)

const LevelTrace = slog.LevelDebug - 4

const (
	// Log level string representations
	traceLevel = "trace"
	debugLevel = "debug"
	infoLevel  = "info"
	warnLevel  = "warn"
	errorLevel = "error"
)

// Logger wraps slog.Logger with helper methods
type Logger struct {
	*slog.Logger
}

// With returns a new logger with additional context fields
func (l *Logger) With(args ...any) *Logger {
	return &Logger{l.Logger.With(args...)}
}

// Panic logs the error and exits the application
func (l *Logger) Panic(err error) {
	l.Error(err.Error())
	os.Exit(1)
}

// Trace logs at [LevelTrace].
func (l *Logger) Trace(msg string, args ...any) {
	l.Log(context.Background(), LevelTrace, msg, args...)
}

// TraceContext logs at [LevelTrace] with the given context.
func (l *Logger) TraceContext(ctx context.Context, msg string, args ...any) {
	l.Log(ctx, LevelTrace, msg, args...)
}

// Default returns the default logger instance
func Default() *Logger {
	return &Logger{slog.Default()}
}

// NewWithContext initializes a new slog and attaches it to the provided context
func NewWithContext(ctx context.Context) (context.Context, *Logger, error) {
	logger, err := New()
	if err != nil {
		return nil, nil, err
	}
	ctx = WithContext(ctx, logger)
	return ctx, logger, nil
}

// New initialize a new slog and set up configuration with logger level and output.
//
// If config level is empty, will be set to DefaultLogLevel
func New() (*Logger, error) {
	logLevel := os.Getenv("LOG_LEVEL")
	logFormat := os.Getenv("LOG_FORMAT")

	parsedLevel, err := parseLevel(logLevel)
	if err != nil {
		return nil, err
	}

	opts := &slog.HandlerOptions{
		Level:       parsedLevel,
		ReplaceAttr: replaceAttr,
	}

	var handler slog.Handler
	switch logFormat {
	case FormatText:
		handler = slog.NewTextHandler(os.Stdout, opts)
	case FormatJSON:
		handler = slog.NewJSONHandler(os.Stdout, opts)
	case FormatPretty:
		handler = tint.NewHandler(os.Stdout, &tint.Options{
			Level:       parsedLevel,
			TimeFormat:  time.DateTime,
			ReplaceAttr: replaceAttr,
		})
	default:
		return nil, fmt.Errorf("invalid log format: %v", logFormat)
	}

	slogger := slog.New(handler)
	slog.SetDefault(slogger)

	logger := &Logger{
		slogger,
	}
	logger.Debug("log configured", "level", parsedLevel, "format", logFormat)
	return logger, nil
}

// WithContext attaches a Logger to the context
func WithContext(ctx context.Context, log *Logger) context.Context {
	return context.WithValue(ctx, loggerKey, log)
}

// FromCtx retrieves the Logger from the context or returns default
func FromCtx(ctx context.Context) *Logger {
	if log, ok := ctx.Value(loggerKey).(*Logger); ok {
		return log
	}
	return Default()
}

// With creates a new logger with the given context fields
func With(args ...any) *Logger {
	return &Logger{slog.With(args...)}
}

func parseLevel(level string) (slog.Level, error) {
	level = strings.ToLower(level)
	switch level {
	case traceLevel:
		return LevelTrace, nil
	case debugLevel:
		return slog.LevelDebug, nil
	case infoLevel:
		return slog.LevelInfo, nil
	case warnLevel:
		return slog.LevelWarn, nil
	case errorLevel:
		return slog.LevelError, nil
	default:
		return DefaultLogLevel, errors.New("invalid log level")
	}
}

func replaceAttr(groups []string, a slog.Attr) slog.Attr {
	if a.Key == slog.LevelKey && len(groups) == 0 {
		level, ok := a.Value.Any().(slog.Level)
		if ok {
			switch level {
			case LevelTrace:
				return tint.Attr(8, slog.String(a.Key, "TRACE"))
			case slog.LevelDebug:
				return tint.Attr(7, slog.String(a.Key, "DEBUG"))
			case slog.LevelInfo:
				return tint.Attr(2, slog.String(a.Key, "INFO"))
			case slog.LevelWarn:
				return tint.Attr(11, slog.String(a.Key, "WARN"))
			case slog.LevelError:
				return tint.Attr(9, slog.String(a.Key, "ERROR"))
			}
		}
	}
	return a
}
