package my_middleware

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/supermarine1377/todoapp/app/common/request_id"
)

func RequestID() echo.MiddlewareFunc {
	return middleware.RequestIDWithConfig(middleware.RequestIDConfig{
		RequestIDHandler: func(c echo.Context, requestID string) {
			req := c.Request()
			newCtx := request_id.Set(req.Context(), requestID)
			newReq := req.WithContext(newCtx)
			c.SetRequest(newReq)
		},
	})
}
