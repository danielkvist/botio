package client

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/danielkvist/botio/models"

	"github.com/dgrijalva/jwt-go"
)

func TestGet(t *testing.T) {
	jwtKey := "0123456789abcdef"
	command := &models.Command{Cmd: "Hi", Response: "World"}
	ts := testGet(t, jwtKey, command)
	defer ts.Close()

	tt := []struct {
		name            string
		key             string
		url             string
		expectedToFail  bool
		expectedCommand *models.Command
	}{
		{
			name:            "normal request",
			key:             jwtKey,
			url:             ts.URL,
			expectedToFail:  false,
			expectedCommand: command,
		},
		{
			name:           "no key",
			url:            ts.URL,
			expectedToFail: true,
		},
		{
			name:           "no URL",
			key:            jwtKey,
			expectedToFail: true,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			command, err := Get(tc.url, tc.key)
			if err != nil {
				if tc.expectedToFail {
					t.Logf("client request expected to fail failed: %v", err)
					return
				}
				t.Fatalf("client request not expected to fail: %v", err)
			}

			if tc.expectedToFail {
				t.Fatalf("client request expected to fail not failed")
			}

			if command.Cmd != tc.expectedCommand.Cmd {
				t.Fatalf("command has wrong Cmd value, expected %q. got=%v", tc.expectedCommand.Cmd, command.Cmd)
			}

			if command.Response != tc.expectedCommand.Response {
				t.Fatalf("command has wrong Response value, expected %q. got=%v", tc.expectedCommand.Response, command.Response)
			}
		})
	}
}

func TestGetAll(t *testing.T) {
	jwtKey := "0123456789abcdef"
	commands := []*models.Command{
		&models.Command{
			Cmd:      "A",
			Response: "a",
		},
		&models.Command{
			Cmd:      "B",
			Response: "b",
		},
		&models.Command{
			Cmd:      "C",
			Response: "c",
		},
	}
	ts := testGet(t, jwtKey, commands)
	defer ts.Close()

	tt := []struct {
		name             string
		key              string
		url              string
		expectedToFail   bool
		expectedCommands []*models.Command
	}{
		{
			name:             "normal request",
			key:              jwtKey,
			url:              ts.URL,
			expectedToFail:   false,
			expectedCommands: commands,
		},
		{
			name:           "no key",
			url:            ts.URL,
			expectedToFail: true,
		},
		{
			name:           "no URL",
			key:            jwtKey,
			expectedToFail: true,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			commands, err := GetAll(tc.url, tc.key)
			if err != nil {
				if tc.expectedToFail {
					t.Logf("client request expected to fail failed: %v", err)
					return
				}
				t.Fatalf("client request not expected to fail: %v", err)
			}

			if tc.expectedToFail {
				t.Fatalf("client request expected to fail not failed")
			}

			for i, ec := range tc.expectedCommands {
				if ec.Cmd != commands[i].Cmd {
					t.Fatalf("command has wrong Cmd value, expected %q. got=%v", ec.Cmd, commands[i].Cmd)
				}

				if ec.Response != commands[i].Response {
					t.Fatalf("command has wrong Response value, expected %q. got=%v", ec.Response, commands[i].Response)
				}
			}
		})
	}
}

