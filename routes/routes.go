package routes

import (
	"my-go-sever/controllers"
	"my-go-sever/repositories"
	"my-go-sever/services"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func SetupRoutes(r *chi.Mux) {
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Krishna your backend working properly.."))
	})

	// Initialize repository, service, and controller
	repo := repositories.NewTodoRepository()
	service := services.NewTodoService(repo)
	controller := controllers.NewTodoController(service)

	r.Route("/todos", func(r chi.Router) {
		r.Post("/", controller.CreateTodoHandler)
		r.Get("/", controller.GetAllTodosHandler)
		r.Get("/{id}", controller.GetTodoHandler)
		r.Put("/{id}", controller.UpdateTodoHandler)
		r.Delete("/{id}", controller.DeleteTodoHandler)
		r.Get("/metrics", controller.GetTodoMetricsController)
	})
}
