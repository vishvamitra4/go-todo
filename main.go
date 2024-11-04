package main

import (
	"log"
	"my-go-sever/database/mongodb"
	"my-go-sever/routes"
	"net/http"
	"os"

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
