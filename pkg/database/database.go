package database

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"log"
)

type config struct {
	host     string
	port     string
	user     string
	password string
	dbname   string
}

func (c config) String() string {
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", c.host, c.port, c.user, c.password, c.dbname)
}

func Connect() *sql.DB {
	//todo read those values from env
	cfg := config{
		host:     "localhost",
		port:     "5432",
		user:     "books",
		password: "books",
		dbname:   "booksdb",
	}
	db, err := sql.Open("postgres", cfg.String())
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Database connected!")
	return db
}
