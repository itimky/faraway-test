package pow_test

import (
	"testing"

	"github.com/itimky/faraway-test/pkg/pow"
	"github.com/itimky/faraway-test/test"
	mocks "github.com/itimky/faraway-test/test/io"
	"github.com/stretchr/testify/suite"
)

type HashCashSuite struct {
	suite.Suite

	randReaderMock *mocks.MockReader
	difficulty     int

	hashCash *pow.HashCash
}

func (s *HashCashSuite) SetupTest() {
	s.randReaderMock = mocks.NewMockReader(s.T())
	s.difficulty = 1
	s.hashCash = pow.NewHashCash(
		s.randReaderMock,
		s.difficulty,
	)
}

func (s *HashCashSuite) Test_GenerateChallenge() {
	testCases := []struct {
		name             string
		expectedResult   string
		expectedErr      error
		randReaderResult []byte
		randReaderErr    error
	}{
		{
			name:             "err: rand reader error",
			expectedResult:   "",
			expectedErr:      test.Err,
			randReaderResult: nil,
			randReaderErr:    test.Err,
		},
		{
			name:             "ok",
			expectedResult:   "000102030405060708090a0b0c0d0e0f",
			expectedErr:      nil,
			randReaderResult: []byte{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15},
			randReaderErr:    nil,
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			s.randReaderMock.EXPECT().Read(make([]byte, 16)).Return(16, tc.randReaderErr).Run(func(b []byte) {
				if tc.randReaderResult != nil {
					copy(b, tc.randReaderResult)
				}
			}).Once()

			result, err := s.hashCash.GenerateChallenge()
			s.Equal(tc.expectedResult, result)
			s.ErrorIs(err, tc.expectedErr)
		})
	}
}

func (s *HashCashSuite) Test_SolveChallenge() {
	testCases := []struct {
		name           string
		expectedResult int
		challenge      string
	}{
		{
			name:           "ok",
			expectedResult: 6,
			challenge:      "000102030405060708090a0b0c0d0e0f",
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			result := s.hashCash.SolveChallenge(tc.challenge, s.difficulty)
			s.Equal(tc.expectedResult, result)
		})
	}
}

func (s *HashCashSuite) Test_ValidateSolution() {
	testCases := []struct {
		name        string
		expectedErr error
		challenge   string
		solution    int
	}{
		{
			name:        "err: invalid solution",
			expectedErr: pow.ErrInvalidSolution,
			challenge:   "000102030405060708090a0b0c0d0e0f",
			solution:    5,
		},
		{
			name:        "ok",
			expectedErr: nil,
			challenge:   "000102030405060708090a0b0c0d0e0f",
			solution:    6,
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			err := s.hashCash.ValidateSolution(tc.challenge, tc.solution)
			s.ErrorIs(err, tc.expectedErr)
		})
	}
}

func TestHashCashSuite(t *testing.T) {
	t.Parallel()

	suite.Run(t, new(HashCashSuite))
}
