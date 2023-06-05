package middleware

import (
	"context"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"grader/pkg/server/logger"
	"grader/pkg/utils"
	"net/http"
)

func getMinLogLevel(r *http.Request) zapcore.Level {
	minLevel := zap.ErrorLevel
	if r.FormValue("level") == "debug" {
		minLevel = zap.DebugLevel
	}
	return minLevel
}

func Logger(l *logger.Logger, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		ctxlogger := l.Zap.With(
			zap.String("logger", "ctxlog"),
			zap.String("trace-id", utils.GetRequestIDFromContext(r.Context())),
			zap.String("request-method", r.Method),
			zap.String("request-url", r.URL.String()),
			zap.String("remote-address", r.RemoteAddr),
			zap.String("email", "user@mail.ru"),
		).WithOptions(
			zap.IncreaseLevel(getMinLogLevel(r)),
			zap.AddCaller(),
			zap.AddStacktrace(zap.ErrorLevel),
		).Sugar()

		ctx := context.WithValue(r.Context(), logger.Key, ctxlogger)

		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	})
}
