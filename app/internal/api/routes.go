package api

import (
	"net/http"

	"github.com/supermarine1377/todoapp/app/internal/api/handlers/healthz"
	"github.com/supermarine1377/todoapp/app/internal/api/handlers/task"
	"github.com/supermarine1377/todoapp/app/internal/api/server"
	"github.com/supermarine1377/todoapp/app/internal/repository"
)

// Route は、APIのルーティングを設定する
func Route(s *server.Server) {
	s.RegisterHandler(healthz.Healthz, "/healthz", http.MethodGet)
	{
		tr := repository.NewTaskRepository(s.DB())
		th := task.NewTaskHandler(tr)
		s.RegisterHandler(th.Create, "/tasks", http.MethodPost)
		s.RegisterHandler(th.List, "/tasks", http.MethodGet)
		s.RegisterHandler(th.Get, "/tasks/:id", http.MethodGet)
	}
}
