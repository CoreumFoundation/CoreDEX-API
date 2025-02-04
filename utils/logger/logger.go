package logger

import (
	"os"
	"slices"
	"sync"
)

const ENV_CONFIG_NAME = "LOG_LEVEL"

// LogLevel represent level of logging.
type LogLevel string

const (
	// Debug has verbose message.
	Debug LogLevel = "debug"
	// Info is default log level.
	Info LogLevel = "info"
	// Warn is for logging warnings.
	Warn LogLevel = "warn"
	// Error is for logging errors.
	Error LogLevel = "error"
	// Fatal is for logging fatal messages. The system shutdown after logging the message.
	Fatal LogLevel = "fatal"
)

var allLogLevels = []LogLevel{
	Debug,
	Info,
	Warn,
	Error,
	Fatal,
}

// Fields to be passed when we want to call WithFields for structured logging.
type Fields map[string]interface{}

// A global variable so that log functions can be directly accessed.
var log Logger
var lock = sync.Mutex{}

func init() {
	lock.Lock()
	defer lock.Unlock()

	logLevelStr := os.Getenv(ENV_CONFIG_NAME)
	logLevel := Error
	if slices.Contains(allLogLevels, LogLevel(logLevelStr)) {
		logLevel = LogLevel(logLevelStr)
	}

	log = newZeroLogger(logLevel)
}

func GlobalSetLevel(level LogLevel) {
	lock.Lock()
	defer lock.Unlock()

	log = log.WithLevel(level)
}

func GlobalSetFields(fields Fields) {
	lock.Lock()
	defer lock.Unlock()

	log = log.WithFields(fields)
}

// Logger is our contract for the logger.
type Logger interface {
	Debugf(format string, args ...interface{})

	Infof(format string, args ...interface{})

	Warnf(format string, args ...interface{})

	Errorf(format string, args ...interface{})

	Fatalf(format string, args ...interface{})

	WithLevel(level LogLevel) Logger

	WithFields(fields Fields) Logger
}

// Debugf logs message with DEBUG log level.
func Debugf(format string, args ...interface{}) {
	log.Debugf(format, args...)
}

// Infof logs message with INFO log level.
func Infof(format string, args ...interface{}) {
	log.Infof(format, args...)
}

// Warnf logs message with WARN log level.
func Warnf(format string, args ...interface{}) {
	log.Warnf(format, args...)
}

// Errorf logs message with ERROR log level.
func Errorf(format string, args ...interface{}) {
	log.Errorf(format, args...)
}

// Fatalf logs message with FATAL log level.
func Fatalf(format string, args ...interface{}) {
	log.Fatalf(format, args...)
}

// WithFields builds nested logger with specified fields.
func WithFields(fields Fields) Logger {
	return log.WithFields(fields)
}
