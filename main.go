package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"sync"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

// Response struct for backend
type Response struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Data    any    `json:"data"`
}

// Todo represents a simple ToDo item
type Todo struct {
	ID    int    `json:"id"`
	Title string `json:"title"`
	Desc  string `json:"desc"`
}

var (
	todos      = []Todo{} // slice to store todos
	todoID     = 1        // incremental ID for new todos
	todosMutex sync.Mutex // mutex to synchronize access to todos
)

func main() {
	r := chi.NewRouter()     // here is new route..
	r.Use(middleware.Logger) //middlewars...
	r.Use(middleware.Recoverer)
	// Routes for getting status..
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Vishva your backend working properly.."))
	})
	// Routes crud operation...
	r.Post("/todos", createTodo)
	r.Get("/todos", getAllTodos)
	r.Get("/todos/{id}", getTodo)
	r.Put("/todos/{id}", updateTodo)
	r.Delete("/todos/{id}", deleteTodo)

	log.Println(" server started on port 3333")
	http.ListenAndServe(":3333", r) // server stated on port 3333...
}

// createTodo handles creating a new ToDo
func createTodo(w http.ResponseWriter, r *http.Request) {
	var todo Todo                                // here we are defining type of todo..
	err := json.NewDecoder(r.Body).Decode(&todo) // decoding the values from request..
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}
	if todo.Title == "" || todo.Desc == "" {
		http.Error(w, "title or desc is missing into this...", http.StatusBadRequest)
		return
	}

	todosMutex.Lock()
	todo.ID = todoID            // setting the todoid..
	todoID++                    // increment the todoID..
	todos = append(todos, todo) // appending the value into todos..
	todosMutex.Unlock()

	response := Response{
		Success: true,
		Message: "New Todo is created..",
		Data:    todo,
	}
	w.WriteHeader(http.StatusCreated) // 201...
	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Println("Something unexpected wring while creating nre todo..")
		http.Error(w, "Internal Error", http.StatusInternalServerError)
	}
}

// getAllTodos handles retrieving all ToDos
func getAllTodos(w http.ResponseWriter, r *http.Request) {
	todosMutex.Lock()
	defer todosMutex.Unlock()

	// creating response...
	response := Response{
		Success: true,
		Message: "All Todo is fetched..",
		Data:    todos,
	}

	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Something Wrong..", http.StatusInternalServerError)
	}
}

// gettinhg a particular todo by its id..
func getTodo(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id")) // getting id and converting into integer..
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}
	log.Println(id)

	todosMutex.Lock()
	defer todosMutex.Unlock()

	for _, todo := range todos {
		if todo.ID == id {
			json.NewEncoder(w).Encode(todo)
			return
		}
	}

	http.Error(w, "ToDo not found", http.StatusNotFound)

}

// updateTodo handles updating a specific ToDo by ID
func updateTodo(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	var updatedTodo Todo
	if err := json.NewDecoder(r.Body).Decode(&updatedTodo); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	log.Println(updatedTodo)

	todosMutex.Lock()
	defer todosMutex.Unlock() // making sure mutex get unblocked...
	for i, todo := range todos {
		if todo.ID == id {
			todos[i].Title = updatedTodo.Title
			todos[i].Desc = updatedTodo.Desc
			json.NewEncoder(w).Encode(todos[i])
			return
		}
	}

	http.Error(w, "ToDo not found", http.StatusNotFound)
}

// deleteTodo handles deleting a specific ToDo by ID
func deleteTodo(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	todosMutex.Lock()
	defer todosMutex.Unlock()
	for i, todo := range todos {
		if todo.ID == id {
			todos = append(todos[:i], todos[i+1:]...) // deleting a particular document.
			json.NewEncoder(w).Encode("Deleted!")
			return
		}
	}

	http.Error(w, "ToDo not found", http.StatusNotFound)
}
