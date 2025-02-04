package logger

import (
	"os"
	"time"

	"github.com/rs/zerolog"
)

type zeroLogLogger struct {
	logger zerolog.Logger
}

func zeroLogLevel(level LogLevel) zerolog.Level {
	switch level {
	case Debug:
		return zerolog.DebugLevel
	case Info:
		return zerolog.InfoLevel
	case Warn:
		return zerolog.WarnLevel
	case Error:
		return zerolog.ErrorLevel
	case Fatal:
		return zerolog.FatalLevel
	default:
		return zerolog.InfoLevel
	}
}

func newZeroLogger(logLevel LogLevel) Logger {
	severityMap := map[zerolog.Level]string{
		zerolog.DebugLevel: "DEBUG",
		zerolog.InfoLevel:  "INFO",
		zerolog.WarnLevel:  "WARNING",
		zerolog.ErrorLevel: "ERROR",
		zerolog.FatalLevel: "CRITICAL",
		zerolog.PanicLevel: "ALERT",
		zerolog.NoLevel:    "DEFAULT",
		zerolog.Disabled:   "DEFAULT",
		zerolog.TraceLevel: "DEBUG",
	}

	// This hook logs severity of log entry in Cloud Logging compatible format.
	hookFunc := zerolog.HookFunc(func(e *zerolog.Event, level zerolog.Level, msg string) {
		if severity, ok := severityMap[level]; ok {
			e.Str("severity", severity)
		}
	})

	zerolog.TimeFieldFormat = time.RFC3339
	logger := zerolog.New(os.Stdout).
		Hook(hookFunc).
		Level(zeroLogLevel(logLevel)).
		With().
		Timestamp().
		Logger()

	return &zeroLogLogger{
		logger: logger,
	}
}

// Debugf logs message with DEBUG log level.
func (l *zeroLogLogger) Debugf(format string, args ...interface{}) {
	l.logger.Debug().Msgf(format, args...)
}

// Infof logs message with INFO log level.
func (l *zeroLogLogger) Infof(format string, args ...interface{}) {
	l.logger.Info().Msgf(format, args...)
}

// Warnf logs message with WARN log level.
func (l *zeroLogLogger) Warnf(format string, args ...interface{}) {
	l.logger.Warn().Msgf(format, args...)
}

// Errorf logs message with ERROR log level.
func (l *zeroLogLogger) Errorf(format string, args ...interface{}) {
	l.logger.Error().Msgf(format, args...)
}

// Fatalf logs message with FATAL log level.
func (l *zeroLogLogger) Fatalf(format string, args ...interface{}) {
	l.logger.Fatal().Msgf(format, args...)
}

// WithFields builds nested logger with specified fields.
func (l *zeroLogLogger) WithFields(fields Fields) Logger {
	return &zeroLogLogger{
		logger: l.logger.With().Fields(map[string]interface{}(fields)).Logger(),
	}
}

func (l *zeroLogLogger) WithLevel(level LogLevel) Logger {
	return &zeroLogLogger{
		logger: l.logger.Level(zeroLogLevel(level)),
	}
}
