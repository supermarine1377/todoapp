// db はデータベースに対する操作を実装する
package db

// TODORepository はTODOのリポジトリ
type TODORepository struct {
}

// NewTODORepository はTODORepositoryを生成する
func NewTODORepository() *TODORepository {
	return &TODORepository{}
}

// Create はTODOを作成する
func (tr TODORepository) Create() error {
	return nil
}
