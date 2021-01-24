package logger

import (
	"os"

	"github.com/sirupsen/logrus"
)

const (
	// logTimeFormat represents time format in log messages
	logTimeFormat  = "2006-01-02 15:04:05.99"
	LogLevelEnvKey = "LOG_LEVEL"
)

type Logger interface {
	Debug(args ...interface{})
	Debugf(format string, args ...interface{})
	Info(args ...interface{})
	Infof(format string, args ...interface{})
	Warn(args ...interface{})
	Warnf(format string, args ...interface{})
	Error(args ...interface{})
	Errorf(format string, args ...interface{})
	Fatal(args ...interface{})
	Fatalf(format string, args ...interface{})
}

type logger struct {
	*logrus.Logger
}

func New() Logger {
	customFormatter := new(logrus.TextFormatter)
	customFormatter.TimestampFormat = logTimeFormat
	customFormatter.FullTimestamp = true
	l := logger{}
	l.Logger = logrus.New()

	l.SetFormatter(customFormatter)
	mode := os.Getenv(LogLevelEnvKey)
	if mode == "info" {
		l.SetLevel(logrus.InfoLevel)
	} else {
		l.SetLevel(logrus.DebugLevel)
	}
	return l
}
