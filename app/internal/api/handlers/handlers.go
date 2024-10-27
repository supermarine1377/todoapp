// package handlers はAPIのハンドラーを提供します。
package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

// Healthz はサーバーの死活監視用のAPIのハンドラー
func Healthz(c echo.Context) error {
	return c.JSON(http.StatusOK, nil)
}
