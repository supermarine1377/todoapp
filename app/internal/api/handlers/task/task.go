// package task はTask関係のAPIハンドラーを提供する
package task

import (
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/supermarine1377/todoapp/app/common/apperrors"
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
	Create() error
}

// NewTaskHandler はTaskHandler を生成する
func NewTaskHandler(tr TaskRepository) *TaskHandler {
	return &TaskHandler{
		tr: tr,
	}
}

// Create はTaskを登録する
func (th *TaskHandler) Create(c echo.Context) error {
	if err := th.tr.Create(); err != nil {
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
