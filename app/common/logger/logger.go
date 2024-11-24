// package logger は、アプリケーション全体で使うロガーを提供する
package logger

import (
	"context"
	"log/slog"
	"os"

	"github.com/supermarine1377/todoapp/app/common/request_id"
)

const requestIDKey = "request_id"

type Handler struct {
	slog.Handler
}

func (h *Handler) Handle(ctx context.Context, r slog.Record) error {
	if requestID := request_id.Get(ctx); requestID != "" {
		r.AddAttrs(slog.Attr{Key: requestIDKey, Value: slog.StringValue(requestID)})
	}
	return h.Handler.Handle(ctx, r)
}

func NewHandler() *Handler {
	handler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		AddSource: true,
		Level:     slog.LevelInfo,
		// Cloud Logging対応
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			if a.Key == slog.MessageKey {
				a.Key = "message"
			}
			if a.Key == slog.LevelKey {
				a.Key = "severity"
			}
			return a
		},
	})
	return &Handler{Handler: handler}
}
