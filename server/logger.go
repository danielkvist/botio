package server

import (
	"net/http"

	"github.com/sirupsen/logrus"
)

func (s *Server) log(w http.ResponseWriter, r *http.Request, status int, err error) {
	log := s.logger.WithFields(logrus.Fields{
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
