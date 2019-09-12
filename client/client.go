// Package client exports functions to manage botio's commands with HTTP methods
// and JWT based authentication.
package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/danielkvist/botio/models"

	"github.com/dgrijalva/jwt-go"
)

// Get receives an URL and a key to perform an HTTP GET request
// using the key for authentication using JWT and returns a *models.Command.
// If something goes wrong it returns a non-nil error.
func Get(url, key string) (*models.Command, error) {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("while creating a new request for %q: %v", url, err)
	}

	if err := reqWithTokenHeader(req, key); err != nil {
		return nil, fmt.Errorf("%v", err)
	}

	c := &http.Client{}
	resp, err := c.Do(req)
	if err != nil {
		return nil, fmt.Errorf("while making a request for %q: %v", url, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("while making a GET request for %q a %v status code was expected. got %q", url, http.StatusOK, resp.Status)
	}

	d, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("while reading the response body from %q: %v", url, err)
	}

	var cmd models.Command
	if err := json.Unmarshal(d, &cmd); err != nil {
		return nil, fmt.Errorf("while unmarshaling the response body from %q: %v", url, err)
	}

	return &cmd, nil
}

// GetAll receives an URL and a key to perform an HTTP GET request
// using the key for authentication using JWT and returns an []*models.Command.
// If something goes wrong it returns a non-nil error.
func GetAll(url, key string) ([]*models.Command, error) {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("while creating a new request for %q: %v", url, err)
	}

	if err := reqWithTokenHeader(req, key); err != nil {
		return nil, fmt.Errorf("%v", err)
	}

	c := &http.Client{}
	resp, err := c.Do(req)
	if err != nil {
		return nil, fmt.Errorf("while making a request for %q: %v", url, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("while making a GET request for %q a %v status code was expected. got %q", url, http.StatusOK, resp.Status)
	}

	d, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("while reading the response body from %q: %v", url, err)
	}

	commands := []*models.Command{}
	if err := json.Unmarshal(d, &commands); err != nil {
		return nil, fmt.Errorf("while unmarshaling the response body from %q: %v", url, err)
	}

	return commands, nil
}

// Post receives an URL and a key with a command and a response and performs
// an HTTP POST request using the command and the response as the body of the request
// and the key for authentication using JWT.
// If something goes wrong it returns a non-nil error.
func Post(url, key, command, response string) (*models.Command, error) {
	if command == "" || response == "" {
		return nil, fmt.Errorf("empty fields are not allowed")
	}

	cmd := models.Command{
		Cmd:      command,
		Response: response,
	}
	cmdData, err := json.Marshal(cmd)
	if err != nil {
		return nil, fmt.Errorf("while marshaling command %q: %v", command, err)
	}

	req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(cmdData))
	if err != nil {
		return nil, fmt.Errorf("while making a request for %q: %v", url, err)

	}

	if err := reqWithTokenHeader(req, key); err != nil {
		return nil, fmt.Errorf("%v", err)
	}

	c := &http.Client{}
	resp, err := c.Do(req)
	if err != nil {
		return nil, fmt.Errorf("while making a request for %q: %v", url, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("while making a POST request for %q a %v status code was expected. got %q", url, http.StatusOK, resp.Status)
	}

	return &cmd, nil
}

// Put receives an URL and a key with a command and a response and performs
// an HTTP POST request using the command and the response as the body of the request
// and the key for authentication using JWT.
// It performs an HTTP POST request instead of an HTTP PUT request
// due to how BoltDB databases work.
// If something goes wrong it returns a non-nil error.
func Put(url, key, command, response string) (*models.Command, error) {
	return Post(url, key, command, response)
}

// Delete receives an URL and a key and performs an HTTP DELETE request
// using the key for authentication using JWT.
// If something goes wrong it returns a non-nil error.
func Delete(url, key string) error {
	req, err := http.NewRequest(http.MethodDelete, url, nil)
	if err != nil {
		return fmt.Errorf("while creating a new request for %q: %v", url, err)
	}

	if err := reqWithTokenHeader(req, key); err != nil {
		return fmt.Errorf("%v", err)
	}

	c := &http.Client{}
	resp, err := c.Do(req)
	if err != nil {
		return fmt.Errorf("while making a request for %q: %v", url, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("while making a DELETE request for %q a %v status code was expected. got %q", url, http.StatusOK, resp.Status)
	}

	return nil
}

func reqWithTokenHeader(r *http.Request, key string) error {
	token, err := generate(key)
	if err != nil {
		return fmt.Errorf("while generating JWT token for authentication: %v", err)
	}

	r.Header.Set("Token", token)

	return nil
}

func generate(key string) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)

	tkStr, err := token.SignedString([]byte(key))
	if err != nil {
		return "", fmt.Errorf("while generating the authentication JWT token: %v", err)
	}

	return tkStr, nil

}
