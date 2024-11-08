package repositories

import (
	"context"
	"log"
	"my-go-sever/database/clickhouse"
	"my-go-sever/database/mongodb"
	"my-go-sever/models"
	"os"
	"time"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// TodoRepository represents a repository for managing todos in ClickHouse.
type TodoRepository struct {
}

func NewTodoRepository() *TodoRepository {
	return &TodoRepository{}
}

// creating a particular todo...

func (r *TodoRepository) Create(todo models.Todo) models.Todo {

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// storing data into clickhouse..
	client, _ := clickhouse.NewClickhouseClient(os.Getenv("CH_HOST"), os.Getenv("CH_PORT"), os.Getenv("CH_USERNAME"), os.Getenv("CH_PASSWORD"), os.Getenv("CH_DB"))

	if client == nil {
		log.Fatal("ClickHouse client is not initialized")
	}

	query := `
        INSERT INTO test_db.todos (id, title, desc, status, created_at, effort_hours)
        VALUES (?, ?, ?, ?, ?, ?)
    `

	_, e := client.ExecContext(ctx, query,
		uuid.New().String(),
		todo.Title,
		todo.Desc,
		todo.Status,
		time.Now().UTC(),
		todo.EffortHours,
	)

	if e != nil {
		log.Fatalf("Error while creating new todo: %v", e)
	}

	// storing data into mongo...

	collection := mongodb.GetCollection("todos")

	newTodo := models.Todo{
		ID:          primitive.NewObjectID(),
		Title:       todo.Title,
		Desc:        todo.Desc,
		Status:      todo.Status,
		CreatedAt:   time.Now().UTC(),
		EffortHours: todo.EffortHours,
	}

	_, e1 := collection.InsertOne(ctx, newTodo)

	if e1 != nil {
		log.Fatal("Error while creating new todo", e1)
	}

	return newTodo

}

// function for getting a all todos..

func (r *TodoRepository) GetAll(flag string) []models.Todo { // flag === true means data is coming from clickhouse
	if flag == "true" {
		return r.GetTodosFromMongo()
	}
	return r.getTodosFromClickHouse()
}

// getting all todos from mongo..
func (r *TodoRepository) GetTodosFromMongo() []models.Todo {
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

// Fetch todos from ClickHouse
func (r *TodoRepository) getTodosFromClickHouse() []models.Todo {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	client, _ := clickhouse.NewClickhouseClient(os.Getenv("CH_HOST"), os.Getenv("CH_PORT"), os.Getenv("CH_USERNAME"), os.Getenv("CH_PASSWORD"), os.Getenv("CH_DB"))
	if client == nil {
		log.Fatal("ClickHouse client is not initialized")
		return nil
	}

	query := `SELECT id, title, desc, status, created_at, effort_hours FROM todos`

	rows, err := client.QueryContext(ctx, query)
	if err != nil {
		log.Fatal("Error while fetching todos from ClickHouse", err)
		return nil
	}
	defer rows.Close()

	var todos []models.Todo
	for rows.Next() {
		var todo models.Todo
		err := rows.Scan(&todo.ID, &todo.Title, &todo.Desc, &todo.Status, &todo.CreatedAt, &todo.EffortHours)
		if err != nil {
			log.Fatal("Error while scanning ClickHouse rows", err)
			return nil
		}
		todos = append(todos, todo)
	}

	if err := rows.Err(); err != nil {
		log.Fatal("Error while iterating over ClickHouse rows", err)
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

//

func (r *TodoRepository) AggregateMetrics() (interface{}, error) {
	collection := mongodb.GetCollection("todos")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	pipeline := mongo.Pipeline{
		// Counting total number of todos
		bson.D{
			{Key: "$group", Value: bson.D{
				{Key: "_id", Value: nil},
				{Key: "totalCount", Value: bson.D{{Key: "$sum", Value: 1}}},
			}},
		},
		bson.D{
			{Key: "$project", Value: bson.D{
				{Key: "_id", Value: 0},
				{Key: "totalCount", Value: 1},
			}},
		},
	}

	cur, error := collection.Aggregate(ctx, pipeline)
	if error != nil {
		log.Fatal("Error while aggragation..")
	}
	defer cur.Close(ctx)

	// parsing the result into slice of map..
	var metrics []bson.M

	if err := cur.All(ctx, &metrics); err != nil {
		log.Fatal("Error while parsing the result...")
	}

	return metrics, nil
}
