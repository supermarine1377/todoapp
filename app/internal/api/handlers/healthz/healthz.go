// package healthz は死活監視用APIのハンドラーを提供します。
package healthz

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

// Healthz はサーバーの死活監視用のAPIのハンドラー
//
//	@Summary サーバーの死活監視用のAPI
//	@Description	サーバーの死活監視用のAPI
//	@Tags			accounts
//	@Accept			json
//	@Produce		json
//	@Success		200
//	@Router			/healthz [get]
func Healthz(c echo.Context) error {
	return c.JSON(http.StatusOK, nil)
}
