package thanks

import (
	"context"
	"fmt"
)

// GenerateThanks: Protobufのメッセージ型ではなく、Goのプリミティブ型を扱う
func (s *ThanksService) GenerateThanks(ctx context.Context, name string) (string, error) {
	if name == "" {
		return "", ErrInvalidInput
	}
	return fmt.Sprintf("Thanks, %s! Thank you for using our service.", name), nil
}
