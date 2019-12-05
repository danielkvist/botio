package server

import (
	"time"

	"github.com/sirupsen/logrus"
)

func (s *server) logInfo(pkg, function, msg string, t time.Duration) {
	s.log.WithFields(logrus.Fields{
		"db":       s.dbPlatform,
		"cache":    s.cachePlatform,
		"grpcAddr": s.listener.Addr().String(),
		"httpPort": s.httpPort,
		"ssl":      s.ssl,
		"pkg":      pkg,
		"func":     function,
		"in":       t,
	}).Info(msg)
}

func (s *server) logError(pkg, function, err, msg string) {
	s.log.WithFields(logrus.Fields{
		"db":       s.dbPlatform,
		"cache":    s.cachePlatform,
		"grpcAddr": s.listener.Addr().String(),
		"httpPort": s.httpPort,
		"ssl":      s.ssl,
		"pkg":      pkg,
		"func":     function,
		"error":    err,
	}).Error(msg)
}

func (s *server) logWarning(pkg, function, err, msg string) {
	s.log.WithFields(logrus.Fields{
		"db":       s.dbPlatform,
		"cache":    s.cachePlatform,
		"grpcAddr": s.listener.Addr().String(),
		"httpPort": s.httpPort,
		"ssl":      s.ssl,
		"pkg":      pkg,
		"func":     function,
		"error":    err,
	}).Warning(msg)
}

func (s *server) logFatal(pkg, function, err, msg string) {
	s.log.WithFields(logrus.Fields{
		"db":       s.dbPlatform,
		"cache":    s.cachePlatform,
		"grpcAddr": s.listener.Addr().String(),
		"httpPort": s.httpPort,
		"ssl":      s.ssl,
		"pkg":      pkg,
		"func":     function,
		"error":    err,
	}).Fatal(msg)
}
