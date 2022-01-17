package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jonathan-innis/go-rest-tutorial/pkg/controllers"
	"github.com/jonathan-innis/go-rest-tutorial/pkg/database"
)

func main() {
	r := mux.NewRouter()
	collection := database.ConnectDB()

	db := database.NewMongoDB(collection)
	bc := controllers.NewBookController(db)

	r.HandleFunc("/api/books", bc.GetBooks).Methods("GET")
	// r.HandleFunc("/api/books/{id}", getBook).Methods("GET")
	r.HandleFunc("/api/books", bc.CreateBook).Methods("POST")
	// r.HandleFunc("/api/books/{id}", updateBook).Methods("PUT")
	// r.HandleFunc("/api/books/{id}", deleteBook).Methods("DELETE")

	// set our port address
	log.Fatal(http.ListenAndServe(":8000", r))
}
