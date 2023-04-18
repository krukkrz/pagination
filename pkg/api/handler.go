package api

import (
	"encoding/json"
	"fmt"
	booksModels "github.com/krukkrz/pagination/pkg/books/model"
	carsModels "github.com/krukkrz/pagination/pkg/cars/model"
	"log"
	"net/http"
	"strconv"
)

type BookRepository interface {
	FetchAll(limit, offset int) ([]booksModels.Book, error)
}

type CarRepository interface {
	FetchAll(cursor, limit int) ([]carsModels.Car, error)
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
	mux.HandleFunc("/cars", s.FetchAllCars)
	log.Printf("Application is ready to listen on port: %s", port)
	return http.ListenAndServe(port, mux)
}

func (s Server) FetchAllBooks(rw http.ResponseWriter, r *http.Request) {
	log.Printf("received a request: %s", r.RequestURI)
	validateGetRequest(rw, r)

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
	urlFormat := "%s?limit=%d&offset=%d"
	response := PaginatedResponse[booksModels.Book]{
		Data: books,
		Links: LinksResponse{
			Next:  fmt.Sprintf(urlFormat, r.URL.Path, limit, nextOffset),
			Prev:  fmt.Sprintf(urlFormat, r.URL.Path, limit, prevOffset),
			First: fmt.Sprintf(urlFormat, r.URL.Path, limit, 0),
		},
	}

	encodeJsonResponse(rw, response)
}

func (s Server) FetchAllCars(rw http.ResponseWriter, r *http.Request) {
	log.Printf("received a request: %s", r.RequestURI)
	validateGetRequest(rw, r)

	cursor, err := strconv.Atoi(r.URL.Query().Get("cursor"))
	if err != nil {
		rw.WriteHeader(http.StatusBadRequest)
	}

	limit, err := strconv.Atoi(r.URL.Query().Get("limit"))
	if err != nil {
		rw.WriteHeader(http.StatusBadRequest)
	}

	log.Printf("received a request with cursor: %d and limit: %d", cursor, limit)

	cars, err := s.carRepository.FetchAll(cursor, limit)
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
	}

	nextCursor, prevCursor := cursor+limit, cursor-limit
	if prevCursor < 1 {
		prevCursor = 1
	}
	urlFormat := "%s?cursor=%d&limit=%d"
	response := PaginatedResponse[carsModels.Car]{
		Data: cars,
		Links: LinksResponse{
			Next:  fmt.Sprintf(urlFormat, r.URL.Path, nextCursor, limit),
			Prev:  fmt.Sprintf(urlFormat, r.URL.Path, prevCursor, limit),
			First: fmt.Sprintf(urlFormat, r.URL.Path, 1, limit),
		},
	}

	encodeJsonResponse(rw, response)
}

func encodeJsonResponse(rw http.ResponseWriter, response interface{}) {
	rw.Header().Set("Content-Type", "application/json")
	json.NewEncoder(rw).Encode(response)
}

func validateGetRequest(rw http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		rw.WriteHeader(http.StatusMethodNotAllowed)
	}
}
