package server

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/itimky/faraway-test/pkg/contract"
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

	repQuoteData, err := json.Marshal(contract.RepQuote{
		Error: "",
		Quote: quote,
	})
	if err != nil {
		return fmt.Errorf("marshal json: %w", err)
	}

	_, err = s.socket.Send(string(repQuoteData), 0)
	if err != nil {
		return fmt.Errorf("send quote: %w", err)
	}

	return nil
}
