package task_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/supermarine1377/todoapp/app/common/apperrors"
	"github.com/supermarine1377/todoapp/app/internal/api/handlers/task"
	"github.com/supermarine1377/todoapp/app/internal/api/handlers/task/mock"
	"go.uber.org/mock/gomock"
)

func TestTaskHandler_Create(t *testing.T) {
	tests := []struct {
		name        string
		prepareMock func(ctx context.Context) *mock.MockTaskRepository
		wantErr     bool
		statusCode  int
	}{
		{
			name: "Bad request status",
			prepareMock: func(ctx context.Context) *mock.MockTaskRepository {
				ctrl := gomock.NewController(t)
				m := mock.NewMockTaskRepository(ctrl)
				m.EXPECT().CreateCtx(ctx, gomock.Any()).Return(apperrors.ErrBadRequest)
				return m
			},
			wantErr:    false,
			statusCode: http.StatusBadRequest,
		},
		{
			name: "Internal server error status",
			prepareMock: func(ctx context.Context) *mock.MockTaskRepository {
				ctrl := gomock.NewController(t)
				m := mock.NewMockTaskRepository(ctrl)
				m.EXPECT().CreateCtx(ctx, gomock.Any()).Return(apperrors.ErrInternalServerError)
				return m
			},
			wantErr:    false,
			statusCode: http.StatusInternalServerError,
		},
		{
			name: "Created status",
			prepareMock: func(ctx context.Context) *mock.MockTaskRepository {
				ctrl := gomock.NewController(t)
				m := mock.NewMockTaskRepository(ctrl)
				m.EXPECT().CreateCtx(ctx, gomock.Any()).Return(nil)
				return m
			},
			wantErr:    false,
			statusCode: http.StatusCreated,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := tt.prepareMock(context.Background())
			th := task.NewTaskHandler(m)

			r := httptest.NewRequest(http.MethodPost, "/", nil)
			rc := httptest.NewRecorder()
			c := echo.New().NewContext(r, rc)

			if err := th.Create(c); (err != nil) != tt.wantErr {
				t.Errorf("TaskHandler.Create() error = %v, wantErr %v", err, tt.wantErr)
			}

			assert.Equal(t, tt.statusCode, rc.Code)
		})
	}
}
