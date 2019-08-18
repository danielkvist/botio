package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/danielkvist/botio/db"
	"github.com/danielkvist/botio/models"

	"github.com/go-chi/chi"
)

func TestGet(t *testing.T) {
	tt := []struct {
		cmd      string
		response string
	}{
		{"start", "Hi!"},
		{"response", "42"},
		{"goodbye", "Ciao!"},
	}

	mdb := db.DBFactory("testing")
	for _, tc := range tt {
		mdb.Set("", tc.cmd, tc.response)
	}

	for _, tc := range tt {
		path := fmt.Sprintf("/api/commands/%s", tc.cmd)
		req, err := http.NewRequest(http.MethodGet, path, nil)
		if err != nil {
			t.Fatalf("while creating a new request to route %q: %v", path, err)
		}

		rec := httptest.NewRecorder()
		router := chi.NewRouter()
		router.HandleFunc("/api/commands/{command}", Get(mdb, ""))
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

	mdb := db.DBFactory("testing")
	for _, tc := range tt {
		mdb.Set("", tc.cmd, tc.response)
	}

	path := fmt.Sprintf("/api/commands")
	req, err := http.NewRequest(http.MethodGet, path, nil)
	if err != nil {
		t.Fatalf("while creating a new request to route %q: %v", path, err)
	}

	rec := httptest.NewRecorder()
	router := chi.NewRouter()
	router.HandleFunc("/api/commands", GetAll(mdb, ""))
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
		c, err := mdb.Get("", cmd.Cmd)
		if err != nil {
			t.Fatalf("expected to find command %q in testing db", cmd.Cmd)
		}

		if cmd.Response != c.Response {
			t.Fatalf("expected command %q to have response %q. got=%q", cmd.Cmd, cmd.Response, c.Response)
		}
	}
}

func TestPost(t *testing.T) {
	mdb := db.DBFactory("testing")
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
	router := chi.NewRouter()
	router.HandleFunc("/api/commands", Post(mdb, ""))
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status of response from %q to be %v. got=%v", path, http.StatusOK, rec.Code)
	}

	c, err := mdb.Get("", cmd.Cmd)
	if err != nil {
		t.Fatalf("expected to find command %q in testing db as result of POST request", cmd.Cmd)
	}

	if cmd.Response != c.Response {
		t.Fatalf("expected command %q to have response %q. got=%q", cmd.Cmd, cmd.Response, c.Response)
	}
}

func TestDelete(t *testing.T) {
	mdb := db.DBFactory("testing")
	cmd := models.Command{
		Cmd:      "start",
		Response: "Hi!",
	}

	mdb.Set("", cmd.Cmd, cmd.Response)

	path := fmt.Sprintf("/api/commands/%s", cmd.Cmd)
	req, err := http.NewRequest(http.MethodDelete, path, nil)
	if err != nil {
		t.Fatalf("while creating a new request to route %q: %v", path, err)
	}

	rec := httptest.NewRecorder()
	router := chi.NewRouter()
	router.HandleFunc("/api/commands/{command}", Delete(mdb, ""))
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status of response from %q to be %v. got=%v", path, http.StatusOK, rec.Code)
	}

	_, err = mdb.Get("", cmd.Cmd)
	if err == nil {
		t.Fatalf("expected to not find command %q in testing db as result of DELETE request", cmd.Cmd)
	}
}
