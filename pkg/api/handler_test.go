package api_test

import (
	"encoding/json"
	"fmt"
	"github.com/krukkrz/pagination/pkg/api"
	"github.com/krukkrz/pagination/pkg/api/internal"
	booksModels "github.com/krukkrz/pagination/pkg/books/model"
	carsModels "github.com/krukkrz/pagination/pkg/cars/model"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func TestFetchAllBooks(t *testing.T) {
	testCases := []struct {
		skip            bool
		name            string
		limit           interface{}
		offset          interface{}
		method          string
		expectedStatus  int
		serviceError    bool
		expectedBooks   []booksModels.Book
		bookServiceMock api.BookRepository
		prevOffset      int
		nextOffset      int
	}{
		{
			name:            "handling only GET request",
			limit:           10,
			offset:          1,
			method:          "POST",
			bookServiceMock: internal.BookServiceMockReturnError(),
			expectedStatus:  http.StatusMethodNotAllowed,
		},
		{
			name:            "api requires offset and limit parameters in path",
			method:          "GET",
			bookServiceMock: internal.BookServiceMockReturnError(),
			expectedStatus:  http.StatusBadRequest,
		},
		{
			name:            "limit and offset must be a number",
			limit:           "l",
			offset:          "o",
			bookServiceMock: internal.BookServiceMockReturnError(),
			expectedStatus:  http.StatusBadRequest,
		},
		{
			name:            "api accepts offset and limit and pass it to bookRepository",
			limit:           10,
			offset:          1,
			bookServiceMock: internal.BookRepositoryMockReturnBooks(10, 1, t),
			expectedStatus:  http.StatusOK,
		},
		{
			name:            "returns internal server error if bookRepository returns error",
			limit:           10,
			offset:          1,
			serviceError:    true,
			bookServiceMock: internal.BookServiceMockReturnError(),
			expectedStatus:  http.StatusInternalServerError,
		},
		{
			name:            "returns books in response if all went good [0-10]",
			limit:           10,
			offset:          0,
			prevOffset:      0,
			nextOffset:      10,
			bookServiceMock: internal.BookRepositoryMockReturnBooks(10, 0, t),
			expectedBooks:   internal.Books,
			expectedStatus:  http.StatusOK,
		},
		{
			name:            "returns books in response if all went good [4-9]",
			limit:           5,
			offset:          4,
			prevOffset:      0,
			nextOffset:      9,
			bookServiceMock: internal.BookRepositoryMockReturnBooks(5, 4, t),
			expectedBooks:   internal.Books,
			expectedStatus:  http.StatusOK,
		},
		{
			name:            "returns books in response if all went good [9-14]",
			limit:           5,
			offset:          9,
			prevOffset:      4,
			nextOffset:      14,
			bookServiceMock: internal.BookRepositoryMockReturnBooks(5, 9, t),
			expectedBooks:   internal.Books,
			expectedStatus:  http.StatusOK,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.skip {
				t.SkipNow()
			}
			parameters := buildBooksParameters(tc.limit, tc.offset)

			req, err := http.NewRequest(tc.method, fmt.Sprintf("/books%s", parameters), nil)
			if err != nil {
				t.Fatal(err)
			}

			bs := tc.bookServiceMock
			cr := internal.CarRepositoryMockReturnCars(1, 1, t)
			srv := api.NewServer(bs, cr)

			rr := httptest.NewRecorder()
			handler := http.HandlerFunc(srv.FetchAllBooks)

			handler.ServeHTTP(rr, req)

			if status := rr.Code; status != tc.expectedStatus {
				t.Fatalf("api returned wrong status code: got %v want %v",
					status, tc.expectedStatus)
			}

			if tc.expectedBooks != nil {
				var actual api.PaginatedResponse[booksModels.Book]
				if err := json.NewDecoder(rr.Body).Decode(&actual); err != nil {
					t.Fatalf("unexpected error while parsing response body: %v", err)
				}

				if !reflect.DeepEqual(actual.Data, tc.expectedBooks) {
					t.Errorf("api returned unexpected body: got %v want %v", actual.Data, tc.expectedBooks)
				}

				expectedUrlFormat := "/books?limit=%v&offset=%d"
				expectedPrev := fmt.Sprintf(expectedUrlFormat, tc.limit, tc.prevOffset)
				if actual.Links.Prev != expectedPrev {
					t.Errorf("unexpected prev link value, got: %s, expected: %s", actual.Links.Prev, expectedPrev)
				}

				expectedNext := fmt.Sprintf(expectedUrlFormat, tc.limit, tc.nextOffset)
				if actual.Links.Next != expectedNext {
					t.Errorf("unexpected next link value, got: %s, expected: %s", actual.Links.Next, expectedNext)
				}
				expectedFirst := fmt.Sprintf(expectedUrlFormat, tc.limit, 0)
				if actual.Links.First != expectedFirst {
					t.Errorf("unexpected first link value, got: %s, expected: %s", actual.Links.First, expectedFirst)
				}
			}
		})
	}
}

