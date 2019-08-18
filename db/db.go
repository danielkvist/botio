package db

import (
	"io"

	"github.com/danielkvist/botio/models"
)

type DB interface {
	Open(path, col string) error
	Set(col, el, val string) (*models.Command, error)
	Get(col, el string) (*models.Command, error)
	GetAll(col string) ([]*models.Command, error)
	Remove(col, el string) error
	Update(col, el string, val string) (*models.Command, error)
	Backup(w io.Writer) (int, error)
}

func DBFactory(env string) DB {
	switch env {
	case "production":
		return &Bolt{}
	case "testing":
		var m Mem
		m = make(map[string]string)
		return m
	default:
		return nil
	}
}
