package thankshandler

import (
	"net/http"

	"github.com/horirinn2000/grpc-connect-envoy/services/thanks/api/gen/thanks/v1/thanksv1connect"
	"github.com/horirinn2000/grpc-connect-envoy/services/thanks/internal/thanks"
)

// Handler は Connect Service Handlerの実装
type Handler struct {
	service thanks.ThanksRepo // 依存関係
}

// NewHandler は Handlerを初期化し、Connect Handlerを返します
func NewHandler(svc thanks.ThanksRepo) (string, http.Handler) {
	h := &Handler{
		service: svc,
	}
	// Connectの生成ロジックをここに置く
	return thanksv1connect.NewThanksServiceHandler(h)
}
