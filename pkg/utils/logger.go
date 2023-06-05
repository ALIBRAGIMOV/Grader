package utils

import (
	"context"
	"go.uber.org/zap"
	"grader/pkg/server/logger"
)

var (
	defaultLogger *zap.SugaredLogger
)

func Init(logger *zap.SugaredLogger) {
	defaultLogger = logger
}

func GetLogger(ctx context.Context) *zap.SugaredLogger {
	if ctx == nil {
		return defaultLogger
	}

	z, ok := ctx.Value(logger.Key).(*zap.SugaredLogger)

	if !ok || z == nil {
		return defaultLogger
	}

	return z
}

func GetRequestIDFromContext(ctx context.Context) string {
	requestID, ok := ctx.Value(logger.RequestIDKey).(string)

	if !ok {
		return "-"
	}
	return requestID
}
