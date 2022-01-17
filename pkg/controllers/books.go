package controllers

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/jonathan-innis/go-rest-tutorial/pkg/database"
	"github.com/jonathan-innis/go-rest-tutorial/pkg/helper"
	"github.com/jonathan-innis/go-rest-tutorial/pkg/models"
)

type BookController struct {
	db database.Interface
}

func NewBookController(db database.Interface) *BookController {
	return &BookController{db: db}
}

func (bc *BookController) GetBooks(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	bookModel := &models.Book{}

	res, err := bc.db.List(context.TODO(), bookModel)
	if err != nil {
		helper.GetError(err, w)
	}
	json.NewEncoder(w).Encode(res)
}

func (bc *BookController) CreateBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var book models.Book

	_ = json.NewDecoder(r.Body).Decode(&book)

	bc.db.Create(context.TODO(), &book)
}
