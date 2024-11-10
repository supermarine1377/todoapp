// repository はデータベースに対する操作を実装する
package repository

import (
	"context"

	"github.com/supermarine1377/todoapp/app/common/apperrors"
	"github.com/supermarine1377/todoapp/app/internal/model/entity/task"
)

// TaskRepository はTaskのリポジトリ
type TaskRepository struct {
	db DB
}

// DB はDBの実装を抽象化する
//
//go:generate mockgen -source=repository.go -destination=./mock/mock.go -package=mock
type DB interface {
	InsertCtx(ctx context.Context, p interface{}) error
	SelectCtx(ctx context.Context, p interface{}, columns []string, offset, limit int) error
}

// NewTaskRepository はTaskRepositoryを生成する
func NewTaskRepository(db DB) *TaskRepository {
	return &TaskRepository{db: db}
}

// Create はTaskを作成する
func (tr *TaskRepository) CreateCtx(ctx context.Context, task *task.Task) error {
	if err := tr.db.InsertCtx(ctx, task); err != nil {
		return apperrors.ErrInternalServerError
	}
	return nil
}

func (tr *TaskRepository) ListCtx(ctx context.Context, offset, limit int) (*task.Tasks, error) {
	var t task.Tasks
	col := []string{"id", "title", "created_at", "updated_at"}
	if err := tr.db.SelectCtx(ctx, &t, col, offset, limit); err != nil {
		return nil, err
	}
	return &t, nil
}
