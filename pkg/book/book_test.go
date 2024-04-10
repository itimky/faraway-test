package book_test

import (
	"testing"

	"github.com/itimky/faraway-test/pkg/book"
	mocks "github.com/itimky/faraway-test/test/pkg/book"
	"github.com/stretchr/testify/suite"
)

type BookSuite struct {
	suite.Suite

	intnerMock *mocks.Mockintner
	quote      string
	book       *book.Book
}

func (s *BookSuite) SetupTest() {
	s.quote = "quote"
	s.intnerMock = mocks.NewMockintner(s.T())
	s.book = book.NewBook(
		s.intnerMock,
		[]string{s.quote},
	)
}

func (s *BookSuite) Test_GetRandomQuote_OK() {
	s.intnerMock.EXPECT().Intn(1).Return(0).Once()

	quote, err := s.book.GetRandomQuote()
	s.Equal(s.quote, quote)
	s.NoError(err)
}

func TestBookSuiteSuite(t *testing.T) {
	t.Parallel()

	suite.Run(t, new(BookSuite))
}
