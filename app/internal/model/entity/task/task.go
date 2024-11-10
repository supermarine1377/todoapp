// package task はタスクを実装する
package task

// ID はタスクのID
type ID int64

// Task はタスクを表す
type Task struct {
	ID        ID     `json:"id"`
	Title     string `validate:"required" json:"title"`
	CreatedAt int64  `json:"created_at"`
	UpdatedAt int64  `json:"updated_at"`
}

type Tasks []*Task
