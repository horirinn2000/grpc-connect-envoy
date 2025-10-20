package thankshandler

import (
	"errors"

	"connectrpc.com/connect"
	"github.com/horirinn2000/grpc-connect-envoy/services/thanks/internal/thanks"
)

// mapErrorToConnectCode は内部エラーをConnectエラーに変換します。
func (h *Handler) mapErrorToConnectCode(err error) *connect.Error {
	// errors.Is を使用して、内部で定義した特定のエラーと比較します
	if errors.Is(err, thanks.ErrNotFound) {
		// リソースが見つからない場合 -> Connect Code NotFound (404)
		return connect.NewError(connect.CodeNotFound, err)
	}

	if errors.Is(err, thanks.ErrInvalidInput) {
		// クライアントからの入力データが不正な場合 -> Connect Code InvalidArgument (400)
		return connect.NewError(connect.CodeInvalidArgument, err)
	}

	// 上記のいずれにも該当しない、想定外または一般的な内部エラーの場合
	// -> Connect Code Internal または CodeUnknown を返すのが標準的です
	return connect.NewError(connect.CodeInternal, errors.New("internal server error"))
}
