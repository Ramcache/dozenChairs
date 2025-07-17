package utils

import (
	"dozenChairs/internal/models"
	"dozenChairs/pkg/config"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func GenerateJWT(user *models.User) (string, error) {
	claims := jwt.MapClaims{
		"user_id":  user.ID,
		"username": user.Username,
		"email":    user.Email,
		"role":     user.Role,
		"exp":      time.Now().Add(24 * time.Hour).Unix(),
	}
	jwtKey := []byte(config.Load().JWTSecret)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtKey)
}
