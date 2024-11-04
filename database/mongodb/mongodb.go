package mongodb

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var Client *mongo.Client // mongo client instance..

func ConnectMongoDb() *mongo.Client { // initialize the mongodb connection and returning mongo-clienr pointer...

	if err := godotenv.Load(); err != nil { // env get loaded or not..
		log.Fatal("Error loading env file..")
	}

	uri := os.Getenv("MONGODB_URI")
	if uri == "" {
		log.Fatal("Mongodb uri is not found in env variables..")
	}

	clientOptions := options.Client().ApplyURI(uri)
	// creating context with 10 second timeout..
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatal("Failed to create mongoDB client..")
	}

	// connection is active or not by sending ping
	if err := client.Ping(ctx, nil); err != nil {
		log.Fatal("Failed to connect mongodb")
	}

	fmt.Println("Connected to DB..")
	Client = client
	return client
}

// getting a particular collection from here...
func GetCollection(collectionName string) *mongo.Collection {
	dbName := os.Getenv("MONGODB_DB")
	if dbName == "" {
		log.Fatal("MONGODB_DB not found in environment variables")
	}
	return Client.Database(dbName).Collection(collectionName)
}
