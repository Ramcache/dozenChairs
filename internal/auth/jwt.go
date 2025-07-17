package auth

import (
    "time"
    "github.com/golang-jwt/jwt/v5"
)

var jwtKey = []byte("my_secret_key")

func GenerateToken(username string, ttl time.Duration) (string, error) {
    claims := jwt.MapClaims{
        "username": username,
        "exp": time.Now().Add(ttl).Unix(),
    }

    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    return token.SignedString(jwtKey)
}

func ValidateToken(tokenString string) (string, error) {
    token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
        return jwtKey, nil
    })

    if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
        return claims["username"].(string), nil
    } else {
        return "", err
    }
}
