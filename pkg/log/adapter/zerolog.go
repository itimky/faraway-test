package logadapter

import (
	"os"

	"github.com/itimky/faraway-test/pkg/log"
	"github.com/rs/zerolog"
)

type ZeroLog struct {
	logger zerolog.Logger
}

func NewZeroLog() ZeroLog {
	return ZeroLog{
		logger: zerolog.New(os.Stderr).With().Timestamp().Logger(),
	}
}

func (l ZeroLog) Info(msg string) {
	l.logger.Info().Msg(msg)
}

func (l ZeroLog) Infof(msg string, args ...interface{}) {
	l.logger.Info().Msgf(msg, args...)
}

func (l ZeroLog) Error(msg string) {
	l.logger.Error().Msg(msg)
}

func (l ZeroLog) Errorf(msg string, args ...interface{}) {
	l.logger.Error().Msgf(msg, args...)
}

func (l ZeroLog) WithField( //nolint: ireturn
	key string,
	value string,
) log.Logger {
	return ZeroLog{
		logger: l.logger.With().Str(key, value).Logger(),
	}
}
