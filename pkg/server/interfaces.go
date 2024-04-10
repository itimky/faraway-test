package server

import "context"

type book interface {
	GetRandomQuote() (string, error)
}

type hashCash interface {
	GenerateChallenge() (string, error)
	ValidateSolution(challenge string, solution int) error
}

type socket interface {
	Recv() ([]byte, error)
	Send(data []byte) error
}

type handler interface {
	Handle(ctx context.Context) error
}
