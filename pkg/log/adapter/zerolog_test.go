package logadapter_test

import (
	"testing"

	logadapter "github.com/itimky/faraway-test/pkg/log/adapter"
	"github.com/stretchr/testify/suite"
)

type ZeroLogSuite struct {
	suite.Suite

	zeroLog logadapter.ZeroLog
}

func (s *ZeroLogSuite) SetupTest() {
	s.zeroLog = logadapter.NewZeroLog()
}

func (s *ZeroLogSuite) Test_Info() {
	s.zeroLog.Info("smoke test")
}

func (s *ZeroLogSuite) Test_Infof() {
	s.zeroLog.Infof("smoke test")
}

func (s *ZeroLogSuite) Test_Error() {
	s.zeroLog.Error("smoke test")
}

func (s *ZeroLogSuite) Test_Errorf() {
	s.zeroLog.Errorf("smoke test")
}

func (s *ZeroLogSuite) Test_Fatalf() {
	// Fatal() calls os.Exit(1) and is not testable
}

func (s *ZeroLogSuite) Test_WithField() {
	_ = s.zeroLog.WithField("key", "value")
}

func TestZeroLogAdapterSuite(t *testing.T) {
	t.Parallel()

	suite.Run(t, new(ZeroLogSuite))
}
