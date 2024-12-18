package task_test

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/supermarine1377/todoapp/app/common/apperrors"
	"github.com/supermarine1377/todoapp/app/internal/api/handlers/task"
	"github.com/supermarine1377/todoapp/app/internal/api/handlers/task/mock"
	entity_task "github.com/supermarine1377/todoapp/app/internal/model/entity/task"
	"go.uber.org/mock/gomock"
)

func TestTaskHandler_Create(t *testing.T) {
	tests := []struct {
		name        string
		req         func() *http.Request
		prepareMock func(ctx context.Context) *mock.MockTaskRepository
		wantErr     bool
		statusCode  int
		message     string
	}{
		{
			name: "Invalid json",
			req: func() *http.Request {
				b := strings.NewReader("{invalid}")
				req := httptest.NewRequest(http.MethodPost, "/", b)
				req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
				return req
			},
			prepareMock: func(ctx context.Context) *mock.MockTaskRepository {
				return nil
			},
			wantErr:    false,
			statusCode: http.StatusBadRequest,
		},
		{
			name: "Missing Content-Type",
			req: func() *http.Request {
				b := strings.NewReader(`{}`)
				return httptest.NewRequest(http.MethodPost, "/", b)
			},
			prepareMock: func(ctx context.Context) *mock.MockTaskRepository {
				return nil
			},
			wantErr:    false,
			statusCode: http.StatusUnsupportedMediaType,
		},
		{
			name: "Missing title",
			req: func() *http.Request {
				b := strings.NewReader(`{}`)
				req := httptest.NewRequest(http.MethodPost, "/", b)
				req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
				return req
			},
			prepareMock: func(ctx context.Context) *mock.MockTaskRepository {
				return nil
			},
			wantErr:    false,
			statusCode: http.StatusBadRequest,
		},
		{
			name: "Task successfully created",
			req: func() *http.Request {
				b := strings.NewReader(`{"title":"hoge"}`)
				req := httptest.NewRequest(http.MethodPost, "/", b)
				req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
				return req
			},
			prepareMock: func(ctx context.Context) *mock.MockTaskRepository {
				ctrl := gomock.NewController(t)
				m := mock.NewMockTaskRepository(ctrl)
				m.EXPECT().CreateCtx(ctx, gomock.Any()).Return(nil)
				return m
			},
			wantErr:    false,
			statusCode: http.StatusCreated,
		},
		{
			name: "Task successfully created",
			req: func() *http.Request {
				b := strings.NewReader(`{"title":"hoge"}`)
				req := httptest.NewRequest(http.MethodPost, "/", b)
				req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
				return req
			},
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
			name: "Task successfully created",
			req: func() *http.Request {
				b := strings.NewReader(`{"title":"hoge"}`)
				req := httptest.NewRequest(http.MethodPost, "/", b)
				req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
				return req
			},
			prepareMock: func(ctx context.Context) *mock.MockTaskRepository {
				ctrl := gomock.NewController(t)
				m := mock.NewMockTaskRepository(ctrl)
				m.EXPECT().CreateCtx(ctx, gomock.Any()).Return(apperrors.ErrInternalServerError)
				return m
			},
			wantErr:    false,
			statusCode: http.StatusInternalServerError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := echo.New()
			m := tt.prepareMock(context.Background())
			th := task.NewTaskHandler(m)
			rc := httptest.NewRecorder()
			c := e.NewContext(tt.req(), rc)

			if err := th.Create(c); (err != nil) != tt.wantErr {
				t.Errorf("TaskHandler.Create() error = %v, wantErr %v", err, tt.wantErr)
			}

			assert.Equal(t, tt.statusCode, rc.Code)
		})
	}
}

var tasksInDB = entity_task.Tasks{{
	ID:        1,
	Title:     "dummy",
	CreatedAt: 1,
	UpdatedAt: 1,
}}

