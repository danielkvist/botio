// Package client exports simple methods to get the commands
// from the server.
package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/danielkvist/botio/models"
)

// Get receives an URL to which make a request using the received username
// and password for basic authentication with the objective
// of get a botio's command and return it.
// If something goes wrong while making the request it returns a non-nil error.
func Get(url, username, password string) (*models.Command, error) {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("while creating a new request for %q: %v", url, err)
	}
	req.SetBasicAuth(username, password)

	c := &http.Client{}

	resp, err := c.Do(req)
	if err != nil {
		return nil, fmt.Errorf("while making a request for %q: %v", url, err)
	}
	defer resp.Body.Close()

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

// GetAll receives an URL to which make a request using the received username
// and password for basic authentication with the objective
// of get all botio's commands and return them.
// If something goes wrong while making the request it returns a non-nil error.
func GetAll(url, username, password string) ([]*models.Command, error) {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("while creating a new request for %q: %v", url, err)
	}
	req.SetBasicAuth(username, password)

	c := &http.Client{}
	resp, err := c.Do(req)
	if err != nil {
		return nil, fmt.Errorf("while making a request for %q: %v", url, err)
	}
	defer resp.Body.Close()

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

// TODO:
func Post(url, username, password, command, response string) (*models.Command, error) {
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
	req.SetBasicAuth(username, password)

	c := &http.Client{}
	resp, err := c.Do(req)
	if err != nil {
		return nil, fmt.Errorf("while making a request for %q: %v", url, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("something went wrong while making POST request to %v to create %q command", url, command)
	}

	return &cmd, nil
}

// TODO:
func Put(url, username, password, response string) (*models.Command, error) {
	return nil, nil
}

// TODO:
func Delete(url, username, password string) error {
	req, err := http.NewRequest(http.MethodDelete, url, nil)
	if err != nil {
		return fmt.Errorf("while creating a new request for %q: %v", url, err)
	}
	req.SetBasicAuth(username, password)

	c := &http.Client{}
	resp, err := c.Do(req)
	if err != nil {
		return fmt.Errorf("while making a request for %q: %v", url, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("something went wrong while making DELETE request to %v", url)
	}

	return nil
}
