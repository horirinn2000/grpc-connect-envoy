package auth

import (
	"context"
	"crypto/rsa"
	"time"
)

type AuthRepo interface {
	Authentication(ctx context.Context, username string, password string) (string, string, error)
	RefreshToken(ctx context.Context, refreshToken string) (string, string, error)
}

type AuthService struct {
	privateKey      *rsa.PrivateKey
	username        string
	password        string
	tokenExp        time.Duration
	refreshTokenExp time.Duration
	refreshToken    map[string]string
}

func NewService(p *rsa.PrivateKey, username, password string, tokenExp, refreshTokenExp time.Duration) *AuthService {
	return &AuthService{
		privateKey:      p,
		username:        username,
		password:        password,
		tokenExp:        tokenExp,
		refreshTokenExp: refreshTokenExp,
		refreshToken:    make(map[string]string),
	}
}
