package server

import (
	"context"
	"fmt"
	"strconv"
)

type POWMiddleware struct {
	socket   socket
	hashCash hashCash
}

func NewPOWMiddleware(
	socket socket,
	hashCash hashCash,
) *POWMiddleware {
	return &POWMiddleware{
		socket:   socket,
		hashCash: hashCash,
	}
}

func (m *POWMiddleware) Handle(
	_ context.Context,
) error {
	_, err := m.socket.Recv()
	if err != nil {
		return fmt.Errorf("recv: %w", err)
	}

	challenge, err := m.hashCash.GenerateChallenge()
	if err != nil {
		return fmt.Errorf("generate challenge: %w", err)
	}

	err = m.socket.Send([]byte(challenge))
	if err != nil {
		return fmt.Errorf("send challenge: %w", err)
	}

	solutionData, err := m.socket.Recv()
	if err != nil {
		return fmt.Errorf("recv solution: %w", err)
	}

	solution, err := strconv.Atoi(string(solutionData))
	if err != nil {
		return fmt.Errorf("convert solution: %w", err)
	}

	err = m.hashCash.ValidateSolution(challenge, solution)
	if err != nil {
		mErr := m.socket.Send([]byte(err.Error()))
		if mErr != nil {
			return fmt.Errorf("send quote: %w", mErr)
		}

		return fmt.Errorf("validate solution: %w", err)
	}

	return nil
}
