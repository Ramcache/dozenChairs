package middlewares

import (
	"context"
	"dozenChairs/internal/auth"
	"dozenChairs/pkg/config"
	"dozenChairs/pkg/httphelper"
	"net/http"
	"strings"
)

type contextKey string

const (
	userIDKey contextKey = "userID"
	roleKey   contextKey = "role"
)

func RequireAuth(jwt *auth.JWTManager) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			cfg := config.LoadConfig()
			if !cfg.AuthEnabled {
				// Пропускаем авторизацию
				ctx := context.WithValue(r.Context(), userIDKey, "debug-user")
				ctx = context.WithValue(ctx, roleKey, "admin") // или "user" — смотри сам
				next.ServeHTTP(w, r.WithContext(ctx))
				return
			}

			authHeader := r.Header.Get("Authorization")
			if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
				httphelper.WriteError(w, http.StatusUnauthorized, "Missing or invalid Authorization header")
				return
			}

			token := strings.TrimPrefix(authHeader, "Bearer ")
			userID, role, err := jwt.ValidateAccess(token)
			if err != nil {
				httphelper.WriteError(w, http.StatusUnauthorized, "Invalid or expired token")
				return
			}

			ctx := context.WithValue(r.Context(), userIDKey, userID)
			ctx = context.WithValue(ctx, roleKey, role)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func RequireRole(requiredRole string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			cfg := config.LoadConfig()
			if !cfg.AuthEnabled {
				next.ServeHTTP(w, r) // пропускаем проверку роли
				return
			}

			role, ok := r.Context().Value(roleKey).(string)
			if !ok || role != requiredRole {
				httphelper.WriteError(w, http.StatusForbidden, "Forbidden: insufficient permissions")
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

func Role() contextKey {
	return roleKey
}

func UserID() contextKey {
	return userIDKey
}
