package controllers

import (
	"encoding/json"
	"my-go-sever/models"
	"my-go-sever/services"
	"my-go-sever/utils"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type TodoController struct {
	service *services.TodoService
}

func NewTodoController(service *services.TodoService) *TodoController {
	return &TodoController{service: service}
}

func (c *TodoController) CreateTodoHandler(w http.ResponseWriter, r *http.Request) {
	var todo models.Todo
	if err := json.NewDecoder(r.Body).Decode(&todo); err != nil || todo.Title == "" || todo.Desc == "" {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	createdTodo := c.service.Create(todo)
	utils.RespondWithJSON(w, http.StatusCreated, models.Response{
		Success: true,
		Message: "New Todo is created.",
		Data:    createdTodo,
	})
}

func (c *TodoController) GetAllTodosHandler(w http.ResponseWriter, r *http.Request) {
	todos := c.service.GetAll()
	utils.RespondWithJSON(w, http.StatusOK, models.Response{
		Success: true,
		Message: "All Todos fetched.",
		Data:    todos,
	})
}

func (c *TodoController) GetTodoHandler(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}
	todo, found := c.service.GetByID(id)
	if !found {
		http.Error(w, "ToDo not found", http.StatusNotFound)
		return
	}
	utils.RespondWithJSON(w, http.StatusOK, todo)
}

func (c *TodoController) UpdateTodoHandler(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}
	var updatedTodo models.Todo
	if err := json.NewDecoder(r.Body).Decode(&updatedTodo); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	todo, found := c.service.Update(id, updatedTodo)
	if !found {
		http.Error(w, "ToDo not found", http.StatusNotFound)
		return
	}
	utils.RespondWithJSON(w, http.StatusOK, todo)
}

func (c *TodoController) DeleteTodoHandler(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	if !c.service.Delete(id) {
		http.Error(w, "ToDo not found", http.StatusNotFound)
		return
	}
	utils.RespondWithJSON(w, http.StatusOK, "Deleted!")
}
