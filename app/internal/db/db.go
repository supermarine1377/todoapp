// db はデータベースに対する操作を実装する
package db

import (
	"context"

	"github.com/supermarine1377/todoapp/app/internal/model/entity/task"
)

// TaskRepository はTaskのリポジトリ
type TaskRepository struct {
}

// NewTaskRepository はTaskRepositoryを生成する
func NewTaskRepository() *TaskRepository {
	return &TaskRepository{}
}

// Create はTaskを作成する
func (tr TaskRepository) CreateCtx(ctx context.Context, task *task.Task) error {
	return nil
}
