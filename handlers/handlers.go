// Package handlers provides a set of HTTP CRUD handlers.
package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"

	"github.com/danielkvist/botio/db"
	"github.com/danielkvist/botio/models"

	"github.com/go-chi/chi"
)

// Get returns an http.HandlerFunc which extracts from the request params
// a command name with which it tries to find an item in the database.
//
// Returns a non-2xx status code when:
// * There's a problem getting the item from the database.
// * There's a problem unmarshaling the item.
// * There's a problem encoding the result to JSON.
func Get(bolter db.Bolter, col string) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")

		command := chi.URLParam(r, "command")
		result, err := bolter.Get(col, command)
		if err != nil {
			http.Error(w, fmt.Sprintf("error while getting command %q from the database", command), http.StatusInternalServerError)
			log.Printf("while getting command %q from collection %q: %v", command, col, err)
			return
		}

		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(result); err != nil {
			w.Write([]byte(fmt.Sprintf("error while encoding JSON response with command %q", command)))
			log.Printf("while encoding JSON response for command %q: %v", command, err)
		}
	}
}

// GetAll returns an http.HandlerFunc which extracts all the commands
// from the commands collection of the database.
//
// Returns a non-2xx status code when:
// * There's a problem while getting the commands from the database.
// * There's a problem encoding the result to JSON.
func GetAll(bolter db.Bolter, col string) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")

		commands, err := bolter.GetAll(col)
		if err != nil {
			http.Error(w, fmt.Sprintf("error while getting all items in collection %q", col), http.StatusInternalServerError)
			log.Printf("while getting all commands from %q: %v", col, err)
			return
		}

		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(commands); err != nil {
			w.Write([]byte(fmt.Sprintf("error while encoding JSON response with all commands on collection %q", col)))
			log.Printf("while encoding JSON response with all commands on collection %q: %v", col, err)
		}
	}
}

// Post returns an http.HandlerFunc that tries to extract a new command
// from the body of the request and insert it into the database.
//
// Returns a non-2xx status code when:
// * The request body is too large.
// * There's a problem closing the request body.
// * There's a problem marshaling the new command.
// * There's a problem adding the new command into the database.
// * There's a problem encoding the result to JSON.
func Post(bolter db.Bolter, col string) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")

		if r.ContentLength == 0 {
			http.Error(w, "bad request body", http.StatusBadRequest)
			return
		}

		body, err := ioutil.ReadAll(io.LimitReader(r.Body, 10240))
		if err != nil {
			http.Error(w, "error while reading request body", http.StatusRequestEntityTooLarge)
			log.Printf("while reading request body: %v", err)
			return
		}

		if err := r.Body.Close(); err != nil {
			log.Printf("while trying to close request body: %v", err)
		}

		var cmd models.Command
		if err := json.Unmarshal(body, &cmd); err != nil {
			http.Error(w, "error while trying to unmarshal request body. entity not added to the database", http.StatusInternalServerError)
			log.Printf("while trying to unmarshal request body: %v", err)
			return
		}

		if cmd.Cmd == "" || cmd.Response == "" {
			http.Error(w, "empty fields are not allowed", http.StatusBadRequest)
			return
		}

		result, err := bolter.Set(col, cmd.Cmd, cmd.Response)
		if err != nil {
			http.Error(w, "error while trying to add command to the database", http.StatusInternalServerError)
			log.Printf("while trying to add command %q with response %q into collection %q: %v", cmd.Cmd, cmd.Response, col, err)
			return
		}

		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(result); err != nil {
			w.Write([]byte(fmt.Sprintf("error while encoding JSON response for added command %q", cmd.Cmd)))
			log.Printf("while encoding JSON response for added command %q: %v", cmd.Cmd, err)
		}
	}
}

// Put returns an http.HandlerFunc that tries to update an existing
// command on the database with the received request body.
//
// Returns a non-2xx status code when:
// * The request body is too large.
// * There's a problem closing the request body.
// * There's a problem unmarshaling the command.
// * There's a problem while updating the command on the database.
// * There's a problem encoding the result to JSON.
func Put(bolter db.Bolter, col string) func(w http.ResponseWriter, r *http.Request) {
	return Post(bolter, col)
}

// Delete returns an http.HandlerFunc which extracts from the request params
// a command to remove from the database.
//
// Returns a non-2xx status code when:
// * There's a problem while removing the command from the database.
func Delete(bolter db.Bolter, col string) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")

		command := chi.URLParam(r, "command")
		if err := bolter.Remove(col, command); err != nil {
			http.Error(w, fmt.Sprintf("error while removing command %q", command), http.StatusInternalServerError)
			log.Printf("while removing command %q from collection %q: %v", command, col, err)
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}

// Backup returns an http.HandlerFunc that will send to the client
// a backup of the database.
//
// Returns a non-2xx status code when if there is a problem while making the backup.
func Backup(bolter db.Bolter, col string) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/octet-stream")
		w.Header().Set("Content-Disposition", `attachment; filename="botio.db"`)

		length, err := bolter.Backup(w)
		w.Header().Set("Content-Length", strconv.Itoa(length))
		if err != nil {
			http.Error(w, fmt.Sprintf("error while triying to backup the dabatase: %v", err), http.StatusInternalServerError)
		}
	}
}
