package server_test

import (
	"context"
	"testing"

	"github.com/itimky/faraway-test/pkg/pow"
	"github.com/itimky/faraway-test/pkg/server"
	"github.com/itimky/faraway-test/test"
	netmocks "github.com/itimky/faraway-test/test/net"
	mocks "github.com/itimky/faraway-test/test/pkg/server"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type POWMiddlewareTest struct {
	suite.Suite

	connMock     *netmocks.MockConn
	hashCashMock *mocks.MockhashCash

	middleware *server.POWMiddleware
}

func (s *POWMiddlewareTest) SetupTest() {
	s.connMock = netmocks.NewMockConn(s.T())
	s.hashCashMock = mocks.NewMockhashCash(s.T())

	s.middleware = server.NewPOWMiddleware(s.hashCashMock)
}

func (s *POWMiddlewareTest) Test_Handle() {
	testCases := []struct {
		name                string
		expectedError       error
		initRecvErr         error
		genChallengeRes     string
		genChallengeErr     error
		sendChallengeParams []byte
		sendChallengeErr    error
		recvSolutionRes     []byte
		recvSolutionErr     error
		validateSolutionVal int
		validateSolutionErr error
	}{
		{
			name:            "err: generate challenge error",
			expectedError:   test.Err,
			genChallengeErr: test.Err,
		},
		{
			name:                "err: send challenge error",
			expectedError:       test.Err,
			genChallengeRes:     "challenge",
			sendChallengeParams: []byte("challenge\n"),
			sendChallengeErr:    test.Err,
		},
		{
			name:                "err: recv solution error",
			expectedError:       test.Err,
			genChallengeRes:     "challenge",
			sendChallengeParams: []byte("challenge\n"),
			recvSolutionErr:     test.Err,
		},
		{
			name:                "err: convert solution error",
			expectedError:       pow.ErrInvalidSolution,
			genChallengeRes:     "challenge",
			sendChallengeParams: []byte("challenge\n"),
			recvSolutionRes:     []byte("not a number\n"),
		},
		{
			name:                "err: validate solution error",
			expectedError:       test.Err,
			genChallengeRes:     "challenge",
			sendChallengeParams: []byte("challenge\n"),
			recvSolutionRes:     []byte("42\n"),
			validateSolutionVal: 42,
			validateSolutionErr: test.Err,
		},
		{
			name:                "ok",
			genChallengeRes:     "challenge",
			sendChallengeParams: []byte("challenge\n"),
			recvSolutionRes:     []byte("42\n"),
			validateSolutionVal: 42,
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			if tc.genChallengeRes != "" || tc.genChallengeErr != nil {
				s.hashCashMock.EXPECT().GenerateChallenge().Return(tc.genChallengeRes, tc.genChallengeErr).Once()
			}

			if tc.sendChallengeParams != nil {
				s.connMock.EXPECT().Write(tc.sendChallengeParams).
					Return(len(tc.sendChallengeParams), tc.sendChallengeErr).Once()
			}

			if tc.recvSolutionRes != nil || tc.recvSolutionErr != nil {
				s.connMock.EXPECT().Read(mock.Anything).Return(len(tc.recvSolutionRes), tc.recvSolutionErr).
					Run(func(buf []byte) {
						copy(buf, tc.recvSolutionRes)
					}).Once()
			}

			if tc.validateSolutionVal != 0 {
				s.hashCashMock.EXPECT().ValidateSolution(tc.genChallengeRes, tc.validateSolutionVal).
					Return(tc.validateSolutionErr).Once()
			}

			err := s.middleware.Handle(context.Background(), s.connMock)

			s.ErrorIs(err, tc.expectedError)
		})
	}
}

func TestPOWMiddlewareTestSuite(t *testing.T) {
	t.Parallel()

	suite.Run(t, new(POWMiddlewareTest))
}
