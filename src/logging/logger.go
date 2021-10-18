package logging

import (
	"context"
	"fmt"
	"os"
	"strings"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var logger *zap.SugaredLogger

var logLevelSeverity = map[string]zapcore.Level{
	"DEBUG":     zapcore.DebugLevel,
	"INFO":      zapcore.InfoLevel,
	"WARNING":   zapcore.WarnLevel,
	"ERROR":     zapcore.ErrorLevel,
	"CRITICAL":  zapcore.DPanicLevel,
	"ALERT":     zapcore.PanicLevel,
	"EMERGENCY": zapcore.FatalLevel,
}

type requestIDType int

const (
	requestIDKey requestIDType = iota
)

// WithRqID returns a context with requestID
func WithRqID(ctx context.Context, requestID string) context.Context {
	return context.WithValue(ctx, requestIDKey, requestID)
}

// NewLogger returns a zap logger
func NewLogger(isDev bool) *zap.SugaredLogger {
	logLevel := strings.ToUpper(os.Getenv("LOG_LEVEL"))
	if logLevel == "" {
		logLevel = "INFO"
	}
	var defaultLogger *zap.SugaredLogger

	if isDev {
		config := zap.NewDevelopmentConfig()
		l, _ := config.Build()
		defaultLogger = l.Sugar()
	} else {
		config := zap.NewProductionEncoderConfig()
		encoder := zapcore.NewJSONEncoder(config)
		atom := zap.NewAtomicLevel()
		defaultLogger = zap.New(zapcore.NewCore(encoder, zapcore.Lock(os.Stdout), atom)).Sugar()
	}
	// flushes buffer
	defer func() {
		err := defaultLogger.Sync()
		fmt.Println(err)
	}()

	return defaultLogger //. .With(zap.String("v", "0.0.1"))
}
