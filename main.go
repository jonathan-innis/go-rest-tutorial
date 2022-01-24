package main

import (
	"context"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jonathan-innis/go-rest-tutorial/pkg/controllers"
	"github.com/jonathan-innis/go-rest-tutorial/pkg/database"
	"github.com/jonathan-innis/go-rest-tutorial/pkg/middleware"
)

func main() {
	r := mux.NewRouter()
	collection := database.ConnectDB(context.Background(), "books")

	db := database.NewMongoDB(collection)
	bc := controllers.NewBookController(db)

	r.Use(middleware.LoggingMiddleware)

	r.HandleFunc("/api/books", bc.GetBooks).Methods("GET")
	r.HandleFunc("/api/books/{id}", bc.GetBook).Methods("GET")
	r.HandleFunc("/api/books", bc.CreateBook).Methods("POST")
	r.HandleFunc("/api/books/{id}", bc.UpdateBook).Methods("PUT")
	r.HandleFunc("/api/books/{id}", bc.DeleteBook).Methods("DELETE")

	// set our port address
	log.Fatal(http.ListenAndServe(":8000", r))
}
