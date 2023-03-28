package books

import (
	"database/sql"
	"fmt"
	"github.com/krukkrz/pagination/pkg/books/model"
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

func (r Repository) FetchAll(limit, offset int) ([]model.Book, error) {
	log.Printf("fetching books with offset: %d and limit: %d", offset, limit)
	query := "SELECT * FROM books LIMIT $1 OFFSET $2;"

	rows, err := r.db.Query(query, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("error occured while running query: %s, error: %v", query, err)
	}
	defer rows.Close()

	var books []model.Book
	for rows.Next() {
		var book model.Book
		if err = rows.Scan(&book.Id, &book.Author, &book.Title, &book.CreatedAt); err != nil {
			return nil, fmt.Errorf("error while parsing rows: %v", err)
		}
		books = append(books, book)
	}

	return books, nil
}