func TestTaskHandler_List(t *testing.T) {
	tests := []struct {
		name        string
		req         func() *http.Request
		prepareMock func(ctx context.Context) *mock.MockTaskRepository
		wantErr     bool
		statusCode  int
		checkRes    bool
	}{
		{
			name: "Fetch tasks with default limit",
			req: func() *http.Request {
				req := httptest.NewRequest(http.MethodGet, "/tasks", nil)
				req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
				return req
			},
			prepareMock: func(ctx context.Context) *mock.MockTaskRepository {
				ctrl := gomock.NewController(t)
				m := mock.NewMockTaskRepository(ctrl)
				m.EXPECT().ListCtx(gomock.Any(), 0, 10).Return(&tasksInDB, nil)
				return m
			},
			wantErr:    false,
			checkRes:   true,
			statusCode: http.StatusOK,
		},
		{
			name: "Fetch tasks with offset of 5 and default limit",
			req: func() *http.Request {
				req := httptest.NewRequest(http.MethodGet, "/tasks?offset=5", nil)
				req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
				return req
			},
			prepareMock: func(ctx context.Context) *mock.MockTaskRepository {
				ctrl := gomock.NewController(t)
				m := mock.NewMockTaskRepository(ctrl)
				m.EXPECT().ListCtx(gomock.Any(), 5, 10)
				return m
			},
			wantErr:    false,
			statusCode: http.StatusOK,
		},
		{
			name: "Fetch tasks with limit=5",
			req: func() *http.Request {
				req := httptest.NewRequest(http.MethodGet, "/tasks?limit=5", nil)
				req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
				return req
			},
			prepareMock: func(ctx context.Context) *mock.MockTaskRepository {
				ctrl := gomock.NewController(t)
				m := mock.NewMockTaskRepository(ctrl)
				m.EXPECT().ListCtx(gomock.Any(), 0, 5)
				return m
			},
			wantErr:    false,
			statusCode: http.StatusOK,
		},
		{
			name: "Fetch tasks with limit=5 and offset=5",
			req: func() *http.Request {
				req := httptest.NewRequest(http.MethodGet, "/tasks?limit=5&offset=5", nil)
				req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
				return req
			},
			prepareMock: func(ctx context.Context) *mock.MockTaskRepository {
				ctrl := gomock.NewController(t)
				m := mock.NewMockTaskRepository(ctrl)
				m.EXPECT().ListCtx(gomock.Any(), 5, 5)
				return m
			},
			wantErr:    false,
			statusCode: http.StatusOK,
		},
		{
			name: "Reject request with non-numeric limit parameter",
			req: func() *http.Request {
				req := httptest.NewRequest(http.MethodGet, "/tasks?limit=invalid", nil)
				req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
				return req
			},
			prepareMock: func(ctx context.Context) *mock.MockTaskRepository {
				return nil
			},
			wantErr:    false,
			statusCode: http.StatusBadRequest,
		},
		{
			name: "Reject request with non-numeric offse parameter",
			req: func() *http.Request {
				req := httptest.NewRequest(http.MethodGet, "/tasks?offset=invalid", nil)
				req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
				return req
			},
			prepareMock: func(ctx context.Context) *mock.MockTaskRepository {
				return nil
			},
			wantErr:    false,
			statusCode: http.StatusBadRequest,
		},
		{
			name: "Encounter internal server error",
			req: func() *http.Request {
				req := httptest.NewRequest(http.MethodGet, "/tasks", nil)
				req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
				return req
			},
			prepareMock: func(ctx context.Context) *mock.MockTaskRepository {
				ctrl := gomock.NewController(t)
				m := mock.NewMockTaskRepository(ctrl)
				m.EXPECT().ListCtx(gomock.Any(), 0, 10).Return(nil, errors.New("any error"))
				return m
			},
			wantErr:    false,
			statusCode: http.StatusInternalServerError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := echo.New()
			m := tt.prepareMock(context.Background())
			th := task.NewTaskHandler(m)
			rc := httptest.NewRecorder()
			c := e.NewContext(tt.req(), rc)

			if err := th.List(c); err != nil {
				t.Errorf("TaskHandler.Create() error = %v, wantErr %v", err, tt.wantErr)
			}
			assert.Equal(t, tt.statusCode, rc.Code)

			if tt.checkRes {
				var tasks entity_task.Tasks
				if err := json.Unmarshal(rc.Body.Bytes(), &tasks); err != nil {
					t.Fatal(err)
				}
				assert.Equal(t, tasks, tasksInDB)
			}
		})
	}
}

