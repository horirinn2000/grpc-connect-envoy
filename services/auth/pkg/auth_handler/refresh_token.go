package authhandler

import (
	"context"
	"errors"

	"connectrpc.com/connect"
	authv1 "github.com/horirinn2000/grpc-connect-envoy/services/auth/api/gen/auth/v1"
)

func (h *Handler) RefreshToken(
	ctx context.Context,
	req *authv1.RefreshTokenRequest,
) (*authv1.RefreshTokenResponse, error) {
	if req.RefreshToken == "" {
		return nil, connect.NewError(connect.CodeInvalidArgument, errors.New("refresh token is required"))
	}

	tokenString, refreshTokenString, err := h.service.RefreshToken(ctx, req.RefreshToken)
	if err != nil {
		return nil, h.mapErrorToConnectCode(err)
	}

	return &authv1.RefreshTokenResponse{
		Token:        tokenString,
		RefreshToken: refreshTokenString,
	}, nil
}
