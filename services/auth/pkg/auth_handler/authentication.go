package authhandler

import (
	"context"
	"errors"

	"connectrpc.com/connect"
	authv1 "github.com/horirinn2000/grpc-connect-envoy/services/auth/api/gen/auth/v1"
)

func (h *Handler) Authenticate(
	ctx context.Context,
	req *authv1.AuthenticateRequest,
) (*authv1.AuthenticateResponse, error) {
	if req.Username == "" {
		return nil, connect.NewError(connect.CodeInvalidArgument, errors.New("username is required"))
	}
	if req.Password == "" {
		return nil, connect.NewError(connect.CodeInvalidArgument, errors.New("password is required"))
	}

	tokenString, refreshTokenString, err := h.service.Authentication(ctx, req.Username, req.Password)
	if err != nil {
		return nil, h.mapErrorToConnectCode(err)
	}
	return &authv1.AuthenticateResponse{
		Token:        tokenString,
		RefreshToken: refreshTokenString,
	}, nil
}
