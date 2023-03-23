package model

import "fmt"

type Book struct {
	Id        int
	Title     string
	Author    string
	CreatedAt string
}

func (b Book) String() string {
	return fmt.Sprintf("{Id: %d, Title: %s, Author: %s, CreatedAt: %s}", b.Id, b.Title, b.Author, b.CreatedAt)
}
