package middleware

import (
	"context"
	"net/http"
	"strings"
)

func Middleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			tokenString := strings.TrimSpace(strings.TrimPrefix(strings.TrimSpace(r.Header.Get("Authorization")), "Bearer"))
			if tokenString == "" {
				next.ServeHTTP(w, r)
				// http.Error(w, "unauthorized", http.StatusUnauthorized)
				return
			}

			r = r.WithContext(context.WithValue(r.Context(), "has_token", true))
			r = r.WithContext(context.WithValue(r.Context(), "token", tokenString))

			next.ServeHTTP(w, r)
		})
	}
}
