package server

import (
	"context"
	"fmt"
	"net"
)

type Handler struct {
	book book
}

func NewHandler(
	book book,
) *Handler {
	return &Handler{
		book: book,
	}
}

func (s *Handler) Handle(
	_ context.Context,
	conn net.Conn,
) error {
	quote, err := s.book.GetRandomQuote()
	if err != nil {
		return fmt.Errorf("get random quote: %w", err)
	}

	_, err = conn.Write([]byte(fmt.Sprintf("%s\n", quote)))
	if err != nil {
		return fmt.Errorf("write quote: %w", err)
	}

	return nil
}
