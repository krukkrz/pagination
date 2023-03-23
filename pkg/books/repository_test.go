package books_test

import (
	"database/sql"
	"fmt"
	"github.com/krukkrz/pagination/pkg/books"
	_ "github.com/proullon/ramsql/driver"
	"testing"
)

func TestFetchAll(t *testing.T) {
	db, err := sql.Open("ramsql", "Test")
	if err != nil {
		t.Fatalf("sql.Open : Error : %s\n", err)
	}
	defer db.Close()
	initDatabaseData(t, db)

	testCases := []struct {
		name                   string
		limit                  int
		offset                 int
		expectedFirstElementId int
		expectedLastElementId  int
		skipCause              string
	}{
		{
			name:                   "should return books from 2 to 6",
			limit:                  5,
			offset:                 1,
			expectedFirstElementId: 2,
			expectedLastElementId:  6,
		},
		{
			name:                   "should return books from 4 to 7",
			limit:                  3,
			offset:                 3,
			expectedFirstElementId: 4,
			expectedLastElementId:  6,
		},
		{
			name:                   "should return books from 4 to 10",
			limit:                  3,
			offset:                 11,
			expectedFirstElementId: 4,
			expectedLastElementId:  10,
			skipCause:              "need to verify that real db will return expected result here",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.skipCause != "" {
				t.Skipf("skipping test because: %s", tc.skipCause)
			}
			repo := books.NewRepository(db)

			actual, err := repo.FetchAll(tc.limit, tc.offset)
			if err != nil {
				t.Fatalf("unexpected error occured: %v", err)
			}

			firstElementId := actual[0].Id
			if firstElementId != tc.expectedFirstElementId {
				t.Errorf("expecting first returned element ID to be: %d, got: %d", tc.expectedFirstElementId, firstElementId)
			}

			lastElementId := actual[len(actual)-1].Id
			if lastElementId != tc.expectedLastElementId {
				t.Errorf("expecting last returned element ID to be: %d, got: %d", tc.expectedLastElementId, lastElementId)
			}
		})
	}
}

func initDatabaseData(t *testing.T, db *sql.DB) {
	initTable := `CREATE TABLE if NOT EXISTS books (book_id serial PRIMARY KEY, title VARCHAR ( 100 ) NOT NULL, author VARCHAR ( 100 ) NOT NULL, created_at TIMESTAMP NOT NULL);`
	initialBatch := []string{
		initTable,
		insertBook(1),
		insertBook(2),
		insertBook(3),
		insertBook(4),
		insertBook(5),
		insertBook(6),
		insertBook(7),
		insertBook(8),
		insertBook(9),
		insertBook(10),
	}
	for _, b := range initialBatch {
		_, err := db.Exec(b)
		if err != nil {
			t.Fatalf("sql.Exec: Error: %s\n", err)
		}
	}
}

func insertBook(id int) string {
	return fmt.Sprintf("INSERT INTO books (book_id, title, author, created_at) VALUES (%d, 'Title-%d', 'Author-%d', current_timestamp);", id, id, id)
}
