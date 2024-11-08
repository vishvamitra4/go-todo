package controllers

import (
	"encoding/json"
	"fmt"
	"log"
	"my-go-sever/database/clickhouse"
	"my-go-sever/database/mongodb"
	"my-go-sever/models"
	"my-go-sever/services"
	"my-go-sever/utils"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"
	"go.mongodb.org/mongo-driver/bson"
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
	flag := r.URL.Query().Get("flag")
	todos := c.service.GetAll(flag)
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

// getting metrics from clickhouse...
func (c *TodoController) GetTodoMetricsControllerClick(w http.ResponseWriter, r *http.Request) {
	metrics, err := clickhouse.AggregateMetricsFromClickHouse()
	if err != nil {
		log.Println(err)
		fmt.Println("Error While Fetching Metrics from clickhouse..")
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(metrics)
}

// getting metrics from mongo...
func (c *TodoController) GetTodoMetricsControllerMongo(w http.ResponseWriter, r *http.Request) {
	collection := mongodb.GetCollection("todos")

	matchFields := bson.D{}
	for key, values := range r.URL.Query() {
		log.Printf(key, values)
		if key != "groupBy" && key != "project" && key != "sort" && key != "limit" {
			matchFields = append(matchFields, bson.E{Key: key, Value: values[0]})
		}
	}

	groupByFields := bson.D{}
	if groupByParam := r.URL.Query().Get("groupBy"); groupByParam != "" {
		groupByKeys := strings.Split(groupByParam, ",")

		idFields := bson.D{}
		for _, field := range groupByKeys {
			idFields = append(idFields, bson.E{Key: field, Value: "$" + field})
		}

		groupByFields = append(groupByFields, bson.E{Key: "_id", Value: idFields})
		groupByFields = append(groupByFields, bson.E{Key: "count", Value: bson.D{{Key: "$sum", Value: 1}}})
	}

	projectFields := bson.D{}
	if projectParam := r.URL.Query().Get("project"); projectParam != "" {
		projectKeys := strings.Split(projectParam, ",")
		for _, field := range projectKeys {
			projectFields = append(projectFields, bson.E{Key: field, Value: 1})
		}
		projectFields = append(projectFields, bson.E{Key: "_id", Value: 0})
	}

	sortFields := bson.D{}
	if sortParam := r.URL.Query().Get("sort"); sortParam != "" {
		sortKeys := strings.Split(sortParam, ",")
		for _, field := range sortKeys {
			fieldParts := strings.Split(field, ":")
			sortDirection, _ := strconv.Atoi(fieldParts[1])
			sortFields = append(sortFields, bson.E{Key: fieldParts[0], Value: sortDirection})
		}
	}

	limit := 10 // Default limit
	if limitParam := r.URL.Query().Get("limit"); limitParam != "" {
		if limitValue, err := strconv.Atoi(limitParam); err == nil {
			limit = limitValue
		}
	}

	log.Printf("Match Fields: %v", matchFields)
	log.Printf("Group By Fields: %v", groupByFields)
	log.Printf("Project Fields: %v", projectFields)
	log.Printf("Sort Fields: %v", sortFields)
	log.Printf("Limit: %d", limit)

	metrics, err := mongodb.AggregationWithOptions(collection, matchFields, groupByFields, projectFields, sortFields, limit)
	if err != nil {
		log.Printf("Aggregation error: %v", err)

		http.Error(w, "Failed to retrieve metrics due to aggregation error", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(metrics)
}
