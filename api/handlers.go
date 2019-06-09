package api

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/danielkvist/botio/db"
	"github.com/danielkvist/botio/models"
	"github.com/gorilla/mux"
)

func Get(bolter db.Bolter, col string) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")

		params := mux.Vars(r)
		command := params["command"]

		result, err := bolter.Get(col, command)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(fmt.Sprintf("error while getting command %q from the database", command)))
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

func GetAll(bolter db.Bolter, col string) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")

		commands, err := bolter.GetAll(col)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(fmt.Sprintf("error while getting all items in collection %q", col)))
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

func Post(bolter db.Bolter, col string) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")

		if r.ContentLength == 0 {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("error bad request body"))
			return
		}

		body, err := ioutil.ReadAll(io.LimitReader(r.Body, 10240))
		if err != nil {
			w.WriteHeader(http.StatusRequestEntityTooLarge)
			w.Write([]byte("error while reading your request body"))
			log.Printf("while reading request body: %v", err)
			return
		}

		if err := r.Body.Close(); err != nil {
			log.Printf("while trying to close request body: %v", err)
		}

		var cmd models.Command
		if err := json.Unmarshal(body, &cmd); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("error while trying to unmarshal request body. entity not added to the database"))
			log.Printf("while trying to unmarshal request body: %v", err)
			return
		}

		if cmd.Cmd == "" || cmd.Response == "" {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("any field empty is not allowed"))
			return
		}

		result, err := bolter.Set(col, cmd.Cmd, cmd.Response)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("error while trying to add command to the database"))
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

func Put(bolter db.Bolter, col string) func(w http.ResponseWriter, r *http.Request) {
	return Post(bolter, col)
}

func Delete(bolter db.Bolter, col string) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")

		params := mux.Vars(r)
		command := params["command"]

		if err := bolter.Remove(col, command); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(fmt.Sprintf("error while removing command %q", command)))
			log.Printf("while removing command %q from collection %q: %v", command, col, err)
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}
