package thanks

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// GenerateThanks: Protobufのメッセージ型ではなく、Goのプリミティブ型を扱う
func Test_GenerateThanks(t *testing.T) {
	// テストケースを構造体のスライスとして定義
	testCases := []struct {
		name        string // テストケース名
		inputName   string // 入力値
		expectedMsg string // 期待されるメッセージ
		expectedErr error  // 期待されるエラー
	}{
		{
			name:        "正常系: 名前を渡した場合",
			inputName:   "hori",
			expectedMsg: "Thanks, hori! Thank you for using our service.",
			expectedErr: nil,
		},
		{
			name:        "異常系: 空文字列を渡した場合",
			inputName:   "",
			expectedMsg: "",
			expectedErr: ErrInvalidInput,
		},
	}

	// 各テストケースをループで実行
	for _, tc := range testCases {
		// t.Run を使うと、テストが失敗したときにどのケースで失敗したかが分かりやすくなります
		t.Run(tc.name, func(t *testing.T) {
			// Arrange (準備)
			s := NewService()
			ctx := context.Background()

			// Act (実行)
			msg, err := s.GenerateThanks(ctx, tc.inputName)

			// Assert (検証)
			assert.Equal(t, tc.expectedMsg, msg)
			// require.ErrorIs を使うと、ラップされたエラーも正しく比較できます
			if tc.expectedErr != nil {
				require.ErrorIs(t, err, tc.expectedErr)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
