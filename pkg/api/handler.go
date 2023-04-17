package api

import (
	"encoding/json"
	"fmt"
	booksModel "github.com/krukkrz/pagination/pkg/books/model"
	carsModel "github.com/krukkrz/pagination/pkg/cars/model"
	"log"
	"net/http"
	"strconv"
)

type BookRepository interface {
	FetchAll(limit, offset int) ([]booksModel.Book, error)
}

type CarRepository interface {
	FetchAll(cursor, limit int) ([]carsModel.Car, error)
}

type Server struct {
	bookRepository BookRepository
	carRepository  CarRepository
}

type PaginatedResponse[T any] struct {
	Data  []T           `json:"data"`
	Links LinksResponse `json:"links"`
}

type LinksResponse struct {
	Prev  string `json:"prev"`
	Next  string `json:"next"`
	First string `json:"first"`
	//Last  string //todo implement last link
}

func NewServer(bookRepository BookRepository, carRepository CarRepository) *Server {
	return &Server{
		bookRepository: bookRepository,
		carRepository:  carRepository,
	}
}

func (s Server) Start(port string) error {
	mux := http.NewServeMux()
	mux.HandleFunc("/books", s.FetchAllBooks)
	log.Printf("Application is ready to listen on port: %s", port)
	return http.ListenAndServe(port, mux)
}

func (s Server) FetchAllBooks(rw http.ResponseWriter, r *http.Request) {
	log.Printf("received a request: %s", r.RequestURI)
	if r.Method != "GET" {
		rw.WriteHeader(http.StatusMethodNotAllowed)
	}

	limit, err := strconv.Atoi(r.URL.Query().Get("limit"))
	if err != nil {
		rw.WriteHeader(http.StatusBadRequest)
	}

	offset, err := strconv.Atoi(r.URL.Query().Get("offset"))
	if err != nil {
		rw.WriteHeader(http.StatusBadRequest)
	}

	log.Printf("received a request with limit: %d and offset: %d", limit, offset)

	books, err := s.bookRepository.FetchAll(limit, offset)
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
	}

	nextOffset, prevOffset := offset+limit, offset-limit
	if prevOffset < 0 {
		prevOffset = 0
	}
	response := PaginatedResponse[booksModel.Book]{
		Data: books,
		Links: LinksResponse{
			Next:  fmt.Sprintf("%s?limit=%d&offset=%d", r.URL.Path, limit, nextOffset),
			Prev:  fmt.Sprintf("%s?limit=%d&offset=%d", r.URL.Path, limit, prevOffset),
			First: fmt.Sprintf("%s?limit=%d&offset=%d", r.URL.Path, limit, 0),
		},
	}

	rw.Header().Set("Content-Type", "application/json")
	json.NewEncoder(rw).Encode(response)
}
