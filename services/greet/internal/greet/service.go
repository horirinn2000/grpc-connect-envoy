package greet

import (
	"context"
)

type GreetRepo interface {
	GenerateGreet(ctx context.Context, name string) (string, error)
}

type GreetService struct{}

func NewService() *GreetService {
	return &GreetService{}
}
