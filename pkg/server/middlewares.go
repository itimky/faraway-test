package server

import (
	"bufio"
	"context"
	"fmt"
	"net"
	"strconv"
	"strings"

	"github.com/itimky/faraway-test/pkg/pow"
)

type POWMiddleware struct {
	hashCash hashCash
}

func NewPOWMiddleware(
	hashCash hashCash,
) *POWMiddleware {
	return &POWMiddleware{
		hashCash: hashCash,
	}
}

func (m *POWMiddleware) Handle(
	_ context.Context,
	conn net.Conn,
) error {
	reader := bufio.NewReader(conn)
	writer := bufio.NewWriter(conn)

	challenge, err := m.hashCash.GenerateChallenge()
	if err != nil {
		return fmt.Errorf("generate challenge: %w", err)
	}

	_, err = fmt.Fprintf(writer, "%s\x03", challenge)
	if err != nil {
		return fmt.Errorf("fprintf: %w", err)
	}

	err = writer.Flush()
	if err != nil {
		return fmt.Errorf("flush: %w", err)
	}

	solutionData, err := reader.ReadString('\x03')
	if err != nil {
		return fmt.Errorf("read string: %w", err)
	}

	solutionData = strings.TrimRight(solutionData, "\x03")

	solution, err := strconv.Atoi(solutionData)
	if err != nil {
		return pow.ErrInvalidSolution
	}

	err = m.hashCash.ValidateSolution(challenge, solution)
	if err != nil {
		return fmt.Errorf("validate solution: %w", err)
	}

	return nil
}
