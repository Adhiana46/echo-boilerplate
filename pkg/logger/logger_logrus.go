package logger

import (
	"errors"
	"io"
	"os"
	"time"

	rotatelogs "github.com/lestrrat/go-file-rotatelogs"
	"github.com/sirupsen/logrus"
)

type logrusInterface interface {
	WithFields(fields logrus.Fields) *logrus.Entry

	Panic(...interface{})
	Panicf(string, ...interface{})
	Fatal(...interface{})
	Fatalf(string, ...interface{})
	Error(...interface{})
	Errorf(string, ...interface{})
	Warn(...interface{})
	Warnf(string, ...interface{})
	Info(...interface{})
	Infof(string, ...interface{})
	Trace(...interface{})
	Tracef(string, ...interface{})
	Debug(...interface{})
	Debugf(string, ...interface{})
}

type loggerLogrus struct {
	logger logrusInterface
}

func getLogrusLogLevel(level string) logrus.Level {
	switch level {
	case LevelPanic:
		return logrus.PanicLevel
	case LevelFatal:
		return logrus.FatalLevel
	case LevelError:
		return logrus.ErrorLevel
	case LevelWarn:
		return logrus.WarnLevel
	case LevelInfo:
		return logrus.InfoLevel
	case LevelTrace:
		return logrus.TraceLevel

	case LevelDebug:
		return logrus.DebugLevel
	}

	// default
	return logrus.DebugLevel
}

func NewLogrusLogger(cfg Config) (Logger, error) {
	if cfg.FileLocation == "" {
		return nil, errors.New("file location cannot be empty")
	}

	writer, err := rotatelogs.New(
		cfg.FileLocation+"/%Y%m%d%H%M.log",
		rotatelogs.WithLinkName(cfg.FileLocation+"/log.log"),
		rotatelogs.WithMaxAge(7*24*time.Hour),
		rotatelogs.WithRotationTime(24*time.Hour),
	)

	if err != nil {
		return nil, err
	}

	l := logrus.New()

	l.SetFormatter(&logrus.JSONFormatter{
		TimestampFormat: "2006-01-02 15:04:05",
		PrettyPrint:     false,
	})

	l.SetOutput(io.MultiWriter(writer, os.Stdout))

	logLevel := getLogrusLogLevel(cfg.Level)

	l.SetLevel(logLevel)

	return &loggerLogrus{
		logger: l,
	}, nil
}

func (l *loggerLogrus) Panic(args ...interface{}) {
	l.logger.Panic(args...)
}

func (l *loggerLogrus) Panicf(format string, args ...interface{}) {
	l.logger.Panicf(format, args...)
}

func (l *loggerLogrus) Fatal(args ...interface{}) {
	l.logger.Fatal(args...)
}

func (l *loggerLogrus) Fatalf(format string, args ...interface{}) {
	l.logger.Fatalf(format, args...)
}

func (l *loggerLogrus) Error(args ...interface{}) {
	l.logger.Error(args...)
}

func (l *loggerLogrus) Errorf(format string, args ...interface{}) {
	l.logger.Errorf(format, args...)
}

func (l *loggerLogrus) Warn(args ...interface{}) {
	l.logger.Warn(args...)
}

func (l *loggerLogrus) Warnf(format string, args ...interface{}) {
	l.logger.Warnf(format, args...)
}

func (l *loggerLogrus) Info(args ...interface{}) {
	l.logger.Info(args...)
}

func (l *loggerLogrus) Infof(format string, args ...interface{}) {
	l.logger.Infof(format, args...)
}

func (l *loggerLogrus) Debug(args ...interface{}) {
	l.logger.Debug(args...)
}

func (l *loggerLogrus) Debugf(format string, args ...interface{}) {
	l.logger.Debugf(format, args...)
}

func (l *loggerLogrus) Trace(args ...interface{}) {
	l.logger.Trace(args...)
}

func (l *loggerLogrus) Tracef(format string, args ...interface{}) {
	l.logger.Tracef(format, args...)
}

func (l *loggerLogrus) WithFields(fields Fields) Logger {
	logrusFields := logrus.Fields{}
	for key, val := range fields {
		logrusFields[key] = val
	}

	return &loggerLogrus{
		logger: l.logger.WithFields(logrusFields),
	}
}
