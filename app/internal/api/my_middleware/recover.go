package my_middleware

import (
	"log/slog"
	"net/http"
	"runtime/debug"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func Recover(logger *slog.Logger) echo.MiddlewareFunc {
	conf := middleware.DefaultRecoverConfig
	conf.LogErrorFunc = func(c echo.Context, err error, _ []byte) error {
		ctx := c.Request().Context()
		stackTrace := debug.Stack()
		logger.LogAttrs(ctx, slog.LevelError, "PANIC ERROR",
			slog.String("error", err.Error()),
			slog.Any("stack_trace", formatStackTrace(stackTrace)),
		)
		if !c.Response().Committed {
			return c.JSON(http.StatusInternalServerError, map[string]string{
				"error": "Internal Server Error",
			})
		}
		return nil
	}

	return middleware.RecoverWithConfig(conf)
}

// Helper function to format stack trace
func formatStackTrace(stackTrace []byte) []map[string]string {
	stackLines := strings.Split(string(stackTrace), "\n")
	var stackFrames []map[string]string

	// Process stack trace in pairs of function and location
	for i := 0; i < len(stackLines)-1; i += 2 {
		function := strings.TrimSpace(stackLines[i])
		location := strings.TrimSpace(stackLines[i+1])
		stackFrames = append(stackFrames, map[string]string{
			"function": function,
			"location": location,
		})
	}

	return stackFrames
}
