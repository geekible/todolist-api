package routers

import (
	"net/http"

	"geekible.todolist/src/config"
	"geekible.todolist/src/domain"
	"geekible.todolist/src/helpers"
	"geekible.todolist/src/services"
	"github.com/go-chi/chi/v5"
)

type ToDoRoutes struct {
	baseEndpoint string
	mux          *chi.Mux
	todoService  *services.ToDoService
	jsonHelper   *helpers.JsonHelpers
}

func InitToDoRoutes(mux *chi.Mux, cfg *config.ServiceConfig) *ToDoRoutes {
	return &ToDoRoutes{
		todoService:  services.InitToDoService(cfg),
		jsonHelper:   helpers.InitJsonHelpers(),
		baseEndpoint: "/todo",
		mux:          mux,
	}
}

func (m *ToDoRoutes) RegisterRoutes() {
	m.mux.Post(m.baseEndpoint, m.createTodo)
}

func (m *ToDoRoutes) createTodo(w http.ResponseWriter, r *http.Request) {
	var todo domain.ToDoEntity
	if err := m.jsonHelper.ReadJSON(w, r, &todo); err != nil {
		m.jsonHelper.ErrorJSON(w, err, http.StatusBadRequest, m.baseEndpoint)
	}

	todo, err := m.todoService.Add(todo)
	if err != nil {
		m.jsonHelper.ErrorJSON(w, err, http.StatusBadRequest, m.baseEndpoint)
	}

	m.jsonHelper.WriteJSON(w, http.StatusCreated, todo, nil)
}
