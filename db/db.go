// Package db exports a DB interface that is implemented by multiple
// databases clients.
package db

import (
	"github.com/danielkvist/botio/proto"
)

// DB represents a database client with basic CRUD methods
// as basic methods to connect and disconnect from the
// database itself.
type DB interface {
	Connect() error
	Add(cmd *proto.BotCommand) error
	Get(cmd *proto.Command) (*proto.BotCommand, error)
	GetAll() (*proto.BotCommands, error)
	Remove(cmd *proto.Command) error
	Update(cmd *proto.BotCommand) error
	Close() error
}

// Create follows the Factory pattern to return a DB
// depending on the received parameter.
func Create(env string) DB {
	switch env {
	case "local":
		return &Bolt{}
	case "postgres":
		return &Postgres{}
	case "testing":
		var m Mem
		m = make(map[string]string)
		return m
	default:
		return nil
	}
}
