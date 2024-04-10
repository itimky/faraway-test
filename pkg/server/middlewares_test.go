package server_test

import (
	"context"
	"testing"

	"github.com/itimky/faraway-test/pkg/server"
	"github.com/itimky/faraway-test/test"
	mocks "github.com/itimky/faraway-test/test/pkg/server"
	"github.com/stretchr/testify/suite"
)

type POWMiddlewareTest struct {
	suite.Suite

	socketMock   *mocks.Mocksocket
	hashCashMock *mocks.MockhashCash

	middleware *server.POWMiddleware
}

func (s *POWMiddlewareTest) SetupTest() {
	s.socketMock = mocks.NewMocksocket(s.T())
	s.hashCashMock = mocks.NewMockhashCash(s.T())

	s.middleware = server.NewPOWMiddleware(s.socketMock, s.hashCashMock)
}

func (s *POWMiddlewareTest) Test_Handle() {
	testCases := []struct {
		name                          string
		expectedError                 error
		initRecvErr                   error
		genChallengeRes               string
		genChallengeErr               error
		sendChallengeParams           []byte
		sendChallengeErr              error
		recvSolutionRes               []byte
		recvSolutionErr               error
		validateSolutionErrSendParams []byte
		validateSolutionErrSendErr    error
		validateSolutionVal           int
		validateSolutionErr           error
	}{
		{
			name:          "err: init recv error",
			expectedError: test.Err,
			initRecvErr:   test.Err,
		},
		{
			name:            "err: generate challenge error",
			expectedError:   test.Err,
			genChallengeErr: test.Err,
		},
		{
			name:                "err: send challenge error",
			expectedError:       test.Err,
			genChallengeRes:     "challenge",
			sendChallengeParams: []byte("challenge"),
			sendChallengeErr:    test.Err,
		},
		{
			name:                "err: recv solution error",
			expectedError:       test.Err,
			genChallengeRes:     "challenge",
			sendChallengeParams: []byte("challenge"),
			recvSolutionErr:     test.Err,
		},
		{
			name:                          "err: validate solution error: send err error",
			expectedError:                 test.Err,
			genChallengeRes:               "challenge",
			sendChallengeParams:           []byte("challenge"),
			recvSolutionRes:               []byte("42"),
			validateSolutionVal:           42,
			validateSolutionErr:           test.Err,
			validateSolutionErrSendParams: []byte(test.Err.Error()),
			validateSolutionErrSendErr:    test.Err,
		},
		{
			name:                          "err: validate solution error",
			expectedError:                 test.Err,
			genChallengeRes:               "challenge",
			sendChallengeParams:           []byte("challenge"),
			recvSolutionRes:               []byte("42"),
			validateSolutionVal:           42,
			validateSolutionErr:           test.Err,
			validateSolutionErrSendParams: []byte(test.Err.Error()),
		},
		{
			name:                "ok",
			genChallengeRes:     "challenge",
			sendChallengeParams: []byte("challenge"),
			recvSolutionRes:     []byte("42"),
			validateSolutionVal: 42,
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			s.socketMock.EXPECT().Recv().Return(nil, tc.initRecvErr).Once()

			if tc.genChallengeRes != "" || tc.genChallengeErr != nil {
				s.hashCashMock.EXPECT().GenerateChallenge().Return(tc.genChallengeRes, tc.genChallengeErr).Once()
			}

			if tc.sendChallengeParams != nil {
				s.socketMock.EXPECT().Send(tc.sendChallengeParams).Return(tc.sendChallengeErr).Once()
			}

			if tc.recvSolutionRes != nil || tc.recvSolutionErr != nil {
				s.socketMock.EXPECT().Recv().Return(tc.recvSolutionRes, tc.recvSolutionErr).Once()
			}

			if tc.validateSolutionVal != 0 {
				s.hashCashMock.EXPECT().ValidateSolution(tc.genChallengeRes, tc.validateSolutionVal).
					Return(tc.validateSolutionErr).Once()

				if tc.validateSolutionErr != nil {
					s.socketMock.EXPECT().Send([]byte(tc.validateSolutionErr.Error())).
						Return(tc.validateSolutionErrSendErr).Once()
				}
			}

			err := s.middleware.Handle(context.Background())

			s.ErrorIs(err, tc.expectedError)
		})
	}
}

func TestPOWMiddlewareTestSuite(t *testing.T) {
	t.Parallel()

	suite.Run(t, new(POWMiddlewareTest))
}
