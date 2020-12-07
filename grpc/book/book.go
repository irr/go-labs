package book

import (
	log "github.com/sirupsen/logrus"

	"golang.org/x/net/context"
)

// Server interface for our service methods
type Server struct {
}

// GetBook logs Book from client and returns new Book
func (s *Server) GetBook(ctx context.Context, input *Book) (*Book, error) {

	log.WithFields(log.Fields{
		"Name": input.Name,
		"Isbn": input.Isbn,
	}).Info("Book data received from client")

	return &Book{Name: "The Great Gatsby", Isbn: 90393}, nil
}
