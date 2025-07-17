package middleware

import (
	"net/http"
)

func OnlyAdmin(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user, ok := r.Context().Value(UserCtxKey).(*JWTUser)
		if !ok || user.Role != "admin" {
			http.Error(w, "Доступ только для администратора", http.StatusForbidden)
			return
		}
		next.ServeHTTP(w, r)
	})
}
