package middleware

import (
	"context"
	"net/http"
	"os"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

type JWTUser struct {
	ID       int
	Username string
	Email    string
	Role     string
}

type contextKey string

const UserCtxKey = contextKey("user")

func JWTAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Требуется авторизация", http.StatusUnauthorized)
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			http.Error(w, "Некорректный заголовок авторизации", http.StatusUnauthorized)
			return
		}

		tokenStr := parts[1]
		claims := jwt.MapClaims{}
		jwtKey := []byte(os.Getenv("JWT_SECRET")) // <-- тут!

		token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
			return jwtKey, nil
		})
		if err != nil || !token.Valid {
			http.Error(w, "Недействительный токен", http.StatusUnauthorized)
			return
		}

		// Достаём все нужные данные из claims
		userID, okID := claims["user_id"].(float64)
		username, okU := claims["username"].(string)
		email, okE := claims["email"].(string)
		role, okR := claims["role"].(string)
		if !okID || !okU || !okE || !okR {
			http.Error(w, "Ошибка данных токена", http.StatusUnauthorized)
			return
		}

		user := &JWTUser{
			ID:       int(userID),
			Username: username,
			Email:    email,
			Role:     role,
		}

		ctx := context.WithValue(r.Context(), UserCtxKey, user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
