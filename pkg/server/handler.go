package server

import (
	"bufio"
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
	writer := bufio.NewWriter(conn)

	quote, err := s.book.GetRandomQuote()
	if err != nil {
		return fmt.Errorf("get random quote: %w", err)
	}

	_, err = fmt.Fprintf(writer, "%s\x03", quote)
	if err != nil {
		return fmt.Errorf("fprintf: %w", err)
	}

	err = writer.Flush()
	if err != nil {
		return fmt.Errorf("flush: %w", err)
	}

	return nil
}
