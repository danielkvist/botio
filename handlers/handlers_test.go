package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/danielkvist/botio/models"

	"github.com/gorilla/mux"
)

type dbt map[string]string

func (db dbt) Set(col string, el string, val string) (*models.Command, error) {
	db[el] = val

	return &models.Command{
		Cmd:      el,
		Response: val,
	}, nil
}

func (db dbt) Get(col string, el string) (*models.Command, error) {
	val, ok := db[el]
	if !ok {
		return nil, fmt.Errorf("element %q not found", el)
	}

	return &models.Command{
		Cmd:      el,
		Response: val,
	}, nil
}

func (db dbt) GetAll(col string) ([]*models.Command, error) {
	var commands []*models.Command

	for k, v := range db {
		tmpCommand := &models.Command{
			Cmd:      k,
			Response: v,
		}

		commands = append(commands, tmpCommand)
	}

	return commands, nil
}

func (db dbt) Remove(col string, el string) error {
	delete(db, el)
	return nil
}

func (db dbt) Update(col string, el string, val string) (*models.Command, error) {
	db[el] = val

	return &models.Command{
		Cmd:      el,
		Response: val,
	}, nil
}

func (db dbt) Backup(w io.Writer) (int, error) {
	return 0, nil
}

func TestGet(t *testing.T) {
	tt := []struct {
		cmd      string
		response string
	}{
		{"start", "Hi!"},
		{"response", "42"},
		{"goodbye", "Ciao!"},
	}

	var db dbt = make(map[string]string, len(tt))
	for _, tc := range tt {
		db[tc.cmd] = tc.response
	}

	for _, tc := range tt {
		path := fmt.Sprintf("/api/commands/%s", tc.cmd)
		req, err := http.NewRequest(http.MethodGet, path, nil)
		if err != nil {
			t.Fatalf("while creating a new request to route %q: %v", path, err)
		}

		rec := httptest.NewRecorder()
		router := mux.NewRouter()
		router.HandleFunc("/api/commands/{command}", Get(db, ""))
		router.ServeHTTP(rec, req)

		if rec.Code != http.StatusOK {
			t.Fatalf("expected status of response from %q to be %v. got=%v", path, http.StatusOK, rec.Code)
		}

		body, err := ioutil.ReadAll(rec.Body)
		if err != nil {
			t.Fatalf("while reading the body of the response from %q: %v", path, err)
		}

		var cmd models.Command
		if err := json.Unmarshal(body, &cmd); err != nil {
			t.Fatalf("while unmarshaling response body from %q: %v", path, err)
		}

		if cmd.Cmd != tc.cmd && cmd.Response != tc.response {
			t.Fatalf("expected to get command %q with response %q. got command %q with response %q", tc.cmd, tc.response, cmd.Cmd, cmd.Response)
		}
	}
}

func TestGetAll(t *testing.T) {
	tt := []struct {
		cmd      string
		response string
	}{
		{"start", "Hi!"},
		{"name", "Go"},
		{"occupation", "Gopher"},
		{"do", "test"},
		{"aaa", "42"},
		{"none", ""},
	}

	var db dbt = make(map[string]string, len(tt))
	for _, tc := range tt {
		db[tc.cmd] = tc.response
	}

	path := fmt.Sprintf("/api/commands")
	req, err := http.NewRequest(http.MethodGet, path, nil)
	if err != nil {
		t.Fatalf("while creating a new request to route %q: %v", path, err)
	}

	rec := httptest.NewRecorder()
	router := mux.NewRouter()
	router.HandleFunc("/api/commands", GetAll(db, ""))
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status of response from %q to be %v. got=%v", path, http.StatusOK, rec.Code)
	}

	body, err := ioutil.ReadAll(rec.Body)
	if err != nil {
		t.Fatalf("while reading the body of the response from %q: %v", path, err)
	}

	var commands []*models.Command
	if err := json.Unmarshal(body, &commands); err != nil {
		t.Fatalf("while unmarshaling response body from %q: %v", path, err)
	}

	for _, cmd := range commands {
		resp, ok := db[cmd.Cmd]
		if !ok {
			t.Fatalf("expected to find command %q in testing db", cmd.Cmd)
		}

		if cmd.Response != resp {
			t.Fatalf("expected command %q to have response %q. got=%q", cmd.Cmd, cmd.Response, resp)
		}
	}
}

func TestPost(t *testing.T) {
	var db dbt = make(map[string]string, 1)

	cmd := models.Command{
		Cmd:      "start",
		Response: "Hi!",
	}

	path := fmt.Sprintf("/api/commands")
	reqBody := []byte(fmt.Sprintf("{\"cmd\": %q, \"response\": %q}", cmd.Cmd, cmd.Response))
	req, err := http.NewRequest(http.MethodPost, path, bytes.NewReader(reqBody))
	if err != nil {
		t.Fatalf("while creating a new request to route %q: %v", path, err)
	}

	rec := httptest.NewRecorder()
	router := mux.NewRouter()
	router.HandleFunc("/api/commands", Post(db, ""))
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status of response from %q to be %v. got=%v", path, http.StatusOK, rec.Code)
	}

	resp, ok := db[cmd.Cmd]
	if !ok {
		t.Fatalf("expected to find command %q in testing db as result of POST request", cmd.Cmd)
	}

	if resp != cmd.Response {
		t.Fatalf("expected command %q to have response %q. got=%q", cmd.Cmd, cmd.Response, resp)
	}
}

func TestDelete(t *testing.T) {
	var db dbt = make(map[string]string, 1)

	cmd := models.Command{
		Cmd:      "start",
		Response: "Hi!",
	}

	db[cmd.Cmd] = cmd.Response

	path := fmt.Sprintf("/api/commands/%s", cmd.Cmd)
	req, err := http.NewRequest(http.MethodDelete, path, nil)
	if err != nil {
		t.Fatalf("while creating a new request to route %q: %v", path, err)
	}

	rec := httptest.NewRecorder()
	router := mux.NewRouter()
	router.HandleFunc("/api/commands/{command}", Delete(db, ""))
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status of response from %q to be %v. got=%v", path, http.StatusOK, rec.Code)
	}

	_, ok := db[cmd.Cmd]
	if ok {
		t.Fatalf("expected to not find command %q in testing db as result of DELETE request", cmd.Cmd)
	}
}
