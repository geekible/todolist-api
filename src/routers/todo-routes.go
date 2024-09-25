package routers

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

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
	m.mux.Put(m.baseEndpoint, m.updateTodo)
	m.mux.Delete(m.baseEndpoint, m.deleteTodo)

	// get's
	m.mux.Get(m.baseEndpoint, m.getById)
	m.mux.Get(fmt.Sprintf("%s/get-by-user-id", m.baseEndpoint), m.getByUserId)
	m.mux.Get(fmt.Sprintf("%s/get-model", m.baseEndpoint), func(w http.ResponseWriter, r *http.Request) {
		m.jsonHelper.WriteJSON(w, http.StatusOK, domain.ToDoEntity{})
	})
}

func (m *ToDoRoutes) createTodo(w http.ResponseWriter, r *http.Request) {
	var todo domain.ToDoEntity
	if err := m.jsonHelper.ReadJSON(w, r, &todo); err != nil {
		m.jsonHelper.ErrorJSON(w, err, http.StatusBadRequest, m.baseEndpoint)
		return
	}

	todo, err := m.todoService.Add(todo)
	if err != nil {
		m.jsonHelper.ErrorJSON(w, err, http.StatusBadRequest, m.baseEndpoint)
		return
	}

	m.jsonHelper.WriteJSON(w, http.StatusCreated, todo, nil)
}

func (m *ToDoRoutes) updateTodo(w http.ResponseWriter, r *http.Request) {
	var todo domain.ToDoEntity
	if err := m.jsonHelper.ReadJSON(w, r, &todo); err != nil {
		m.jsonHelper.ErrorJSON(w, err, http.StatusBadRequest, m.baseEndpoint)
		return
	}

	if err := m.todoService.Update(todo); err != nil {
		m.jsonHelper.ErrorJSON(w, err, http.StatusBadRequest, m.baseEndpoint)
		return
	}

	m.jsonHelper.WriteJSON(w, http.StatusCreated, todo, nil)
}

func (m *ToDoRoutes) deleteTodo(w http.ResponseWriter, r *http.Request) {
	var todo domain.ToDoEntity
	if err := m.jsonHelper.ReadJSON(w, r, &todo); err != nil {
		m.jsonHelper.ErrorJSON(w, err, http.StatusBadRequest, m.baseEndpoint)
		return
	}

	if err := m.todoService.Delete(todo); err != nil {
		m.jsonHelper.ErrorJSON(w, err, http.StatusBadRequest, m.baseEndpoint)
		return
	}

	m.jsonHelper.WriteJSON(w, http.StatusCreated, todo, nil)
}

func (m *ToDoRoutes) getById(w http.ResponseWriter, r *http.Request) {
	userId, err := strconv.Atoi(r.URL.Query().Get("userId"))
	if err != nil {
		m.jsonHelper.ErrorJSON(w, errors.New("userId must be supplied"), http.StatusBadRequest, m.baseEndpoint)
		return
	}

	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil {
		m.jsonHelper.ErrorJSON(w, errors.New("userId must be supplied"), http.StatusBadRequest, m.baseEndpoint)
		return
	}

	todo, err := m.todoService.GetById(uint(id), uint(userId))
	if err != nil {
		m.jsonHelper.ErrorJSON(w, err, http.StatusBadRequest, m.baseEndpoint)
		return
	}

	m.jsonHelper.WriteJSON(w, http.StatusOK, todo, nil)
}

func (m *ToDoRoutes) getByUserId(w http.ResponseWriter, r *http.Request) {
	userId, err := strconv.Atoi(r.URL.Query().Get("userId"))
	if err != nil {
		m.jsonHelper.ErrorJSON(w, errors.New("userId must be supplied"), http.StatusBadRequest, m.baseEndpoint)
		return
	}

	todo, err := m.todoService.GetByUserId(uint(userId), 0, 100)
	if err != nil {
		m.jsonHelper.ErrorJSON(w, err, http.StatusBadRequest, m.baseEndpoint)
		return
	}

	m.jsonHelper.WriteJSON(w, http.StatusOK, todo, nil)
}
