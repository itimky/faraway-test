package server_test

import (
	"context"
	"testing"

	"github.com/itimky/faraway-test/pkg/server"
	"github.com/itimky/faraway-test/test"
	mocks "github.com/itimky/faraway-test/test/pkg/server"
	"github.com/stretchr/testify/suite"
)

type HandlerSuite struct {
	suite.Suite

	socketMock *mocks.Mocksocket
	bookMock   *mocks.Mockbook

	handler *server.Handler
}

func (s *HandlerSuite) SetupTest() {
	s.socketMock = mocks.NewMocksocket(s.T())
	s.bookMock = mocks.NewMockbook(s.T())
	s.handler = server.NewHandler(s.socketMock, s.bookMock)
}

func (s *HandlerSuite) Test_Handle() {
	testCases := []struct {
		name        string
		expectedErr error
		quoteResult string
		quoteError  error
		sendParams  []byte
		sendError   error
	}{
		{
			name:        "err: get random quote error",
			expectedErr: test.Err,
			quoteError:  test.Err,
		},
		{
			name:        "err: send error",
			expectedErr: test.Err,
			quoteResult: "quote",
			sendParams:  []byte("quote"),
			sendError:   test.Err,
		},
		{
			name:        "ok",
			quoteResult: "quote",
			sendParams:  []byte("quote"),
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			s.bookMock.EXPECT().GetRandomQuote().Return(tc.quoteResult, tc.quoteError).Once()

			if tc.sendParams != nil {
				s.socketMock.EXPECT().Send(tc.sendParams).Return(tc.sendError).Once()
			}

			err := s.handler.Handle(context.Background())

			s.ErrorIs(err, tc.expectedErr)
		})
	}
}

func TestHandlerSuite(t *testing.T) {
	t.Parallel()

	suite.Run(t, new(HandlerSuite))
}
