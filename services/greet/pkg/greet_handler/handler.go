package greethandler

import (
	"net/http"

	"github.com/horirinn2000/grpc-connect-envoy/services/greet/api/gen/greet/v1/greetv1connect"
	"github.com/horirinn2000/grpc-connect-envoy/services/greet/internal/greet"
)

// Handler は Connect Service Handlerの実装
type Handler struct {
	service greet.GreetRepo // 依存関係
}

// NewHandler は Handlerを初期化し、Connect Handlerを返します
func NewHandler(svc greet.GreetRepo) (string, http.Handler) {
	h := &Handler{
		service: svc,
	}
	// Connectの生成ロジックをここに置く
	return greetv1connect.NewGreetServiceHandler(h)
}
