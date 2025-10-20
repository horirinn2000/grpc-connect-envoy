package greethandler

import (
	"context"
	"errors"
	"testing"

	"connectrpc.com/connect"
	"github.com/google/go-cmp/cmp"
	greetv1 "github.com/horirinn2000/grpc-connect-envoy/services/greet/api/gen/greet/v1"
	"github.com/horirinn2000/grpc-connect-envoy/services/greet/internal/greet"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/testing/protocmp"
)

// mockGreetService は greet.GreetRepo インターフェースを満たす偽物のサービス
type mockGreetService struct {
	// このモックがどのような値を返すかを設定するためのフィールド
	GenerateGreetFunc func(ctx context.Context, name string) (string, error)
}

// GreetRepoインターフェースのメソッドを実装
func (m *mockGreetService) GenerateGreet(ctx context.Context, name string) (string, error) {
	// もし関数が設定されていれば、それを呼び出す
	if m.GenerateGreetFunc != nil {
		return m.GenerateGreetFunc(ctx, name)
	}
	// デフォルトの動作（何も設定されていない場合）
	return "", errors.New("GenerateGreetFunc is not set")
}

func Test_Greet(t *testing.T) {
	// テストケースを構造体のスライスとして定義
	testCases := []struct {
		name              string                                                 // テストケース名
		input             *greetv1.GreetRequest                                  // 入力値
		mockGenerateGreet func(ctx context.Context, name string) (string, error) // モックの振る舞いを定義
		expected          *greetv1.GreetResponse                                 // 期待されるメッセージ
		expectedErr       error                                                  // 期待されるエラー
	}{
		{
			name:  "正常系: 名前を渡した場合",
			input: &greetv1.GreetRequest{Name: "hori"},
			mockGenerateGreet: func(ctx context.Context, name string) (string, error) {
				require.Equal(t, "hori", name) // サービスが正しい引数で呼ばれたか
				return "Greet, hori! Thank you for using our service.", nil
			},
			expected: &greetv1.GreetResponse{
				Greeting: "Greet, hori! Thank you for using our service.",
			},
			expectedErr: nil,
		},
		{
			name:        "異常系: ハンドラでの入力値チェック",
			input:       &greetv1.GreetRequest{Name: ""},
			expected:    nil,
			expectedErr: connect.NewError(connect.CodeInvalidArgument, errors.New("name is required")),
		},
		{
			name:  "異常系: サービスがエラーを返す",
			input: &greetv1.GreetRequest{Name: "error-case"},
			mockGenerateGreet: func(ctx context.Context, name string) (string, error) {
				// 例えば internal/greet.ErrNotFound を返すケース
				return "", greet.ErrNotFound // internal/greet からエラーをインポートする必要がある
			},
			expected: nil,
			// mapErrorToConnectCode で変換された後のエラーを期待値とする
			expectedErr: connect.NewError(connect.CodeNotFound, greet.ErrNotFound),
		},
	}

	// 各テストケースをループで実行
	for _, tc := range testCases {
		// t.Run を使うと、テストが失敗したときにどのケースで失敗したかが分かりやすくなります
		t.Run(tc.name, func(t *testing.T) {
			// Arrange (準備)
			s := &mockGreetService{GenerateGreetFunc: tc.mockGenerateGreet}
			h := Handler{service: s}

			// Act (実行)
			res, err := h.Greet(context.Background(), tc.input)

			// Assert (検証)
			if diff := cmp.Diff(tc.expected, res, protocmp.Transform()); diff != "" {
				t.Errorf("Greet() mismatch (-want +got):\n%s", diff)
			}
			// require.ErrorIs を使うと、ラップされたエラーも正しく比較できます
			if tc.expectedErr != nil {
				// 1. エラーが存在することを確認
				require.Error(t, err)

				// 2. エラーを Connect エラーに変換（キャスト）
				connErr := new(connect.Error)

				// 3. 変換可能か、つまり Connect エラーであることを確認
				// errors.As を使って *connect.Error として取り出す
				if errors.As(err, &connErr) {
					// 4. Connect コードが期待通りかを確認
					expectedConnectErr := tc.expectedErr.(*connect.Error)
					require.Equal(t, expectedConnectErr.Code(), connErr.Code(), "Connect Code mismatch")

					// 5. エラーメッセージのテキスト部分が期待通りかを確認
					// Connectエラーは通常、コードに続いてメッセージを持つため、Message()で確認する
					require.Equal(t, expectedConnectErr.Message(), connErr.Message(), "Error message mismatch")
				} else {
					// Connectエラーではない場合のエラー
					t.Fatalf("Returned error is not a *connect.Error: %v", err)
				}
			} else {
				require.NoError(t, err)
			}
		})
	}
}
