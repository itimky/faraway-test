package client

type hashCash interface {
	SolveChallenge(challenge string, difficulty int) int
}

type socket interface {
	Recv() ([]byte, error)
	Send(data []byte) error
}
