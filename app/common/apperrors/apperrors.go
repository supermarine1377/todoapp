// package apperror は各パッケージで共通で用いるエラーを定義する
package apperrors

import "errors"

var (
	// ErrBadRequest は不正なリクエストに対するエラーを表す
	ErrBadRequest = errors.New("Bad request")
	// ErrNotFound はリクエストに対するデータが見つからなかったことを表す
	ErrNotFound = errors.New("Not found")
	// ErrInternalServerError はサーバー内部でエラーが起きたことを表す
	ErrInternalServerError = errors.New("Internal server error")
)
