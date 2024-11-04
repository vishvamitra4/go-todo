package repositories

import (
	"context"
	"log"
	"my-go-sever/database/mongodb"
	"my-go-sever/models"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// var (
// 	todos      = []models.Todo{}
// 	todosMutex sync.Mutex
// )

type TodoRepository struct{}

func NewTodoRepository() *TodoRepository {
	return &TodoRepository{}
}

// creating a particular todo...

func (r *TodoRepository) Create(todo models.Todo) models.Todo {
	collection := mongodb.GetCollection("todos")

	// setting the context...
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	newTodo := models.Todo{
		ID:          primitive.NewObjectID(),
		Title:       todo.Title,
		Desc:        todo.Desc,
		Status:      "Pending",
		CreatedAt:   time.Now().UTC(),
		EffortHours: 0,
	}
	_, err := collection.InsertOne(ctx, newTodo)

	if err != nil {
		log.Fatal("Error while creating new todo", err)
	}
	return newTodo
}

// function for getting a all todos..

func (r *TodoRepository) GetAll() []models.Todo {
	collection := mongodb.GetCollection("todos")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cur, err := collection.Find(ctx, bson.D{})
	if err != nil {
		return nil
	}
	defer cur.Close(ctx)
	var todos []models.Todo
	if err := cur.All(ctx, &todos); err != nil {
		log.Fatal("Error while fetching all todos", err)
		return nil
	}

	return todos

}

// function for getting a particular todo by its id...

func (r *TodoRepository) GetByID(id string) (models.Todo, bool) {
	collection := mongodb.GetCollection("todos")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	var todo models.Todo
	objectid, _ := primitive.ObjectIDFromHex(id)
	err := collection.FindOne(ctx, bson.M{"_id": objectid}).Decode(&todo)
	if err != nil {
		log.Fatal("Error while getting a particular todo by id", err)
		return models.Todo{}, false
	}

	return todo, true
}

// function for updating a particular todo..
func (r *TodoRepository) Update(id string, updatedTodo models.Todo) (any, bool) {
	collection := mongodb.GetCollection("todos")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	update := bson.M{
		"$set": bson.M{
			"title":       updatedTodo.Title,
			"description": updatedTodo.Desc,
		},
	}
	objid, _ := primitive.ObjectIDFromHex(id)
	result, err := collection.UpdateOne(ctx, bson.M{"_id": objid}, update)
	if err != nil {
		log.Fatal("Error while updating..", err)
	}

	return result, true

}

// function for deleting a particular todo..
func (r *TodoRepository) Delete(id string) bool {
	collection := mongodb.GetCollection("todos")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	objid, _ := primitive.ObjectIDFromHex(id)
	_, err := collection.DeleteOne(ctx, bson.M{"_id": objid})
	if err != nil {
		log.Fatal("Error while deleting..", err)
	}

	return true
}
