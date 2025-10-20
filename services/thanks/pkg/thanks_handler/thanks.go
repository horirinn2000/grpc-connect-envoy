package thankshandler

import (
	"context"
	"errors"

	"connectrpc.com/connect"
	thanksv1 "github.com/horirinn2000/grpc-connect-envoy/services/thanks/api/gen/thanks/v1"
)

// Thanks は protoc-gen-connect-goによって要求されるメソッド
func (h *Handler) Thanks(
	ctx context.Context,
	req *thanksv1.ThanksRequest,
) (*thanksv1.ThanksResponse, error) {
	// 処理は内部サービスに委譲
	if req.Name == "" {
		return nil, connect.NewError(connect.CodeInvalidArgument, errors.New("name is required"))
	}
	message, err := h.service.GenerateThanks(ctx, req.Name)
	if err != nil {
		return nil, h.mapErrorToConnectCode(err)
	}
	res := &thanksv1.ThanksResponse{Thanks: message}
	return res, nil
}
