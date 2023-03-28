package api_test

import (
	"encoding/json"
	"fmt"
	"github.com/krukkrz/pagination/pkg/books/api"
	internal2 "github.com/krukkrz/pagination/pkg/books/api/internal"
	"github.com/krukkrz/pagination/pkg/books/model"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func TestFetchAllBooksLimitAndOffset(t *testing.T) {
	testCases := []struct {
		skip            bool
		name            string
		limit           interface{}
		offset          interface{}
		method          string
		expectedStatus  int
		serviceError    bool
		expectedBooks   []model.Book
		bookServiceMock api.BookRepository
		prevOffset      int
		nextOffset      int
	}{
		{
			name:            "handling only GET request",
			limit:           10,
			offset:          1,
			method:          "POST",
			bookServiceMock: internal2.BookServiceMockReturnError(),
			expectedStatus:  http.StatusMethodNotAllowed,
		},
		{
			name:            "api requires offset and limit parameters in path",
			method:          "GET",
			bookServiceMock: internal2.BookServiceMockReturnError(),
			expectedStatus:  http.StatusBadRequest,
		},
		{
			name:            "limit and offset must be a number",
			limit:           "l",
			offset:          "o",
			bookServiceMock: internal2.BookServiceMockReturnError(),
			expectedStatus:  http.StatusBadRequest,
		},
		{
			name:            "api accepts offset and limit and pass it to repository",
			limit:           10,
			offset:          1,
			bookServiceMock: internal2.BookServiceMockReturnBooks(10, 1, t),
			expectedStatus:  http.StatusOK,
		},
		{
			name:            "returns internal server error if repository returns error",
			limit:           10,
			offset:          1,
			serviceError:    true,
			bookServiceMock: internal2.BookServiceMockReturnError(),
			expectedStatus:  http.StatusInternalServerError,
		},
		{
			name:            "returns books in response if all went good [0-10]",
			limit:           10,
			offset:          0,
			prevOffset:      0,
			nextOffset:      10,
			bookServiceMock: internal2.BookServiceMockReturnBooks(10, 0, t),
			expectedBooks:   internal2.Books,
			expectedStatus:  http.StatusOK,
		},
		{
			name:            "returns books in response if all went good [4-9]",
			limit:           5,
			offset:          4,
			prevOffset:      0,
			nextOffset:      9,
			bookServiceMock: internal2.BookServiceMockReturnBooks(5, 4, t),
			expectedBooks:   internal2.Books,
			expectedStatus:  http.StatusOK,
		},
		{
			name:            "returns books in response if all went good [9-14]",
			limit:           5,
			offset:          9,
			prevOffset:      4,
			nextOffset:      14,
			bookServiceMock: internal2.BookServiceMockReturnBooks(5, 9, t),
			expectedBooks:   internal2.Books,
			expectedStatus:  http.StatusOK,
		},
		{
			name: "if no parameters in the request, return first page with default 10 elements", //todo implement this test
			skip: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.skip {
				t.SkipNow()
			}
			parameters := buildParameters(tc.limit, tc.offset)

			req, err := http.NewRequest(tc.method, fmt.Sprintf("/books%s", parameters), nil)
			if err != nil {
				t.Fatal(err)
			}

			bs := tc.bookServiceMock
			srv := api.NewServer(bs)

			rr := httptest.NewRecorder()
			handler := http.HandlerFunc(srv.FetchAllBooks)

			handler.ServeHTTP(rr, req)

			if status := rr.Code; status != tc.expectedStatus {
				t.Fatalf("api returned wrong status code: got %v want %v",
					status, tc.expectedStatus)
			}

			if tc.expectedBooks != nil {
				var actual api.PaginatedResponse
				if err := json.NewDecoder(rr.Body).Decode(&actual); err != nil {
					t.Fatalf("unexpected error while parsing response body: %v", err)
				}

				if !reflect.DeepEqual(actual.Data, tc.expectedBooks) {
					t.Errorf("api returned unexpected body: got %v want %v", actual.Data, tc.expectedBooks)
				}

				expectedPrev := fmt.Sprintf("/books?limit=%v&offset=%d", tc.limit, tc.prevOffset)
				if actual.Links.Prev != expectedPrev {
					t.Errorf("unexpected prev link value, got: %s, expected: %s", actual.Links.Prev, expectedPrev)
				}

				expectedNext := fmt.Sprintf("/books?limit=%v&offset=%d", tc.limit, tc.nextOffset)
				if actual.Links.Next != expectedNext {
					t.Errorf("unexpected next link value, got: %s, expected: %s", actual.Links.Next, expectedNext)
				}
				expectedFirst := fmt.Sprintf("/books?limit=%v&offset=%d", tc.limit, 0)
				if actual.Links.First != expectedFirst {
					t.Errorf("unexpected first link value, got: %s, expected: %s", actual.Links.First, expectedFirst)
				}
			}
		})
	}
}

func buildParameters(limit interface{}, offset interface{}) string {
	return fmt.Sprintf("?limit=%d&offset=%d", limit, offset)
}
