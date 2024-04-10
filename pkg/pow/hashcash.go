package pow

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"strings"
)

const (
	DefaultDifficulty int = 3
	randBytes         int = 16
)

var ErrInvalidSolution = errors.New("invalid solution")

type HashCash struct {
	randReader io.Reader
	difficulty int
}

func NewHashCash(
	randReader io.Reader,
	difficulty int,
) *HashCash {
	return &HashCash{
		randReader: randReader,
		difficulty: difficulty,
	}
}

func (hc *HashCash) GenerateChallenge() (string, error) {
	buf := make([]byte, randBytes)

	_, err := hc.randReader.Read(buf)
	if err != nil {
		return "", fmt.Errorf("read: %w", err)
	}

	return hex.EncodeToString(buf), nil
}

func (hc *HashCash) SolveChallenge(
	challenge string,
	difficulty int,
) int {
	prefix := strings.Repeat("0", difficulty)

	solution := 0

	for {
		testString := fmt.Sprintf("%s:%d", challenge, solution)
		hash := sha256.Sum256([]byte(testString))
		hexHash := hex.EncodeToString(hash[:])

		if strings.HasPrefix(hexHash, prefix) {
			return solution
		}

		solution++
	}
}

func (hc *HashCash) ValidateSolution(
	challenge string,
	solution int,
) error {
	hashInput := fmt.Sprintf("%s:%d", challenge, solution)
	hashBytes := sha256.Sum256([]byte(hashInput))
	hash := hex.EncodeToString(hashBytes[:])
	prefix := strings.Repeat("0", hc.difficulty)

	if !strings.HasPrefix(hash, prefix) {
		return ErrInvalidSolution
	}

	return nil
}
