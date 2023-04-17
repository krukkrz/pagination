package main

import (
	"github.com/krukkrz/pagination/pkg/api"
	"github.com/krukkrz/pagination/pkg/books"
	"github.com/krukkrz/pagination/pkg/database"
	"log"
)

func main() {
	log.Println("Starting application...")
	db := database.Connect()
	repository := books.NewRepository(db)
	server := api.NewServer(repository)
	server.Start(":8000")

	//todo dockerize everything
}
