package greethandler

import (
	"context"
	"errors"

	"connectrpc.com/connect"
	greetv1 "github.com/horirinn2000/grpc-connect-envoy/services/greet/api/gen/greet/v1"
)

// Thanks は protoc-gen-connect-goによって要求されるメソッド
func (h *Handler) Greet(
	ctx context.Context,
	req *greetv1.GreetRequest,
) (*greetv1.GreetResponse, error) {
	// 処理は内部サービスに委譲
	if req.Name == "" {
		return nil, connect.NewError(connect.CodeInvalidArgument, errors.New("name is required"))
	}
	message, err := h.service.GenerateGreet(ctx, req.Name)
	if err != nil {
		return nil, h.mapErrorToConnectCode(err)
	}
	res := &greetv1.GreetResponse{Greeting: message}
	return res, nil
}
