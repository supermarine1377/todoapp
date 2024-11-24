package request_id_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/supermarine1377/todoapp/app/common/request_id"
)

func TestSet_Get(t *testing.T) {
	var (
		requestID = "hoge"
		ctx       = context.Background()
	)
	ctx = request_id.Set(ctx, requestID)

	got := request_id.Get(ctx)
	assert.Equal(t, requestID, got)
}
