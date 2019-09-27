package server

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/danielkvist/botio/models"

	"github.com/go-chi/chi"
)

func (s *Server) handleGet() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")

		command := chi.URLParam(r, "command")
		result, err := s.db.Get(command)
		if err != nil {
			s.errResp(w, r, "error while getting item", err, http.StatusInternalServerError)
			return
		}

		s.encodedResp(w, r, result)
	}
}

func (s *Server) handleGetAll() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")

		commands, err := s.db.GetAll()
		if err != nil {
			s.errResp(w, r, "error while getting items", err, http.StatusInternalServerError)
			return
		}

		s.encodedResp(w, r, commands)
	}
}

func (s *Server) handlePost() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")

		if r.ContentLength == 0 {
			s.errResp(w, r, "bad request body", fmt.Errorf("bad request body"), http.StatusBadRequest)
			return
		}

		body, err := ioutil.ReadAll(io.LimitReader(r.Body, 10240))
		if err != nil {
			s.errResp(w, r, "error while reading request body", err, http.StatusRequestEntityTooLarge)
			return
		}

		if err := r.Body.Close(); err != nil {
			s.logger.LogRequest(w, r, 0, err)
		}

		var cmd models.Command
		if err := json.Unmarshal(body, &cmd); err != nil {
			s.errResp(w, r, "error while trying to process request body", err, http.StatusUnprocessableEntity)
			return
		}

		if cmd.Cmd == "" || cmd.Response == "" {
			s.errResp(w, r, "empty fields not allowed", fmt.Errorf("empty fields not allowed"), http.StatusBadRequest)
			return
		}

		result, err := s.db.Add(cmd.Cmd, cmd.Response)
		if err != nil {
			s.errResp(w, r, "error while adding command to the database", err, http.StatusInternalServerError)
			return
		}

		s.encodedResp(w, r, result)
	}
}

func (s *Server) handlePut() func(w http.ResponseWriter, r *http.Request) {
	return s.handlePost()
}

func (s *Server) handleDelete() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")

		command := chi.URLParam(r, "command")
		if err := s.db.Remove(command); err != nil {
			s.errResp(w, r, "error while removing command", err, http.StatusInternalServerError)
			return
		}

		s.encodedResp(w, r, nil)
	}
}

func (s *Server) errResp(w http.ResponseWriter, r *http.Request, msg string, err error, status int) {
	http.Error(w, msg, status)
	s.logger.LogRequest(w, r, status, err)
}

func (s *Server) encodedResp(w http.ResponseWriter, r *http.Request, data interface{}) {
	w.WriteHeader(http.StatusOK)
	if data == nil {
		s.logger.LogRequest(w, r, http.StatusOK, nil)
		return
	}

	if err := json.NewEncoder(w).Encode(data); err != nil {
		s.errResp(w, r, "error while encoding data for response", err, http.StatusInternalServerError)
		return
	}

	s.logger.LogRequest(w, r, http.StatusOK, nil)
}
