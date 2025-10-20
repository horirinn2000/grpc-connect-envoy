package thanks

import (
	"context"
)

type ThanksRepo interface {
	GenerateThanks(ctx context.Context, name string) (string, error)
}

// Service は CoreService の具体的な実装
type ThanksService struct{}

func NewService() *ThanksService {
	return &ThanksService{}
}
