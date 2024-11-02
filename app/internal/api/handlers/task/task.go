// package task はTask関係のAPIハンドラーを提供する
package task

import (
	"context"
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/supermarine1377/todoapp/app/common/apperrors"
	"github.com/supermarine1377/todoapp/app/internal/model/entity/task"
)

// TaskHandler はTask関係のAPIハンドラーを表す
type TaskHandler struct {
	tr TaskRepository
}

// TaskRepository はTaskのリポジトリを抽象化する
//
//go:generate mockgen -source=task.go -destination=./mock/mock.go -package=mock
type TaskRepository interface {
	// task を作成する
	CreateCtx(ctx context.Context, task *task.Task) error
}

// NewTaskHandler はTaskHandler を生成する
func NewTaskHandler(tr TaskRepository) *TaskHandler {
	return &TaskHandler{
		tr: tr,
	}
}

// Create はTaskを登録する
func (th *TaskHandler) Create(c echo.Context) error {
	ctx := c.Request().Context()
	// TODO リクエストからtaskを生成
	task := &task.Task{}
	if err := th.tr.CreateCtx(ctx, task); err != nil {
		if errors.Is(err, apperrors.ErrBadRequest) {
			return c.JSON(http.StatusBadRequest, nil)
		}
		if errors.Is(err, apperrors.ErrInternalServerError) {
			c.Logger().Error("Failed to create Task: %w", err)
			return c.JSON(http.StatusInternalServerError, nil)
		}
	}
	return c.JSON(http.StatusCreated, nil)
}