func TestPost(t *testing.T) {
	jwtKey := "0123456789abcdef"

	ts := testGet(t, jwtKey, nil)
	defer ts.Close()

	tt := []struct {
		name           string
		key            string
		url            string
		command        string
		response       string
		expectedToFail bool
	}{
		{
			name:           "normal request",
			key:            jwtKey,
			url:            ts.URL,
			command:        "Hi",
			response:       "Hello",
			expectedToFail: false,
		},
		{
			name:           "no key",
			url:            ts.URL,
			command:        "Hi",
			response:       "Hello",
			expectedToFail: true,
		},
		{
			name:           "no URL",
			key:            jwtKey,
			command:        "Hi",
			response:       "Hello",
			expectedToFail: true,
		},
		{
			name:           "no command",
			key:            jwtKey,
			url:            ts.URL,
			response:       "Hello",
			expectedToFail: true,
		},
		{
			name:           "no response",
			key:            jwtKey,
			url:            ts.URL,
			command:        "Hi",
			expectedToFail: true,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			command, err := Post(tc.url, tc.key, tc.command, tc.response)
			if err != nil {
				if tc.expectedToFail {
					t.Logf("client request expected to fail failed: %v", err)
					return
				}
				t.Fatalf("client request not expected to fail failed: %v", err)
			}

			if tc.expectedToFail {
				t.Fatalf("client request expected to fail not failed")
			}

			if command.Cmd != tc.command {
				t.Fatalf("command has wrong Cmd value, expected %q. got=%v", tc.command, command.Cmd)
			}

			if command.Response != tc.response {
				t.Fatalf("command has wrong Response value, expected %q. got=%v", tc.response, command.Response)
			}
		})
	}
}

func TestPut(t *testing.T) {
	jwtKey := "0123456789abcdef"

	ts := testGet(t, jwtKey, nil)
	defer ts.Close()

	tt := []struct {
		name           string
		key            string
		url            string
		command        string
		response       string
		expectedToFail bool
	}{
		{
			name:           "normal request",
			key:            jwtKey,
			url:            ts.URL,
			command:        "Hi",
			response:       "Hello",
			expectedToFail: false,
		},
		{
			name:           "no key",
			url:            ts.URL,
			command:        "Hi",
			response:       "Hello",
			expectedToFail: true,
		},
		{
			name:           "no URL",
			key:            jwtKey,
			command:        "Hi",
			response:       "Hello",
			expectedToFail: true,
		},
		{
			name:           "no command",
			key:            jwtKey,
			url:            ts.URL,
			response:       "Hello",
			expectedToFail: true,
		},
		{
			name:           "no response",
			key:            jwtKey,
			url:            ts.URL,
			command:        "Hi",
			expectedToFail: true,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			command, err := Put(tc.url, tc.key, tc.command, tc.response)
			if err != nil {
				if tc.expectedToFail {
					t.Logf("client request expected to fail failed: %v", err)
					return
				}
				t.Fatalf("client request not expected to fail failed: %v", err)
			}

			if tc.expectedToFail {
				t.Fatalf("client request expected to fail not failed")
			}

			if command.Cmd != tc.command {
				t.Fatalf("command has wrong Cmd value, expected %q. got=%v", tc.command, command.Cmd)
			}

			if command.Response != tc.response {
				t.Fatalf("command has wrong Response value, expected %q. got=%v", tc.response, command.Response)
			}
		})
	}
}

func TestDelete(t *testing.T) {
	jwtKey := "0123456789abcdef"

	ts := testGet(t, jwtKey, nil)
	defer ts.Close()

	tt := []struct {
		name           string
		key            string
		url            string
		expectedToFail bool
	}{
		{
			name:           "normal request",
			key:            jwtKey,
			url:            ts.URL,
			expectedToFail: false,
		},
		{
			name:           "no key",
			url:            ts.URL,
			expectedToFail: true,
		},
		{
			name:           "no URL",
			key:            jwtKey,
			expectedToFail: true,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			if err := Delete(tc.url, tc.key); err != nil {
				if tc.expectedToFail {
					t.Logf("client request expected to fail failed: %v", err)
					return
				}
				t.Fatalf("client request not expected to fail: %v", err)
			}

			if tc.expectedToFail {
				t.Fatalf("client request expected to fail not failed")
			}
		})
	}
}

func testGet(t *testing.T, key string, data interface{}) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		if r.Header["Token"] == nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		token, err := jwt.Parse(r.Header["Token"][0], func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("error while parsing JWT token")
			}

			return []byte(key), nil
		})

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		if !token.Valid {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		if err := json.NewEncoder(w).Encode(data); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
	}))
}
