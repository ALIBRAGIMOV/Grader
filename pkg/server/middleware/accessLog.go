package middleware

import (
	"grader/pkg/utils"
	"net/http"
	"time"
)

func AccessLog(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)

		utils.GetLogger(r.Context()).Infow(r.URL.Path,
			"method", r.Method,
			"remote_addr", r.RemoteAddr,
			"url", r.URL.Path,
			"work_time", time.Since(start),
		)
	})
}
