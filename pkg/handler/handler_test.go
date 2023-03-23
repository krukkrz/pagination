package handler_test

import (
	"encoding/json"
	"fmt"
	"github.com/krukkrz/pagination/pkg/handler"
	"github.com/krukkrz/pagination/pkg/handler/internal"
	"github.com/krukkrz/pagination/pkg/model"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func TestFetchAllBooksLimitAndOffset(t *testing.T) {
	testCases := []struct {
		name            string
		limit           interface{}
		offset          interface{}
		method          string
		expectedStatus  int
		serviceError    bool
		expectedBooks   []model.Book
		bookServiceMock handler.Service
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
			name:            "handler requires offset and limit parameters in path",
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
			name:            "handler accepts offset and limit and pass it to service",
			limit:           10,
			offset:          1,
			bookServiceMock: internal.BookServiceMockReturnBooks(10, 1, t),
			expectedStatus:  http.StatusOK,
		},
		{
			name:            "returns internal server error if service returns error",
			limit:           10,
			offset:          1,
			serviceError:    true,
			bookServiceMock: internal.BookServiceMockReturnError(),
			expectedStatus:  http.StatusInternalServerError,
		},

		{
			name:            "returns books in response if all went good",
			limit:           10,
			offset:          1,
			bookServiceMock: internal.BookServiceMockReturnBooks(10, 1, t),
			expectedBooks:   internal.Books,
			expectedStatus:  http.StatusOK,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			parameters := buildParameters(tc.limit, tc.offset)

			req, err := http.NewRequest(tc.method, fmt.Sprintf("/books%s", parameters), nil)
			if err != nil {
				t.Fatal(err)
			}

			bs := tc.bookServiceMock
			srv := handler.NewServer(bs)

			rr := httptest.NewRecorder()
			handler := http.HandlerFunc(srv.FetchAllBooksWithLimitAndOffset)

			handler.ServeHTTP(rr, req)

			if status := rr.Code; status != tc.expectedStatus {
				t.Fatalf("handler returned wrong status code: got %v want %v",
					status, tc.expectedStatus)
			}

			if tc.expectedBooks != nil {

				var actual []model.Book
				if err := json.NewDecoder(rr.Body).Decode(&actual); err != nil {
					t.Fatalf("unexpected error while parsing response body: %v", err)
				}

				if !reflect.DeepEqual(actual, tc.expectedBooks) {
					t.Errorf("handler returned unexpected body: got %v want %v", actual, tc.expectedBooks)
				}
			}
		})
	}
}

func buildParameters(limit interface{}, offset interface{}) string {
	if offset != 0 && limit != 0 {
		return fmt.Sprintf("?limit=%d&offset=%d", limit, offset)
	}
	return ""
}
