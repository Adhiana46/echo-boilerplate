package logger

import (
	logGo "log"
)

var log Logger

type Fields map[string]interface{}

// Level
const (
	LevelPanic = "panic"
	LevelFatal = "fatal"
	LevelError = "error"
	LevelWarn  = "warn"
	LevelInfo  = "info"
	LevelDebug = "debug"
	LevelTrace = "trace"
)

type Config struct {
	Level           string
	TimestampFormat string
	FileLocation    string
}

type Logger interface {
	Panic(args ...interface{})
	Panicf(format string, args ...interface{})

	Fatal(args ...interface{})
	Fatalf(format string, args ...interface{})

	Error(args ...interface{})
	Errorf(format string, args ...interface{})

	Warn(args ...interface{})
	Warnf(format string, args ...interface{})

	Info(args ...interface{})
	Infof(format string, args ...interface{})

	Debug(args ...interface{})
	Debugf(format string, args ...interface{})

	Trace(args ...interface{})
	Tracef(format string, args ...interface{})

	WithFields(fields Fields) Logger
}

func SetLogger(l Logger) {
	log = l
}

func Panic(args ...interface{}) {
	if log == nil {
		logGo.Println(args...)
	} else {
		log.Panic(args...)
	}
}

func Panicf(format string, args ...interface{}) {
	if log == nil {
		logGo.Printf(format, args...)
	} else {
		log.Panicf(format, args...)
	}
}

func Fatal(args ...interface{}) {
	if log == nil {
		logGo.Println(args...)
	} else {
		log.Fatal(args...)
	}
}

func Fatalf(format string, args ...interface{}) {
	if log == nil {
		logGo.Printf(format, args...)
	} else {
		log.Fatalf(format, args...)
	}
}

func Error(args ...interface{}) {
	if log == nil {
		logGo.Println(args...)
	} else {
		log.Error(args...)
	}
}

func Errorf(format string, args ...interface{}) {
	if log == nil {
		logGo.Printf(format, args...)
	} else {
		log.Errorf(format, args...)
	}
}

func Warn(args ...interface{}) {
	if log == nil {
		logGo.Println(args...)
	} else {
		log.Warn(args...)
	}
}

func Warnf(format string, args ...interface{}) {
	if log == nil {
		logGo.Printf(format, args...)
	} else {
		log.Warnf(format, args...)
	}
}

func Info(args ...interface{}) {
	if log == nil {
		logGo.Println(args...)
	} else {
		log.Info(args...)
	}
}

func Infof(format string, args ...interface{}) {
	if log == nil {
		logGo.Printf(format, args...)
	} else {
		log.Infof(format, args...)
	}
}

func Debug(args ...interface{}) {
	if log == nil {
		logGo.Println(args...)
	} else {
		log.Debug(args...)
	}
}

func Debugf(format string, args ...interface{}) {
	if log == nil {
		logGo.Printf(format, args...)
	} else {
		log.Debugf(format, args...)
	}
}

func Trace(args ...interface{}) {
	if log == nil {
		logGo.Println(args...)
	} else {
		log.Trace(args...)
	}
}

func Tracef(format string, args ...interface{}) {
	if log == nil {
		logGo.Printf(format, args...)
	} else {
		log.Tracef(format, args...)
	}
}

func Println(args ...interface{}) {
	if log == nil {
		logGo.Println(args...)
	} else {
		log.Info(args...)
	}
}

func Printf(format string, args ...interface{}) {
	if log == nil {
		logGo.Printf(format, args...)
	} else {
		log.Infof(format, args...)
	}
}

func WithFields(fields Fields) Logger {
	return log.WithFields(fields)
}
