package bot

import "github.com/sirupsen/logrus"

func log(logger *logrus.Logger, platform, id, msg, response string) {
	logger.WithFields(logrus.Fields{
		"platform": platform,
		"id":       id,
		"msg":      msg,
		"response": response,
	}).Info("command handled")
}
