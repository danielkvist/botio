// Package models exports a command struct for bot's commands.
package models

// Command represents a bot command and his response
type Command struct {
	Cmd      string `json:"cmd"`
	Response string `json:"response"`
}
