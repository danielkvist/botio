// Package models exports a very simple Command struct.
package models

// Command represents a bot command (Cmd)
// and its response (Response).
type Command struct {
	Cmd      string `json:"cmd"`
	Response string `json:"response"`
}
