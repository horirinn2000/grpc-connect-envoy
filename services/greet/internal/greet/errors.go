package greet

import "errors"

// NewErrorは、リソースが見つからない場合に内部サービスが返すエラー
var ErrNotFound = errors.New("greet resource not found")

// NewErrorは、入力値が不正で処理できない場合に返すエラー
var ErrInvalidInput = errors.New("invalid input data")
