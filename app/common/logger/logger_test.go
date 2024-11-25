package logger

import (
	"bytes"
	"context"
	"log/slog"
	"testing"

	"github.com/supermarine1377/todoapp/app/common/request_id"
)

const requestID = "request_id"

func TestHandler_Handle(t *testing.T) {
	var infoBuffer, errorBuffer bytes.Buffer

	infoHandler := slog.NewJSONHandler(&infoBuffer, nil)
	errorHandler := slog.NewJSONHandler(&errorBuffer, nil)

	h := &Handler{
		infoHandler: infoHandler,
		errorHander: errorHandler,
	}

	ctx := request_id.Set(context.Background(), requestID)

	record := slog.Record{
		Level:   slog.LevelInfo,
		Message: "info log",
	}
	if err := h.Handle(ctx, record); err != nil {
		t.Fatalf("Handle() error = %v", err)
	}

	if !bytes.Contains(infoBuffer.Bytes(), []byte("info log")) {
		t.Errorf("Expected info log in infoHandler, got: %s", infoBuffer.String())
	}
	if bytes.Contains(errorBuffer.Bytes(), []byte("info log")) {
		t.Errorf("Info log was written to errorHandler")
	}

	infoBuffer.Reset()
	errorBuffer.Reset()

	record = slog.Record{
		Level:   slog.LevelError,
		Message: "error log",
	}
	if err := h.Handle(ctx, record); err != nil {
		t.Fatalf("Handle() error = %v", err)
	}

	if !bytes.Contains(errorBuffer.Bytes(), []byte("error log")) {
		t.Errorf("Expected error log in errorHandler, got: %s", errorBuffer.String())
	}
	if bytes.Contains(infoBuffer.Bytes(), []byte("error log")) {
		t.Errorf("Error log was written to infoHandler")
	}
}
