// db はデータベースに対する操作を実装する
package db

// TaskRepository はTaskのリポジトリ
type TaskRepository struct {
}

// NewTaskRepository はTaskRepositoryを生成する
func NewTaskRepository() *TaskRepository {
	return &TaskRepository{}
}

// Create はTaskを作成する
func (tr TaskRepository) Create() error {
	return nil
}
