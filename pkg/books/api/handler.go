package api

import (
	"encoding/json"
	"github.com/krukkrz/pagination/pkg/books/model"
	"log"
	"net/http"
	"strconv"
)

type BookRepository interface {
	FetchAll(limit, offset int) ([]model.Book, error)
}

type Server struct {
	repository BookRepository
}

func NewServer(repository BookRepository) *Server {
	return &Server{
		repository: repository,
	}
}

func (s Server) FetchAllBooks(rw http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		rw.WriteHeader(http.StatusMethodNotAllowed)
	}

	limit, err := strconv.Atoi(r.FormValue("limit"))
	if err != nil {
		rw.WriteHeader(http.StatusBadRequest)
	}

	offset, err := strconv.Atoi(r.FormValue("offset"))
	if err != nil {
		rw.WriteHeader(http.StatusBadRequest)
	}

	log.Printf("received a request with limit: %d and offset: %d", limit, offset)

	books, err := s.repository.FetchAll(limit, offset)
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
	}

	rw.Header().Set("Content-Type", "application/json")
	json.NewEncoder(rw).Encode(books)
}
