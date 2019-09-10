package db

import (
	"fmt"
	"io"

	"github.com/danielkvist/botio/models"
)

// Mem is a simple in-memory map for testing that satisfies the DB interface.
type Mem map[string]string

func (m Mem) Open(path, col string) error {
	return nil
}

func (m Mem) Set(el, val string) (*models.Command, error) {
	m[el] = val

	return &models.Command{
		Cmd:      el,
		Response: val,
	}, nil
}

func (m Mem) Get(el string) (*models.Command, error) {
	val, ok := m[el]
	if !ok {
		return nil, fmt.Errorf("element %q not found", el)
	}

	return &models.Command{
		Cmd:      el,
		Response: val,
	}, nil
}

func (m Mem) GetAll() ([]*models.Command, error) {
	var commands []*models.Command

	for k, v := range m {
		c := &models.Command{
			Cmd:      k,
			Response: v,
		}

		commands = append(commands, c)
	}

	return commands, nil
}

func (m Mem) Remove(el string) error {
	delete(m, el)
	return nil
}

func (m Mem) Update(el, val string) (*models.Command, error) {
	m[el] = val

	return &models.Command{
		Cmd:      el,
		Response: val,
	}, nil
}

func (m Mem) Backup(w io.Writer) (int, error) {
	return 0, nil
}
