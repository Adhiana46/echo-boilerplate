package logger

import (
	"errors"
	"io"
	"os"
	"time"

	rotatelogs "github.com/lestrrat/go-file-rotatelogs"
	"github.com/sirupsen/logrus"
)

type loggerLogrus struct {
	logger       logrus.Logger
	logrusFields []logrus.Fields
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

	l := logrus.Logger{}

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
	if len(l.logrusFields) > 0 {
		finalLogrusFields := logrus.Fields{}
		for _, logrusFields := range l.logrusFields {
			for key, val := range logrusFields {
				finalLogrusFields[key] = val
			}
		}

		l.logger.WithFields(finalLogrusFields).Panic(args...)

		// reset fields
		l.logrusFields = []logrus.Fields{}
	} else {
		l.logger.Panic(args...)
	}

}

func (l *loggerLogrus) Panicf(format string, args ...interface{}) {
	if len(l.logrusFields) > 0 {
		finalLogrusFields := logrus.Fields{}
		for _, logrusFields := range l.logrusFields {
			for key, val := range logrusFields {
				finalLogrusFields[key] = val
			}
		}

		l.logger.WithFields(finalLogrusFields).Panicf(format, args...)

		// reset fields
		l.logrusFields = []logrus.Fields{}
	} else {
		l.logger.Panicf(format, args...)
	}

}

func (l *loggerLogrus) Fatal(args ...interface{}) {
	if len(l.logrusFields) > 0 {
		finalLogrusFields := logrus.Fields{}
		for _, logrusFields := range l.logrusFields {
			for key, val := range logrusFields {
				finalLogrusFields[key] = val
			}
		}

		l.logger.WithFields(finalLogrusFields).Fatal(args...)

		// reset fields
		l.logrusFields = []logrus.Fields{}
	} else {
		l.logger.Fatal(args...)
	}

}

func (l *loggerLogrus) Fatalf(format string, args ...interface{}) {
	if len(l.logrusFields) > 0 {
		finalLogrusFields := logrus.Fields{}
		for _, logrusFields := range l.logrusFields {
			for key, val := range logrusFields {
				finalLogrusFields[key] = val
			}
		}

		l.logger.WithFields(finalLogrusFields).Fatalf(format, args...)

		// reset fields
		l.logrusFields = []logrus.Fields{}
	} else {
		l.logger.Fatalf(format, args...)
	}

}

func (l *loggerLogrus) Error(args ...interface{}) {
	if len(l.logrusFields) > 0 {
		finalLogrusFields := logrus.Fields{}
		for _, logrusFields := range l.logrusFields {
			for key, val := range logrusFields {
				finalLogrusFields[key] = val
			}
		}

		l.logger.WithFields(finalLogrusFields).Error(args...)

		// reset fields
		l.logrusFields = []logrus.Fields{}
	} else {
		l.logger.Error(args...)
	}

}

func (l *loggerLogrus) Errorf(format string, args ...interface{}) {
	if len(l.logrusFields) > 0 {
		finalLogrusFields := logrus.Fields{}
		for _, logrusFields := range l.logrusFields {
			for key, val := range logrusFields {
				finalLogrusFields[key] = val
			}
		}

		l.logger.WithFields(finalLogrusFields).Errorf(format, args...)

		// reset fields
		l.logrusFields = []logrus.Fields{}
	} else {
		l.logger.Errorf(format, args...)
	}

}

func (l *loggerLogrus) Warn(args ...interface{}) {
	if len(l.logrusFields) > 0 {
		finalLogrusFields := logrus.Fields{}
		for _, logrusFields := range l.logrusFields {
			for key, val := range logrusFields {
				finalLogrusFields[key] = val
			}
		}

		l.logger.WithFields(finalLogrusFields).Warn(args...)

		// reset fields
		l.logrusFields = []logrus.Fields{}
	} else {
		l.logger.Warn(args...)
	}

}

func (l *loggerLogrus) Warnf(format string, args ...interface{}) {
	if len(l.logrusFields) > 0 {
		finalLogrusFields := logrus.Fields{}
		for _, logrusFields := range l.logrusFields {
			for key, val := range logrusFields {
				finalLogrusFields[key] = val
			}
		}

		l.logger.WithFields(finalLogrusFields).Warnf(format, args...)

		// reset fields
		l.logrusFields = []logrus.Fields{}
	} else {
		l.logger.Warnf(format, args...)
	}

}

func (l *loggerLogrus) Info(args ...interface{}) {
	if len(l.logrusFields) > 0 {
		finalLogrusFields := logrus.Fields{}
		for _, logrusFields := range l.logrusFields {
			for key, val := range logrusFields {
				finalLogrusFields[key] = val
			}
		}

		l.logger.WithFields(finalLogrusFields).Info(args...)

		// reset fields
		l.logrusFields = []logrus.Fields{}
	} else {
		l.logger.Info(args...)
	}

}

func (l *loggerLogrus) Infof(format string, args ...interface{}) {
	if len(l.logrusFields) > 0 {
		finalLogrusFields := logrus.Fields{}
		for _, logrusFields := range l.logrusFields {
			for key, val := range logrusFields {
				finalLogrusFields[key] = val
			}
		}

		l.logger.WithFields(finalLogrusFields).Infof(format, args...)

		// reset fields
		l.logrusFields = []logrus.Fields{}
	} else {
		l.logger.Infof(format, args...)
	}

}

func (l *loggerLogrus) Debug(args ...interface{}) {
	if len(l.logrusFields) > 0 {
		finalLogrusFields := logrus.Fields{}
		for _, logrusFields := range l.logrusFields {
			for key, val := range logrusFields {
				finalLogrusFields[key] = val
			}
		}

		l.logger.WithFields(finalLogrusFields).Debug(args...)

		// reset fields
		l.logrusFields = []logrus.Fields{}
	} else {
		l.logger.Debug(args...)
	}

}

func (l *loggerLogrus) Debugf(format string, args ...interface{}) {
	if len(l.logrusFields) > 0 {
		finalLogrusFields := logrus.Fields{}
		for _, logrusFields := range l.logrusFields {
			for key, val := range logrusFields {
				finalLogrusFields[key] = val
			}
		}

		l.logger.WithFields(finalLogrusFields).Debugf(format, args...)

		// reset fields
		l.logrusFields = []logrus.Fields{}
	} else {
		l.logger.Debugf(format, args...)
	}

}

func (l *loggerLogrus) Trace(args ...interface{}) {
	if len(l.logrusFields) > 0 {
		finalLogrusFields := logrus.Fields{}
		for _, logrusFields := range l.logrusFields {
			for key, val := range logrusFields {
				finalLogrusFields[key] = val
			}
		}

		l.logger.WithFields(finalLogrusFields).Trace(args...)

		// reset fields
		l.logrusFields = []logrus.Fields{}
	} else {
		l.logger.Trace(args...)
	}

}

func (l *loggerLogrus) Tracef(format string, args ...interface{}) {
	if len(l.logrusFields) > 0 {
		finalLogrusFields := logrus.Fields{}
		for _, logrusFields := range l.logrusFields {
			for key, val := range logrusFields {
				finalLogrusFields[key] = val
			}
		}

		l.logger.WithFields(finalLogrusFields).Tracef(format, args...)

		// reset fields
		l.logrusFields = []logrus.Fields{}
	} else {
		l.logger.Tracef(format, args...)
	}

}

func (l *loggerLogrus) WithFields(fields Fields) Logger {
	logrusFields := logrus.Fields{}
	for key, val := range fields {
		logrusFields[key] = val
	}

	l.logrusFields = append(l.logrusFields, logrusFields)

	return l
}
