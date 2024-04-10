package server

import (
	"context"
	"fmt"
)

type Handler struct {
	socket socket
	book   book
}

func NewHandler(
	socket socket,
	book book,
) *Handler {
	return &Handler{
		socket: socket,
		book:   book,
	}
}

func (s *Handler) Handle(
	_ context.Context,
) error {
	quote, err := s.book.GetRandomQuote()
	if err != nil {
		return fmt.Errorf("get random quote: %w", err)
	}

	err = s.socket.Send([]byte(quote))
	if err != nil {
		return fmt.Errorf("send quote: %w", err)
	}

	return nil
}