var taskInDB = entity_task.Task{ID: 1}

func TestTaskHandler_Get(t *testing.T) {
	e := echo.New()

	tests := []struct {
		name        string
		id          string
		req         func() *http.Request
		prepareMock func() *mock.MockTaskRepository
		wantErr     bool
		statusCode  int
		resBody     entity_task.Task
	}{
		{
			name: "Fetch task successfully",
			id:   "1",
			req: func() *http.Request {
				return httptest.NewRequest(http.MethodGet, "/tasks/1", nil)
			},
			prepareMock: func() *mock.MockTaskRepository {
				ctrl := gomock.NewController(t)
				m := mock.NewMockTaskRepository(ctrl)
				m.EXPECT().GetCtx(gomock.Any(), 1).Return(&taskInDB, nil)
				return m
			},
			wantErr:    false,
			statusCode: http.StatusOK,
			resBody:    taskInDB,
		},
		{
			name: "Invalid request",
			id:   "hoge",
			req: func() *http.Request {
				return httptest.NewRequest(http.MethodGet, "/tasks/hoge", nil)
			},
			prepareMock: func() *mock.MockTaskRepository { return nil },
			wantErr:     false,
			statusCode:  http.StatusBadRequest,
		},
		{
			name: "Internal server error",
			id:   "1",
			req: func() *http.Request {
				return httptest.NewRequest(http.MethodGet, "/tasks/1", nil)
			},
			prepareMock: func() *mock.MockTaskRepository {
				ctrl := gomock.NewController(t)
				m := mock.NewMockTaskRepository(ctrl)
				m.EXPECT().GetCtx(gomock.Any(), 1).Return(nil, errors.New("error"))
				return m
			},
			statusCode: http.StatusInternalServerError,
		},
		{
			name: "Not found",
			id:   "1",
			req: func() *http.Request {
				return httptest.NewRequest(http.MethodGet, "/tasks/1", nil)
			},
			prepareMock: func() *mock.MockTaskRepository {
				ctrl := gomock.NewController(t)
				m := mock.NewMockTaskRepository(ctrl)
				m.EXPECT().GetCtx(gomock.Any(), 1).Return(nil, apperrors.ErrNotFound)
				return m
			},
			wantErr:    false,
			statusCode: http.StatusNotFound,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := tt.prepareMock()
			th := task.NewTaskHandler(m)
			rc := httptest.NewRecorder()
			c := e.NewContext(tt.req(), rc)
			c.SetParamNames("id")
			c.SetParamValues(tt.id)

			if err := th.Get(c); (err != nil) != tt.wantErr {
				t.Errorf("TaskHandler.Get() error = %v, wantErr %v", err, tt.wantErr)
			}

			res := rc.Result()
			assert.Equal(t, tt.statusCode, res.StatusCode)
			if tt.resBody != (entity_task.Task{}) {
				var got entity_task.Task
				var buff bytes.Buffer
				defer func() {
					_ = res.Body.Close()
				}()
				_, _ = io.Copy(&buff, res.Body)
				if err := json.Unmarshal(buff.Bytes(), &got); err != nil {
					t.Fatal(err)
				}
				assert.Equal(t, tt.resBody, got)
			}
		})
	}
}
