package controllers

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jonathan-innis/go-rest-tutorial/pkg/database"
	"github.com/jonathan-innis/go-rest-tutorial/pkg/helper"
	"github.com/jonathan-innis/go-rest-tutorial/pkg/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type BookController struct {
	db database.Interface
}

func NewBookController(db database.Interface) *BookController {
	return &BookController{db: db}
}

func (bc *BookController) CreateBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var book models.Book
	if err := json.NewDecoder(r.Body).Decode(&book); err != nil {
		helper.GetError(w, http.StatusBadRequest, "Invalid request payload with error: "+err.Error())
		return
	}
	book.ID = primitive.NilObjectID

	createdID, err := bc.db.Create(context.TODO(), &book)
	if err != nil {
		helper.GetInternalError(w, err)
		return
	}

	// Update the ID of the book that we just created
	book.ID = createdID
	w.Header().Add("Location", r.Host+"/api/books/"+book.ID.Hex())
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(book)
}

func (bc *BookController) UpdateBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	book := &models.Book{}

	if err := json.NewDecoder(r.Body).Decode(book); err != nil {
		helper.GetError(w, http.StatusBadRequest, "Invalid request payload with error: "+err.Error())
		return
	}

	// Remove id from payload
	book.ID = primitive.NilObjectID

	params := mux.Vars(r)
	if id, ok := params["id"]; ok {
		// Validate if the id is a valid ObjectID
		if !primitive.IsValidObjectID(id) {
			helper.GetError(w, http.StatusBadRequest, "ID specified must be a valid ObjectID")
			return
		}

		if err := bc.db.Update(context.TODO(), book, id); err != nil {
			helper.GetInternalError(w, err)
			return
		}
		oid, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			helper.GetInternalError(w, err)
			return
		}
		book.ID = oid
		w.WriteHeader(http.StatusAccepted)
		json.NewEncoder(w).Encode(book)
		return
	}
	helper.GetError(w, http.StatusBadRequest, "ID is required")
}

func (bc *BookController) GetBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	params := mux.Vars(r)
	if id, ok := params["id"]; ok {
		// Validate if the id is a valid ObjectID
		if !primitive.IsValidObjectID(id) {
			helper.GetError(w, http.StatusBadRequest, "ID specified must be a valid ObjectID")
			return
		}

		book := &models.Book{}
		found, err := bc.db.Get(context.TODO(), id, book)
		if err != nil {
			helper.GetInternalError(w, err)
			return
		} else if !found {
			helper.GetError(w, http.StatusNotFound, "Book not found")
			return
		}
		json.NewEncoder(w).Encode(book)
		return
	}
	helper.GetError(w, http.StatusBadRequest, "ID is required")
}

func (bc *BookController) GetBooks(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var books []models.Book
	err := bc.db.List(context.TODO(), &books)
	if err != nil {
		helper.GetInternalError(w, err)
	}
	json.NewEncoder(w).Encode(books)
}

func (bc *BookController) DeleteBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	params := mux.Vars(r)
	if id, ok := params["id"]; ok {
		if found, err := bc.db.Delete(context.TODO(), id); err != nil {
			helper.GetInternalError(w, err)
			return
		} else if !found {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		w.WriteHeader(http.StatusOK)
		return
	}
	helper.GetError(w, http.StatusBadRequest, "ID is required")
}
