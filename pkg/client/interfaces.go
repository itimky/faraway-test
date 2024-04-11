package client

type hashCash interface {
	SolveChallenge(challenge string, difficulty int) int
}
