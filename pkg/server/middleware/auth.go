package middleware

import (
	"context"
	"grader/pkg/server/session"
	"net/http"
	"strings"
)

var (
	noAuthUrls = map[string]struct{}{
		"/login":  struct{}{},
		"/signup": struct{}{},
	}
	noSessUrls = map[string]struct{}{
		"/login":  struct{}{},
		"/signup": struct{}{},
	}
	noURLPrefixes = []string{
		"/api/",
	}
)

func CheckURLPrefix(url string, prefixes []string) bool {
	for _, prefix := range prefixes {
		if strings.HasPrefix(url, prefix) {
			return true
		}
	}

	return false
}

func Auth(sm session.Manager, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Cookie")

		path := r.URL.Path
		_, withoutAuth := noAuthUrls[path]
		isWhiteUrl := CheckURLPrefix(path, noURLPrefixes)

		if authHeader == "" && !withoutAuth && !isWhiteUrl {
			http.Error(w, "Authorization header not found", http.StatusUnauthorized)
			return
		}

		token := strings.TrimPrefix(authHeader, "token=")

		sess, err := sm.Check(token)

		_, withoutSess := noSessUrls[r.URL.Path]

		if err != nil && !withoutSess && !isWhiteUrl {
			http.Error(w, "Authorization error", http.StatusUnauthorized)

			return
		}

		ctx := context.WithValue(r.Context(), session.Key, sess)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
