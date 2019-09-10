// Package logger returns a wrapper for a logger.
package logger

import (
	"net/http"
	"os"

	"github.com/sirupsen/logrus"
)

// Logger is a simple wrapper for a logger like logrus.
type Logger struct {
	l *logrus.Logger
}

// New returns a new *Logger totally initialized.
func New() *Logger {
	logger := &Logger{}
	l := logrus.New()
	l.SetFormatter(&logrus.TextFormatter{
		FullTimestamp:    true,
		QuoteEmptyFields: true,
		TimestampFormat:  "02-01-2006 15:04:05",
	})
	l.Out = os.Stdout

	logger.l = l
	return logger
}

// LogMsg logs a message information received by a bot.
func (l *Logger) LogMsg(platform string, id string, msg string, response string) {
	log := l.l.WithFields(logrus.Fields{
		"platform": platform,
		"id":       id,
		"msg":      msg,
		"response": response,
	})

	log.Info("command handled")
}

// LogRequest logs a request information handled by a handler.
func (l *Logger) LogRequest(w http.ResponseWriter, r *http.Request, status int, err error) {
	log := l.l.WithFields(logrus.Fields{
		"method":   r.Method,
		"url":      r.URL.String(),
		"host":     r.Host,
		"from":     r.RemoteAddr,
		"status":   status,
		"content":  w.Header().Get("Content-Type"),
		"protocol": r.Proto,
	})

	if err != nil {
		log.Error(err.Error())
		return
	}

	log.Info("request handled successfully")
}
