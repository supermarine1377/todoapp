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
	infoHandler slog.Handler
	errorHander slog.Handler
}

func (h *Handler) Handle(ctx context.Context, r slog.Record) error {
	if requestID := request_id.Get(ctx); requestID != "" {
		r.AddAttrs(slog.Attr{Key: requestIDKey, Value: slog.StringValue(requestID)})
	}
	if r.Level >= slog.LevelWarn {
		return h.errorHander.Handle(ctx, r)
	}
	return h.infoHandler.Handle(ctx, r)
}

func (h *Handler) Enabled(ctx context.Context, level slog.Level) bool {
	if level >= slog.LevelWarn {
		return h.errorHander.Enabled(ctx, level)
	}
	return h.infoHandler.Enabled(ctx, level)
}

func (h *Handler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &Handler{
		infoHandler: h.infoHandler.WithAttrs(attrs),
		errorHander: h.errorHander.WithAttrs(attrs),
	}
}

// WithGroup creates a new handler with the given group name for both info and error handlers.
func (h *Handler) WithGroup(name string) slog.Handler {
	return &Handler{
		infoHandler: h.infoHandler.WithGroup(name),
		errorHander: h.errorHander.WithGroup(name),
	}
}

var handlerOption = slog.HandlerOptions{
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
}

func NewHandler() *Handler {
	var (
		info = slog.NewJSONHandler(os.Stdout, &handlerOption)
		err  = slog.NewJSONHandler(os.Stderr, &handlerOption)
	)
	return &Handler{infoHandler: info, errorHander: err}
}
