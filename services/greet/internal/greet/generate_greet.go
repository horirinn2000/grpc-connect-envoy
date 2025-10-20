package greet

import (
	"context"
	"fmt"
)

// GenerateThanks: Protobufのメッセージ型ではなく、Goのプリミティブ型を扱う
func (s *GreetService) GenerateGreet(ctx context.Context, name string) (string, error) {
	if name == "" {
		return "", ErrInvalidInput
	}
	return fmt.Sprintf("Greet, %s! Thank you for using our service.", name), nil
}
