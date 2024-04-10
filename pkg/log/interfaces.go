package log

type Logger interface {
	Info(msg string)
	Infof(msg string, args ...interface{})
	Error(msg string)
	Errorf(msg string, args ...interface{})

	WithField(key string, value string) Logger
}
