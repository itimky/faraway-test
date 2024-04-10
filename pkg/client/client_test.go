package client_test

import (
	"context"
	"github.com/itimky/faraway-test/pkg/client"
	"github.com/itimky/faraway-test/pkg/pow"
	"github.com/itimky/faraway-test/test"
	mocks "github.com/itimky/faraway-test/test/pkg/client"
	"github.com/stretchr/testify/suite"
	"testing"
)

type ClientSuite struct {
	suite.Suite

	socketMock   *mocks.Mocksocket
	hashCashMock *mocks.MockhashCash

	client *client.Client
}

func (s *ClientSuite) SetupTest() {
	s.socketMock = mocks.NewMocksocket(s.T())
	s.hashCashMock = mocks.NewMockhashCash(s.T())

	s.client = client.NewClient(s.socketMock, s.hashCashMock)
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
			name:        "err: init send error",
			expectedErr: test.Err,
			initSendErr: test.Err,
		},
		{
			name:             "err: challenge recv error",
			expectedErr:      test.Err,
			challengeRecvErr: test.Err,
		},
		{
			name:                "err: solution send error",
			expectedErr:         test.Err,
			challengeRecvResult: []byte("challenge"),
			solveParams:         "challenge",
			solveResult:         42,
			solutionSendParams:  []byte("42"),
			solutionSendErr:     test.Err,
		},
		{
			name:                "err: quote recv error",
			expectedErr:         test.Err,
			challengeRecvResult: []byte("challenge"),
			solveParams:         "challenge",
			solveResult:         42,
			solutionSendParams:  []byte("42"),
			quoteRecvErr:        test.Err,
		},
		{
			name:                "ok",
			expectedResult:      "quote",
			challengeRecvResult: []byte("challenge"),
			solveParams:         "challenge",
			solveResult:         42,
			solutionSendParams:  []byte("42"),
			quoteRecvResult:     []byte("quote"),
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			s.socketMock.EXPECT().Send([]byte("")).Return(tc.initSendErr).Once()

			if tc.challengeRecvResult != nil || tc.challengeRecvErr != nil {
				s.socketMock.EXPECT().Recv().Return(tc.challengeRecvResult, tc.challengeRecvErr).Once()
			}

			if tc.solveParams != "" {
				s.hashCashMock.EXPECT().SolveChallenge(tc.solveParams, pow.DefaultDifficulty).Return(tc.solveResult)
			}

			if tc.solutionSendParams != nil {
				s.socketMock.EXPECT().Send(tc.solutionSendParams).Return(tc.solutionSendErr).Once()
			}

			if tc.quoteRecvResult != nil || tc.quoteRecvErr != nil {
				s.socketMock.EXPECT().Recv().Return(tc.quoteRecvResult, tc.quoteRecvErr).Once()
			}

			result, err := s.client.GetRandomQuote(context.Background())
			s.Equal(tc.expectedResult, result)
			s.ErrorIs(err, tc.expectedErr)
		})
	}
}

func TestClientSuite(t *testing.T) {
	t.Parallel()

	suite.Run(t, new(ClientSuite))
}
