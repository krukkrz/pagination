package internal

import (
	"fmt"
	"github.com/krukkrz/pagination/pkg/api"
	books "github.com/krukkrz/pagination/pkg/books/model"
	cars "github.com/krukkrz/pagination/pkg/cars/model"
	"log"
	"testing"
)

type BookRepositorySuccessMock struct {
	expectedLimit  int
	expectedOffset int
	t              *testing.T
}

func (b BookRepositorySuccessMock) FetchAll(limit, offset int) ([]books.Book, error) {
	if b.expectedOffset != offset {
		b.t.Fatalf("incorrect offset expecting: %d, got: %d", b.expectedOffset, offset)
	}

	if b.expectedLimit != limit {
		b.t.Fatalf("incorrect limit expecting: %d, got: %d", b.expectedLimit, limit)
	}
	return Books, nil
}

func BookRepositoryMockReturnBooks(expectedLimit, expectedOffset int, t *testing.T) api.BookRepository {
	return &BookRepositorySuccessMock{
		expectedLimit:  expectedLimit,
		expectedOffset: expectedOffset,
		t:              t,
	}
}

type BookRepositoryErrorMock struct{}

func (b BookRepositoryErrorMock) FetchAll(limit, offset int) ([]books.Book, error) {
	return nil, fmt.Errorf("mocked error")
}

func BookServiceMockReturnError() api.BookRepository {
	return &BookRepositoryErrorMock{}
}

type CarRepositorySuccessMock struct {
	expectedLimit  int
	expectedCursor int
	t              *testing.T
}

func (b CarRepositorySuccessMock) FetchAll(cursor, limit int) ([]cars.Car, error) {
	if b.expectedCursor != cursor {
		log.Printf("incorrect cursor expecting: %d, got: %d", b.expectedCursor, cursor)
		b.t.Fail()
	}

	if b.expectedLimit != limit {
		log.Printf("incorrect limit expecting: %d, got: %d", b.expectedLimit, limit)
		b.t.Fail()
	}
	return Cars, nil
}

func CarRepositoryMockReturnCars(expectedLimit, expectedCursor int, t *testing.T) api.CarRepository {
	return &CarRepositorySuccessMock{
		expectedCursor: expectedCursor,
		expectedLimit:  expectedLimit,
		t:              t,
	}
}

type CarRepositoryErrorMock struct{}

func (b CarRepositoryErrorMock) FetchAll(cursor, limit int) ([]cars.Car, error) {
	return nil, fmt.Errorf("mocked error")
}

func CarRepositoryMockReturnError() api.CarRepository {
	return &CarRepositoryErrorMock{}
}
