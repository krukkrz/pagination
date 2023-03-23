package internal

import (
	"fmt"
	"github.com/krukkrz/pagination/pkg/handler"
	"github.com/krukkrz/pagination/pkg/model"
	"testing"
)

type BookServiceSuccessMock struct {
	expectedLimit  int
	expectedOffset int
	t              *testing.T
}

func (b BookServiceSuccessMock) FetchAllLimitAndOffset(limit, offset int) ([]model.Book, error) {
	if b.expectedOffset != offset {
		b.t.Fatalf("incorrect offset expecting: %d, got: %d", b.expectedOffset, offset)
	}

	if b.expectedLimit != limit {
		b.t.Fatalf("incorrect limit expecting: %d, got: %d", b.expectedLimit, limit)
	}
	return Books, nil
}

func BookServiceMockReturnBooks(expectedLimit, expectedOffset int, t *testing.T) handler.Service {
	return &BookServiceSuccessMock{
		expectedLimit:  expectedLimit,
		expectedOffset: expectedOffset,
		t:              t,
	}
}

type BookServiceErrorMock struct {
}

func (b BookServiceErrorMock) FetchAllLimitAndOffset(limit, offset int) ([]model.Book, error) {
	return nil, fmt.Errorf("mocked error")
}

func BookServiceMockReturnError() handler.Service {
	return &BookServiceErrorMock{}
}
