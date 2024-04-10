package client

import (
	"context"
	"fmt"
	"strconv"

	"github.com/itimky/faraway-test/pkg/pow"
)

type Client struct {
	socket   socket
	hashCash hashCash
}

func NewClient(
	socket socket,
	hashCash hashCash,
) *Client {
	return &Client{
		socket:   socket,
		hashCash: hashCash,
	}
}

func (c *Client) GetRandomQuote(
	_ context.Context,
) (string, error) {
	err := c.socket.Send([]byte(""))
	if err != nil {
		return "", fmt.Errorf("send: %w", err)
	}

	challenge, err := c.socket.Recv()
	if err != nil {
		return "", fmt.Errorf("recv: %w", err)
	}

	solution := c.hashCash.SolveChallenge(string(challenge), pow.DefaultDifficulty)

	err = c.socket.Send([]byte(strconv.Itoa(solution)))
	if err != nil {
		return "", fmt.Errorf("send: %w", err)
	}

	quote, err := c.socket.Recv()
	if err != nil {
		return "", fmt.Errorf("recv: %w", err)
	}

	return string(quote), nil
}
