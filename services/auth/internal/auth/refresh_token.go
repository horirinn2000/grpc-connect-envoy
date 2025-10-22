package auth

import (
	"context"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func (s *AuthService) RefreshToken(ctx context.Context, refreshToken string) (string, string, error) {
	username, ok := s.refreshToken[refreshToken]
	if !ok {
		return "", "", ErrInvalidInput
	}

	delete(s.refreshToken, refreshToken)

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, jwt.MapClaims{
		"sub": username,
		"iss": "auth-service",
		"aud": []string{"greet-service"},
		"exp": time.Now().Add(s.tokenExp).Unix(),
		"nbf": time.Now().Unix(),
		"iat": time.Now().Unix(),
	})

	tokenString, err := token.SignedString(s.privateKey)
	if err != nil {
		return "", "", ErrInternal
	}

	newRefreshToken := jwt.NewWithClaims(jwt.SigningMethodRS256, jwt.MapClaims{
		"sub": username,
		"iss": "auth-service",
		"aud": []string{"greet-service"},
		"exp": time.Now().Add(s.refreshTokenExp).Unix(),
		"nbf": time.Now().Unix(),
		"iat": time.Now().Unix(),
	})

	newRefreshTokenString, err := newRefreshToken.SignedString(s.privateKey)
	if err != nil {
		return "", "", ErrInternal
	}

	s.refreshToken[newRefreshTokenString] = username

	return tokenString, newRefreshTokenString, nil
}
