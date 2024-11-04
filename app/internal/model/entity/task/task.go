// package task はタスクを実装する
package task

import "time"

// ID はタスクのID
type ID int64

// Task はタスクを表す
type Task struct {
	ID        ID
	Title     string `validate:"required" json:"title"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

type Tasks []*Task
