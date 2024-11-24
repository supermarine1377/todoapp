package request_id

import (
	"context"
)

type key struct{}

func Set(ctx context.Context, requestID string) context.Context {
	return context.WithValue(ctx, key{}, requestID)
}

func Get(ctx context.Context) string {
	if requestID, ok := ctx.Value(key{}).(string); ok {
		return requestID
	}
	return ""
}
