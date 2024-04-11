package server_test

import (
	"context"
	"testing"

	"github.com/itimky/faraway-test/pkg/server"
	"github.com/itimky/faraway-test/test"
	netmocks "github.com/itimky/faraway-test/test/net"
	mocks "github.com/itimky/faraway-test/test/pkg/server"
	"github.com/stretchr/testify/suite"
)

type HandlerSuite struct {
	suite.Suite

	connMock *netmocks.MockConn
	bookMock *mocks.Mockbook

	handler *server.Handler
}

func (s *HandlerSuite) SetupTest() {
	s.connMock = netmocks.NewMockConn(s.T())
	s.bookMock = mocks.NewMockbook(s.T())
	s.handler = server.NewHandler(s.bookMock)
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
			sendParams:  []byte("quote\n"),
			sendError:   test.Err,
		},
		{
			name:        "ok",
			quoteResult: "quote",
			sendParams:  []byte("quote\n"),
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			s.bookMock.EXPECT().GetRandomQuote().Return(tc.quoteResult, tc.quoteError).Once()

			if tc.sendParams != nil {
				s.connMock.EXPECT().Write(tc.sendParams).Return(len(tc.sendParams), tc.sendError).Once()
			}

			err := s.handler.Handle(context.Background(), s.connMock)

			s.ErrorIs(err, tc.expectedErr)
		})
	}
}

func TestHandlerSuite(t *testing.T) {
	t.Parallel()

	suite.Run(t, new(HandlerSuite))
}
