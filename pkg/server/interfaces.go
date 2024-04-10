package server

import "github.com/pebbe/zmq4"

type book interface {
	GetRandomQuote() (string, error)
}

type hashCash interface {
	GenerateChallenge() (string, error)
	ValidateSolution(challenge string, solution int) bool
}

type recver interface {
	Recv(flags zmq4.Flag) (string, error)
}

type sender interface {
	Send(data string, flags zmq4.Flag) (int, error)
}
