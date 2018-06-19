package log

// Level defines all possible log levels
type Level uint8

// Logger defines the logging interface.
type Logger interface {
	SetLevel(level Level)
	Debugf(format string, args ...interface{})
	Infof(format string, args ...interface{})
	Warnf(format string, args ...interface{})
	Errorf(format string, args ...interface{})
	Error(err error)
	IsEnabled(level Level) bool
	Write(p []byte) (int, error)
}

const (
	// DEBUG for verbose logs
	DEBUG Level = iota + 1
	// INFO for WARN+ERROR+INFO logs
	INFO
	// WARN for WARN+ERROR logs
	WARN
	// ERROR for ERROR only logs
	ERROR
	// NONE is used to disable logs
	NONE
)