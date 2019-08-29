package logger

import "github.com/sirupsen/logrus"

type Logger struct {
	l *logrus.Logger
}

func New() *Logger {
	logger := &Logger{}
	l := logrus.New()
	l.SetFormatter(&logrus.TextFormatter{
		FullTimestamp:  true,
		DisableSorting: true,
	})

	logger.l = l
	return logger
}

func (l *Logger) Info(msg string) {
	l.l.Info(msg)
}

func (l *Logger) Warning(msg string) {
	l.l.Warning(msg)
}

func (l *Logger) Error(msg string) {
	l.l.Error(msg)
}

func (l *Logger) Fatal(msg string) {
	l.l.Fatal(msg)
}

func (l *Logger) Panic(msg string) {
	l.l.Panic(msg)
}
