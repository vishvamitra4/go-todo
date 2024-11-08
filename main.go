package main

import (
	"log"
	"my-go-sever/database/clickhouse"
	"my-go-sever/database/mongodb"
	"my-go-sever/routes"
	"net/http"
	"os"

	_ "github.com/ClickHouse/clickhouse-go/v2"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/joho/godotenv"
)

func main() {

	if err := godotenv.Load(); err != nil {
		log.Fatal("error while loading env files..")
	}
	// connecting DB
	mongodb.ConnectMongoDb()

	// conecting clickhouse..
	client, err := clickhouse.NewClickhouseClient(os.Getenv("CH_HOST"), os.Getenv("CH_PORT"), os.Getenv("CH_USERNAME"), os.Getenv("CH_PASSWORD"), os.Getenv("CH_DB"))
	if err != nil {
		log.Fatalf("Error creating ClickHouse client: %v", err)
	}
	if client == nil {
		log.Fatal("ClickHouse client is nil")
	}

	defer client.Close()

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	// Load routes
	routes.SetupRoutes(r)

	PORT := os.Getenv("PORT")
	if PORT == "" {
		PORT = ":3333"
	}

	log.Println("Server started on port", PORT)
	http.ListenAndServe(PORT, r)
}
