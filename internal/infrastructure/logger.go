package infrastructure

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/hydr0g3nz/mini_bank/internal/domain/infra"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// LoggerConfig holds configuration for logger
type LoggerConfig struct {
	IsProduction bool
	EnableFile   bool   // Optional file logging
	LogDir       string // Optional custom log directory
}

// Logger implements the AppLogger interface using zap
type Logger struct {
	zap *zap.Logger
}

// NewLogger creates a new logger instance with optional file logging
func NewLogger(config LoggerConfig) (*Logger, error) {
	var zapConfig zap.Config

	if config.IsProduction {
		zapConfig = zap.NewProductionConfig()
		zapConfig.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	} else {
		zapConfig = zap.NewDevelopmentConfig()
		zapConfig.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	}

	// Create cores for logging
	cores := []zapcore.Core{}

	// Add file core only if enabled
	if config.EnableFile {
		fileCore, err := createFileCore(zapConfig, config.LogDir)
		if err != nil {
			// Log warning but continue without file logging
			fmt.Fprintf(os.Stderr, "Warning: Failed to create file logger: %v. Continuing with console logging only.\n", err)
		} else {
			cores = append(cores, fileCore)
		}
	}

	// Always add stdout core
	stdoutCore := zapcore.NewCore(
		zapcore.NewConsoleEncoder(zapConfig.EncoderConfig),
		zapcore.AddSync(os.Stdout),
		zapConfig.Level,
	)
	cores = append(cores, stdoutCore)

	// Combine cores
	core := zapcore.NewTee(cores...)

	// Create logger with caller skip
	zapLogger := zap.New(core, zap.AddCallerSkip(1))

	return &Logger{zapLogger}, nil
}

// NewSimpleLogger creates a logger with console output only (no file logging)
func NewSimpleLogger(isProduction bool) (*Logger, error) {
	return NewLogger(LoggerConfig{
		IsProduction: isProduction,
		EnableFile:   false, // ไม่สร้างไฟล์ log
	})
}

// createFileCore creates a file-based logging core
func createFileCore(config zap.Config, logDir string) (zapcore.Core, error) {
	// Use default log directory if not specified
	if logDir == "" {
		logDir = "logs"
	}

	// Create log directory if it doesn't exist
	if err := os.MkdirAll(logDir, 0755); err != nil {
		return nil, fmt.Errorf("can't create log directory: %w", err)
	}

	// Set up log file path
	logFile := filepath.Join(logDir, fmt.Sprintf("app-%s.log", time.Now().Format("2006-01-02")))

	// Open log file
	file, err := os.OpenFile(logFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, fmt.Errorf("failed to open log file %s: %w", logFile, err)
	}

	// Create file core
	fileCore := zapcore.NewCore(
		zapcore.NewJSONEncoder(config.EncoderConfig),
		zapcore.AddSync(file),
		config.Level,
	)

	return fileCore, nil
}

// Implement the AppLogger methods
func (l *Logger) Debug(msg string, fields ...interface{}) {
	l.zap.Debug(msg, toZapFields(fields...)...)
}

func (l *Logger) Debugf(format string, args ...interface{}) {
	l.zap.Sugar().Debugf(format, args...)
}

func (l *Logger) Info(msg string, fields ...interface{}) {
	l.zap.Info(msg, toZapFields(fields...)...)
}

func (l *Logger) Infof(format string, args ...interface{}) {
	l.zap.Sugar().Infof(format, args...)
}

func (l *Logger) Warn(msg string, fields ...interface{}) {
	l.zap.Warn(msg, toZapFields(fields...)...)
}

func (l *Logger) Warnf(format string, args ...interface{}) {
	l.zap.Sugar().Warnf(format, args...)
}

func (l *Logger) Error(msg string, fields ...interface{}) {
	l.zap.Error(msg, toZapFields(fields...)...)
}

func (l *Logger) Errorf(format string, args ...interface{}) {
	l.zap.Sugar().Errorf(format, args...)
}

func (l *Logger) Fatal(msg string, fields ...interface{}) {
	l.zap.Fatal(msg, toZapFields(fields...)...)
}

func (l *Logger) Fatalf(format string, args ...interface{}) {
	l.zap.Sugar().Fatalf(format, args...)
}

func (l *Logger) Sync() error {
	return l.zap.Sync()
}

func (l *Logger) With(fields ...interface{}) infra.Logger {
	return &Logger{
		zap: l.zap.With(toZapFields(fields...)...),
	}
}

func (l *Logger) Close() error {
	return l.zap.Sync()
}

func toZapFields(fields ...interface{}) []zapcore.Field {
	zapFields := make([]zapcore.Field, 0, len(fields)/2)

	for i := 0; i < len(fields)-1; i += 2 {
		key, ok := fields[i].(string)
		if !ok {
			continue // ข้ามถ้า key ไม่ใช่ string
		}
		zapFields = append(zapFields, zap.Any(key, fields[i+1]))
	}

	return zapFields
}
