package services

import (
	"my-go-sever/models"
	"my-go-sever/repositories"
)

type TodoService struct {
	repo *repositories.TodoRepository
}

func NewTodoService(repo *repositories.TodoRepository) *TodoService {
	return &TodoService{repo: repo}
}

func (s *TodoService) Create(todo models.Todo) models.Todo {
	return s.repo.Create(todo)
}

func (s *TodoService) GetAll() []models.Todo {
	return s.repo.GetAll()
}

func (s *TodoService) GetByID(id string) (models.Todo, bool) {
	return s.repo.GetByID(id)
}

func (s *TodoService) Update(id string, updatedTodo models.Todo) (any, bool) {
	return s.repo.Update(id, updatedTodo)
}

func (s *TodoService) Delete(id string) bool {
	return s.repo.Delete(id)
}

func (s *TodoService) GetTodoMetrics() (interface{}, error) {
	return s.repo.AggregateMetrics()
}