func TestFetchAllCars(t *testing.T) {
	testCases := []struct {
		skip              bool
		name              string
		cursor            interface{}
		limit             interface{}
		method            string
		expectedStatus    int
		serviceError      bool
		expectedCars      []carsModels.Car
		carRepositoryMock api.CarRepository
		prevCursor        int
		nextCursor        int
	}{
		{
			name:              "handling only GET request",
			limit:             10,
			cursor:            1,
			method:            "POST",
			carRepositoryMock: internal.CarRepositoryMockReturnError(),
			expectedStatus:    http.StatusMethodNotAllowed,
		},
		{
			name:              "api requires offset and limit parameters in path",
			method:            "GET",
			carRepositoryMock: internal.CarRepositoryMockReturnError(),
			expectedStatus:    http.StatusBadRequest,
		},
		{
			name:              "limit and offset must be a number",
			limit:             "l",
			cursor:            "o",
			carRepositoryMock: internal.CarRepositoryMockReturnError(),
			expectedStatus:    http.StatusBadRequest,
		},
		{
			name:              "api accepts offset and limit and pass it to bookRepository",
			limit:             10,
			cursor:            1,
			carRepositoryMock: internal.CarRepositoryMockReturnCars(10, 1, t),
			expectedStatus:    http.StatusOK,
		},
		{
			name:              "returns internal server error if bookRepository returns error",
			limit:             10,
			cursor:            1,
			serviceError:      true,
			carRepositoryMock: internal.CarRepositoryMockReturnError(),
			expectedStatus:    http.StatusInternalServerError,
		},
		{
			name:              "returns books in response if all went good [0-10]",
			limit:             10,
			cursor:            1,
			prevCursor:        1,
			nextCursor:        11,
			carRepositoryMock: internal.CarRepositoryMockReturnCars(10, 1, t),
			expectedCars:      internal.Cars,
			expectedStatus:    http.StatusOK,
		},
		{
			name:              "returns books in response if all went good [4-9]",
			limit:             5,
			cursor:            4,
			prevCursor:        1,
			nextCursor:        9,
			carRepositoryMock: internal.CarRepositoryMockReturnCars(5, 4, t),
			expectedCars:      internal.Cars,
			expectedStatus:    http.StatusOK,
		},
		{
			name:              "returns books in response if all went good [9-14]",
			limit:             5,
			cursor:            9,
			prevCursor:        4,
			nextCursor:        14,
			carRepositoryMock: internal.CarRepositoryMockReturnCars(5, 9, t),
			expectedCars:      internal.Cars,
			expectedStatus:    http.StatusOK,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.skip {
				t.SkipNow()
			}
			parameters := buildCarsParameters(tc.cursor, tc.limit)

			req, err := http.NewRequest(tc.method, fmt.Sprintf("/carsModels%s", parameters), nil)
			if err != nil {
				t.Fatal(err)
			}

			br := internal.BookRepositoryMockReturnBooks(1, 1, t)
			cr := tc.carRepositoryMock
			srv := api.NewServer(br, cr)

			rr := httptest.NewRecorder()
			handler := http.HandlerFunc(srv.FetchAllCars)

			handler.ServeHTTP(rr, req)

			if status := rr.Code; status != tc.expectedStatus {
				t.Fatalf("api returned wrong status code: got %v want %v",
					status, tc.expectedStatus)
			}

			if tc.expectedCars != nil {
				var actual api.PaginatedResponse[carsModels.Car]
				if err := json.NewDecoder(rr.Body).Decode(&actual); err != nil {
					t.Fatalf("unexpected error while parsing response body: %v", err)
				}

				if !reflect.DeepEqual(actual.Data, tc.expectedCars) {
					t.Errorf("api returned unexpected body: got %v want %v", actual.Data, tc.expectedCars)
				}

				expectedUrlFormat := "/carsModels?cursor=%v&limit=%d"
				expectedPrev := fmt.Sprintf(expectedUrlFormat, tc.prevCursor, tc.limit)
				if actual.Links.Prev != expectedPrev {
					t.Errorf("unexpected prev link value, got: %s, expected: %s", actual.Links.Prev, expectedPrev)
				}

				expectedNext := fmt.Sprintf(expectedUrlFormat, tc.nextCursor, tc.limit)
				if actual.Links.Next != expectedNext {
					t.Errorf("unexpected next link value, got: %s, expected: %s", actual.Links.Next, expectedNext)
				}
				expectedFirst := fmt.Sprintf(expectedUrlFormat, 1, tc.limit)
				if actual.Links.First != expectedFirst {
					t.Errorf("unexpected first link value, got: %s, expected: %s", actual.Links.First, expectedFirst)
				}
			}
		})
	}
}

func buildBooksParameters(limit, offset interface{}) string {
	return fmt.Sprintf("?limit=%d&offset=%d", limit, offset)
}

func buildCarsParameters(cursor, limit interface{}) string {
	return fmt.Sprintf("?cursor=%d&limit=%d", cursor, limit)
}
