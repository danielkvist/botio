package bot

import (
	"time"

	"github.com/sirupsen/logrus"
)

func logInfo(log *logrus.Logger, platform, id, command, response, msg string, t time.Duration) {
	log.WithFields(logrus.Fields{
		"platform": platform,
		"chatID":   id,
		"received": command,
		"sended":   response,
		"in":       t,
	}).Info(msg)
}

func logError(log *logrus.Logger, platform, pkg, function, id, command, err, msg string) {
	log.WithFields(logrus.Fields{
		"pkg":      pkg,
		"platform": platform,
		"func":     function,
		"chatID":   id,
		"received": command,
		"error":    err,
	}).Error(msg)
}

func logWarning(log *logrus.Logger, platform, pkg, function, id, command, err, msg string) {
	log.WithFields(logrus.Fields{
		"pkg":      pkg,
		"platform": platform,
		"func":     function,
		"chatID":   id,
		"received": command,
		"error":    err,
	}).Warning(msg)
}

func logFatal(log *logrus.Logger, platform, pkg, function, id, command, err, msg string) {
	log.WithFields(logrus.Fields{
		"pkg":      pkg,
		"platform": platform,
		"func":     function,
		"chatID":   id,
		"received": command,
		"error":    err,
	}).Fatal(msg)
}
