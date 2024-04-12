package client_test

import (
	"context"
	"testing"

	"github.com/itimky/faraway-test/pkg/client"
	"github.com/itimky/faraway-test/pkg/pow"
	"github.com/itimky/faraway-test/test"
	netmocks "github.com/itimky/faraway-test/test/net"
	mocks "github.com/itimky/faraway-test/test/pkg/client"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type ClientSuite struct {
	suite.Suite

	connMock     *netmocks.MockConn
	hashCashMock *mocks.MockhashCash

	client *client.Client
}

func (s *ClientSuite) SetupTest() {
	s.connMock = netmocks.NewMockConn(s.T())
	s.hashCashMock = mocks.NewMockhashCash(s.T())

	s.client = client.NewClient(s.hashCashMock)
}

func (s *ClientSuite) Test_GetRandomQuote() {
	testCases := []struct {
		name                string
		expectedResult      string
		expectedErr         error
		initSendErr         error
		challengeRecvResult []byte
		challengeRecvErr    error
		solveParams         string
		solveResult         int
		solutionSendParams  []byte
		solutionSendErr     error
		quoteRecvResult     []byte
		quoteRecvErr        error
	}{
		{
			name:             "err: challenge recv error",
			expectedErr:      test.Err,
			challengeRecvErr: test.Err,
		},
		{
			name:                "err: solution send error",
			expectedErr:         test.Err,
			challengeRecvResult: []byte("challenge\x03"),
			solveParams:         "challenge",
			solveResult:         42,
			solutionSendParams:  []byte("42\x03"),
			solutionSendErr:     test.Err,
		},
		{
			name:                "err: quote recv error",
			expectedErr:         test.Err,
			challengeRecvResult: []byte("challenge\x03"),
			solveParams:         "challenge",
			solveResult:         42,
			solutionSendParams:  []byte("42\x03"),
			quoteRecvErr:        test.Err,
		},
		{
			name:                "ok",
			expectedResult:      "quote",
			challengeRecvResult: []byte("challenge\x03"),
			solveParams:         "challenge",
			solveResult:         42,
			solutionSendParams:  []byte("42\x03"),
			quoteRecvResult:     []byte("quote\x03"),
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			s.connMock.EXPECT().Read(mock.Anything).Return(len(tc.challengeRecvResult), tc.challengeRecvErr).
				Run(func(b []byte) {
					copy(b, tc.challengeRecvResult)
				}).Once()

			if tc.solveParams != "" {
				s.hashCashMock.EXPECT().SolveChallenge(tc.solveParams, pow.DefaultDifficulty).Return(tc.solveResult)
			}

			if tc.solutionSendParams != nil {
				s.connMock.EXPECT().Write(tc.solutionSendParams).Return(len(tc.solutionSendParams), tc.solutionSendErr).
					Once()
			}

			if tc.quoteRecvResult != nil || tc.quoteRecvErr != nil {
				s.connMock.EXPECT().Read(mock.Anything).Return(len(tc.quoteRecvResult), tc.quoteRecvErr).
					Run(func(b []byte) {
						copy(b, tc.quoteRecvResult)
					}).Once()
			}

			result, err := s.client.GetRandomQuote(context.Background(), s.connMock)
			s.Equal(tc.expectedResult, result)
			s.ErrorIs(err, tc.expectedErr)
		})
	}
}

func TestClientSuite(t *testing.T) {
	t.Parallel()

	suite.Run(t, new(ClientSuite))
}
