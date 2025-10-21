package auth

import (
	"context"
	"crypto/rsa"
)

type AuthRepo interface {
	Authentication(ctx context.Context, username string, password string) (string, error)
}

type AuthService struct {
	privateKey *rsa.PrivateKey
}

func NewService(p *rsa.PrivateKey) *AuthService {
	return &AuthService{privateKey: p}
}
