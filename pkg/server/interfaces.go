package server

type book interface {
	GetRandomQuote() (string, error)
}

type hashCash interface {
	GenerateChallenge() (string, error)
	ValidateSolution(challenge string, solution int) error
}
