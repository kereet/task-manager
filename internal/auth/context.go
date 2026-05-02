package auth

import (
	"errors"
	"net/http"
	"strings"
)

func ExtractTokenFromRequest(r *http.Request) (string, error) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return "", errors.New("authorization header is required")
	}

	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		return "", errors.New("invalid authorization format. Use: Bearer <token>")
	}

	return parts[1], nil
}

func GetUserFromRequest(r *http.Request) (*Claims, error) {
	tokenString, err := ExtractTokenFromRequest(r)
	if err != nil {
		return nil, err
	}

	claims, err := ValidateToken(tokenString)
	if err != nil {
		return nil, err
	}

	return claims, nil
}
