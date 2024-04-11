package client

import (
	"bufio"
	"context"
	"fmt"
	"net"
	"strings"

	"github.com/itimky/faraway-test/pkg/pow"
)

type Client struct {
	hashCash hashCash
}

func NewClient(
	hashCash hashCash,
) *Client {
	return &Client{
		hashCash: hashCash,
	}
}

func (c *Client) GetRandomQuote(
	_ context.Context,
	conn net.Conn,
) (string, error) {
	reader := bufio.NewReader(conn)
	writer := bufio.NewWriter(conn)

	// Receive the PoW challenge from the server
	challenge, err := reader.ReadString('\n')
	if err != nil {
		return "", fmt.Errorf("read string: %w", err)
	}

	challenge = strings.TrimSpace(challenge)

	solution := c.hashCash.SolveChallenge(challenge, pow.DefaultDifficulty)

	_, err = fmt.Fprintln(writer, solution)
	if err != nil {
		return "", fmt.Errorf("fprintln: %w", err)
	}

	err = writer.Flush()
	if err != nil {
		return "", fmt.Errorf("flush: %w", err)
	}

	quote, err := reader.ReadString('\n')
	if err != nil {
		return "", fmt.Errorf("recv: %w", err)
	}

	return strings.TrimSpace(quote), nil
}
