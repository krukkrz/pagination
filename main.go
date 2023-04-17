package main

import (
	"github.com/krukkrz/pagination/pkg/api"
	"github.com/krukkrz/pagination/pkg/books"
	"github.com/krukkrz/pagination/pkg/cars"
	"github.com/krukkrz/pagination/pkg/database"
	"log"
)

func main() {
	log.Println("Starting application...")
	db := database.Connect()
	bookRepository := books.NewRepository(db)
	carRepository := cars.NewRepository(db)
	server := api.NewServer(bookRepository, carRepository)
	server.Start(":8000")

	//todo dockerize everything
}
