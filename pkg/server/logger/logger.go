package logger

import "go.uber.org/zap"

const (
	RequestIDKey = "requestID"
	Key          = "logger"
)

type Logger struct {
	Zap   *zap.Logger
	Level int
}
