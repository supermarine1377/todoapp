// package todo はTODO関係のAPIハンドラーを提供する
package todo

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/supermarine1377/todoapp/app/common/apperrors"
	"github.com/supermarine1377/todoapp/app/internal/api/handlers/todo/mock"
	"go.uber.org/mock/gomock"
)

func TestTODOHandler_Create(t *testing.T) {
	tests := []struct {
		name        string
		prepareMock func() *mock.MockTODORepository
		wantErr     bool
		statusCode  int
	}{
		{
			name: "Bad request status",
			prepareMock: func() *mock.MockTODORepository {
				ctrl := gomock.NewController(t)
				m := mock.NewMockTODORepository(ctrl)
				m.EXPECT().Create().Return(apperrors.ErrBadRequest)
				return m
			},
			wantErr:    false,
			statusCode: http.StatusBadRequest,
		},
		{
			name: "Internal server error status",
			prepareMock: func() *mock.MockTODORepository {
				ctrl := gomock.NewController(t)
				m := mock.NewMockTODORepository(ctrl)
				m.EXPECT().Create().Return(apperrors.ErrInternalServerError)
				return m
			},
			wantErr:    false,
			statusCode: http.StatusInternalServerError,
		},
		{
			name: "Created status",
			prepareMock: func() *mock.MockTODORepository {
				ctrl := gomock.NewController(t)
				m := mock.NewMockTODORepository(ctrl)
				m.EXPECT().Create().Return(nil)
				return m
			},
			wantErr:    false,
			statusCode: http.StatusCreated,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := tt.prepareMock()
			th := NewTODOHandler(m)

			r := httptest.NewRequest(http.MethodPost, "/", nil)
			rc := httptest.NewRecorder()
			c := echo.New().NewContext(r, rc)

			if err := th.Create(c); (err != nil) != tt.wantErr {
				t.Errorf("TODOHandler.Create() error = %v, wantErr %v", err, tt.wantErr)
			}

			assert.Equal(t, tt.statusCode, rc.Code)
		})
	}
}
