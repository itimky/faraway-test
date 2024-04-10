package server

import (
	"context"
	"fmt"
)

type Handler struct {
	sender sender
	book   book
}

func NewHandler(
	sender sender,
	book book,
) *Handler {
	return &Handler{
		sender: sender,
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

	_, err = s.sender.Send(quote, 0)
	if err != nil {
		return fmt.Errorf("send quote: %w", err)
	}

	return nil
}
