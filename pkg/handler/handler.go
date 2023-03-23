package handler

import (
	"encoding/json"
	"github.com/krukkrz/pagination/pkg/model"
	"log"
	"net/http"
	"strconv"
)

type Service interface {
	FetchAllLimitAndOffset(int, int) ([]model.Book, error)
}

type Server struct {
	service Service
}

func NewServer(service Service) *Server {
	return &Server{
		service: service,
	}
}

func (s Server) FetchAllBooksWithLimitAndOffset(rw http.ResponseWriter, r *http.Request) {
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

	books, err := s.service.FetchAllLimitAndOffset(limit, offset)
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
	}

	rw.Header().Set("Content-Type", "application/json")
	json.NewEncoder(rw).Encode(books)
}
