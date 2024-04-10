package server

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/itimky/faraway-test/pkg/contract"
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

	repChallengeData, err := json.Marshal(contract.RepChallenge{
		Challenge: challenge,
	})
	if err != nil {
		return fmt.Errorf("marshal json: %w", err)
	}

	err = m.socket.Send(repChallengeData)
	if err != nil {
		return fmt.Errorf("send challenge: %w", err)
	}

	reqQuoteData, err := m.socket.Recv()
	if err != nil {
		return fmt.Errorf("recv solution: %w", err)
	}

	var reqQuote contract.ReqQuote

	err = json.Unmarshal(reqQuoteData, &reqQuote)
	if err != nil {
		return fmt.Errorf("unmarshal json: %w", err)
	}

	err = m.hashCash.ValidateSolution(challenge, int(reqQuote.Solution))
	if err != nil {
		repQuoteData, mErr := json.Marshal(contract.RepQuote{
			Error: err.Error(),
			Quote: "",
		})
		if mErr != nil {
			return fmt.Errorf("marshal json: %w", mErr)
		}

		mErr = m.socket.Send(repQuoteData)
		if mErr != nil {
			return fmt.Errorf("send quote: %w", mErr)
		}

		return fmt.Errorf("validate solution: %w", err)
	}

	return nil
}
