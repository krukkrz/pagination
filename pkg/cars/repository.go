package cars

import (
	"database/sql"
	"fmt"
	"github.com/krukkrz/pagination/pkg/cars/model"
	"log"
)

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{
		db: db,
	}
}

func (r Repository) FetchAll(cursor, limit int) ([]model.Car, error) {
	log.Printf("fetching cars with cursor: %d and limit: %d", cursor, limit)
	query := "SELECT * FROM cars WHERE car_id >= $1 LIMIT $2;"

	rows, err := r.db.Query(query, cursor, limit)
	if err != nil {
		return nil, fmt.Errorf("error occured while running query: %s, error: %v", query, err)
	}
	defer rows.Close()

	var cars []model.Car
	for rows.Next() {
		var car model.Car
		if err = rows.Scan(&car.Id, &car.Brand, &car.Model, &car.CreatedAt); err != nil {
			return nil, fmt.Errorf("error while parsing rows: %v", err)
		}
		cars = append(cars, car)
	}

	return cars, nil
}
