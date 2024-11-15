// package task はTask関係のAPIハンドラーを提供する
package task

import (
	"context"
	"errors"
	"net/http"
	"strconv"

	"github.com/go-playground/validator/v10"
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
	// task の一覧を返却する
	ListCtx(ctx context.Context, offset, limit int) (*task.Tasks, error)
	// taskを取得する
	GetCtx(ctx context.Context, id int) (*task.Task, error)
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
	var task task.Task
	if err := c.Bind(&task); err != nil {
		switch err {
		case echo.ErrUnsupportedMediaType:
			return c.JSON(http.StatusUnsupportedMediaType, nil)
		case echo.ErrBadRequest:
			return c.JSON(http.StatusBadRequest, nil)
		default:
			return c.JSON(http.StatusBadRequest, nil)
		}
	}
	val := validator.New()
	if err := val.Struct(task); err != nil {
		return c.JSON(http.StatusBadRequest, "Missing required fileds")
	}

	if err := th.tr.CreateCtx(ctx, &task); err != nil {
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

// List はTaskの一覧を返却する
func (th *TaskHandler) List(c echo.Context) error {
	var offset int
	var err error
	if o := c.QueryParam("offset"); o != "" {
		offset, err = strconv.Atoi(o)
		if err != nil {
			return c.JSON(http.StatusBadRequest, nil)
		}
	}

	limit := 10
	if l := c.QueryParam("limit"); l != "" {
		limit, err = strconv.Atoi(l)
		if err != nil {
			return c.JSON(http.StatusBadRequest, nil)
		}
	}

	ctx := c.Request().Context()
	tasks, err := th.tr.ListCtx(ctx, offset, limit)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, nil)
	}

	return c.JSON(http.StatusOK, tasks)
}

// Get はTaskを取得する
func (th *TaskHandler) Get(c echo.Context) error {
	var id int
	var err error

	if id, err = strconv.Atoi(c.Param("id")); err != nil {
		return c.JSON(http.StatusBadRequest, nil)
	}
	ctx := c.Request().Context()
	task, err := th.tr.GetCtx(ctx, id)
	if err != nil {
		if errors.Is(err, apperrors.ErrNotFound) {
			return c.JSON(http.StatusNotFound, nil)
		}
		return c.JSON(http.StatusInternalServerError, nil)
	}
	return c.JSON(http.StatusOK, task)
}
