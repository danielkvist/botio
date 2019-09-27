// Package db exports a DB interface to manage different databases easily.
package db

import (
	"github.com/danielkvist/botio/models"
)

// DB represents a database with basic CRUD methods.
type DB interface {
	Connect() error
	Add(el, val string) (*models.Command, error)
	Get(el string) (*models.Command, error)
	GetAll() ([]*models.Command, error)
	Remove(el string) error
	Update(el, val string) (*models.Command, error)
	Close() error
}

// Create returns a database that satisfies the DB interface
// depending on the received environment.
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
