package auth

import (
	"context"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func (s *AuthService) Authentication(ctx context.Context, username string, password string) (string, error) {
	if username != s.username || password != s.password {
		return "", ErrInvalidInput
	}

	// Create a new token object, specifying signing method and the claims
	// you would like it to contain.
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, jwt.MapClaims{
		"sub": username,
		"iss": "auth-service",
		"aud": []string{"greet-service"},
		"exp": time.Now().Add(s.tokenExp).Unix(),
		"nbf": time.Now().Unix(),
		"iat": time.Now().Unix(),
	})

	// Sign and get the complete encoded token as a string using the private key
	tokenString, err := token.SignedString(s.privateKey)
	if err != nil {
		return "", ErrInternal
	}

	return tokenString, nil
}
