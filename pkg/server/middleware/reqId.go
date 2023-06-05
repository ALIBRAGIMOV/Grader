package middleware

import (
	"context"
	"crypto/rand"
	"fmt"
	"grader/pkg/server/logger"
	"net/http"
)

func ReqID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestID := r.Header.Get("X-Request-ID")
		if requestID == "" {
			requestID = randBytesHex(16)
			r.Header.Set("X-Request-ID", requestID)
			r.Header.Set("trace-id", requestID)
			w.Header().Set("trace-id", requestID)
			w.Header().Set("X-Request-ID", requestID)
		}

		ctx := context.WithValue(r.Context(), logger.RequestIDKey, requestID)
		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)
	})
}

func randBytesHex(n int) string {
	return fmt.Sprintf("%x", randBytes(n))
}

func randBytes(n int) []byte {
	res := make([]byte, n)

	_, err := rand.Read(res)
	if err != nil {
		return nil
	}

	return res
}
