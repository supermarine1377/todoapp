// package healthz は死活監視用APIのハンドラーを提供します。
package healthz

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

// Healthz はサーバーの死活監視用のAPIのハンドラー
func Healthz(c echo.Context) error {
	return c.JSON(http.StatusOK, nil)
}
