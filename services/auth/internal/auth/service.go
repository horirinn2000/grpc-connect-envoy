package auth

import (
	"context"
	"crypto/rsa"
	"time"
)

type AuthRepo interface {
	Authentication(ctx context.Context, username string, password string) (string, error)
}

type AuthService struct {
	privateKey *rsa.PrivateKey
	username   string
	password   string
	tokenExp   time.Duration
}

func NewService(p *rsa.PrivateKey, username, password string, tokenExp time.Duration) *AuthService {
	return &AuthService{
		privateKey: p,
		username:   username,
		password:   password,
		tokenExp:   tokenExp,
	}
}
