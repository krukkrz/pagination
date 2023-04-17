package cars_test

import (
	"database/sql"
	"fmt"
	"github.com/krukkrz/pagination/pkg/cars"
	_ "github.com/proullon/ramsql/driver"
	"testing"
)

func TestFetchAll(t *testing.T) {
	db, err := sql.Open("ramsql", "Test cars")
	if err != nil {
		t.Fatalf("sql.Open : Error : %s\n", err)
	}
	defer db.Close()
	initDatabaseData(t, db)

	testCases := []struct {
		name                   string
		cursor                 int
		limit                  int
		expectedFirstElementId int
		expectedLastElementId  int
	}{
		{
			name:                   "returns cars from 2 to 6",
			cursor:                 2,
			limit:                  5,
			expectedFirstElementId: 2,
			expectedLastElementId:  6,
		},
		{
			name:                   "returns cars from 6 to 7",
			cursor:                 6,
			limit:                  2,
			expectedFirstElementId: 6,
			expectedLastElementId:  7,
		},
		{
			name:                   "returns only one car",
			cursor:                 6,
			limit:                  1,
			expectedFirstElementId: 6,
			expectedLastElementId:  6,
		},
		{
			name:   "returns no cars",
			cursor: 6,
			limit:  0,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			repo := cars.NewRepository(db)

			actual, err := repo.FetchAll(tc.cursor, tc.limit)
			if err != nil {
				t.Fatalf("unexpected error occured: %v", err)
			}

			if tc.limit == 0 {
				if len(actual) != 0 {
					t.Errorf("expecting no elements returned, got: %d", len(actual))
				}
			} else {
				firstElementId := actual[0].Id
				if firstElementId != tc.expectedFirstElementId {
					t.Errorf("expecting first returned element ID to be: %d, got: %d", tc.expectedFirstElementId, firstElementId)
				}

				lastElementId := actual[len(actual)-1].Id
				if lastElementId != tc.expectedLastElementId {
					t.Errorf("expecting last returned element ID to be: %d, got: %d", tc.expectedLastElementId, lastElementId)
				}
			}
		})
	}
}

func initDatabaseData(t *testing.T, db *sql.DB) {
	initTable := `CREATE TABLE if NOT EXISTS cars (car_id serial PRIMARY KEY, brand VARCHAR ( 100 ) NOT NULL, model VARCHAR ( 100 ) NOT NULL, created_at TIMESTAMP NOT NULL);`
	initialBatch := []string{
		initTable,
		insertCar(1),
		insertCar(2),
		insertCar(3),
		insertCar(4),
		insertCar(5),
		insertCar(6),
		insertCar(7),
		insertCar(8),
		insertCar(9),
		insertCar(10),
	}
	for _, b := range initialBatch {
		_, err := db.Exec(b)
		if err != nil {
			t.Fatalf("sql.Exec: Error: %s\n", err)
		}
	}
}

func insertCar(id int) string {
	return fmt.Sprintf("INSERT INTO cars (car_id, brand, model, created_at) VALUES (%d, 'Brand-%d', 'Model-%d', current_timestamp);", id, id, id)
}
