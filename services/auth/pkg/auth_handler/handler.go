package authhandler

import (
	"net/http"

	"github.com/horirinn2000/grpc-connect-envoy/services/auth/api/gen/auth/v1/authv1connect"
	"github.com/horirinn2000/grpc-connect-envoy/services/auth/internal/auth"
)

// Handler は Connect Service Handlerの実装
type Handler struct {
	service auth.AuthRepo // 依存関係
}

// NewHandler は Handlerを初期化し、Connect Handlerを返します
func NewHandler(svc auth.AuthRepo) (string, http.Handler) {
	h := &Handler{
		service: svc,
	}
	// Connectの生成ロジックをここに置く
	return authv1connect.NewAuthServiceHandler(h)
}
