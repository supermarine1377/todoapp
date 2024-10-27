// package todo はTODO関係のAPIハンドラーを提供する
package todo

import (
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/supermarine1377/todoapp/app/common/apperrors"
)

// TODOHandler はTODO関係のAPIハンドラーを表す
type TODOHandler struct {
	tr TODORepository
}

// TODORepository はTODOのリポジトリを抽象化する
//
//go:generate mockgen -source=todo.go -destination=./mock/mock.go -package=mock
type TODORepository interface {
	// TODO を作成する
	Create() error
}

// NewTODOHandler はTODOHandler を生成する
func NewTODOHandler(tr TODORepository) *TODOHandler {
	return &TODOHandler{
		tr: tr,
	}
}

// Create はTODOを登録する
func (th *TODOHandler) Create(c echo.Context) error {
	if err := th.tr.Create(); err != nil {
		if errors.Is(err, apperrors.ErrBadRequest) {
			return c.JSON(http.StatusBadRequest, nil)
		}
		if errors.Is(err, apperrors.ErrInternalServerError) {
			c.Logger().Error("Failed to create TODO: %w", err)
			return c.JSON(http.StatusInternalServerError, nil)
		}
	}
	return c.JSON(http.StatusCreated, nil)
}
