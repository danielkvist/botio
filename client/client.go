// Package client exports very basic functions to get commands with HTTP methods.
package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/danielkvist/botio/models"

	jwt "github.com/dgrijalva/jwt-go"
)

func Get(url, key string) (*models.Command, error) {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("while creating a new request for %q: %v", url, err)
	}

	token, err := genJWT(key)
	if err != nil {
		return nil, fmt.Errorf("while generating JWT token for authentication for %q: %v", url, err)
	}

	req.Header.Set("Token", token)

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

func GetAll(url, key string) ([]*models.Command, error) {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("while creating a new request for %q: %v", url, err)
	}

	token, err := genJWT(key)
	if err != nil {
		return nil, fmt.Errorf("while generating JWT token for authentication for %q: %v", url, err)
	}

	req.Header.Set("Token", token)

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

func Post(url, key, command, response string) (*models.Command, error) {
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

	token, err := genJWT(key)
	if err != nil {
		return nil, fmt.Errorf("while generating JWT token for authentication for %q: %v", url, err)
	}

	req.Header.Set("Token", token)

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

func Put(url, key, command, response string) (*models.Command, error) {
	return Post(url, key, command, response)
}

func Delete(url, key string) error {
	req, err := http.NewRequest(http.MethodDelete, url, nil)
	if err != nil {
		return fmt.Errorf("while creating a new request for %q: %v", url, err)
	}

	token, err := genJWT(key)
	if err != nil {
		return fmt.Errorf("while generating JWT token for authentication for %q: %v", url, err)
	}

	req.Header.Set("Token", token)

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

func genJWT(key string) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)

	tkStr, err := token.SignedString([]byte(key))
	if err != nil {
		return "", fmt.Errorf("while generating the authentication JWT token: %v", err)
	}

	return tkStr, nil
}
