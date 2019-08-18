package tdb

import (
	"fmt"
	"io"

	"github.com/danielkvist/botio/models"
)

type TDB map[string]string

func (db TDB) Set(col string, el string, val string) (*models.Command, error) {
	db[el] = val

	return &models.Command{
		Cmd:      el,
		Response: val,
	}, nil
}

func (db TDB) Get(col string, el string) (*models.Command, error) {
	val, ok := db[el]
	if !ok {
		return nil, fmt.Errorf("element %q not found", el)
	}

	return &models.Command{
		Cmd:      el,
		Response: val,
	}, nil
}

func (db TDB) GetAll(col string) ([]*models.Command, error) {
	var commands []*models.Command

	for k, v := range db {
		tmpCommand := &models.Command{
			Cmd:      k,
			Response: v,
		}

		commands = append(commands, tmpCommand)
	}

	return commands, nil
}

func (db TDB) Remove(col string, el string) error {
	delete(db, el)
	return nil
}

func (db TDB) Update(col string, el string, val string) (*models.Command, error) {
	db[el] = val

	return &models.Command{
		Cmd:      el,
		Response: val,
	}, nil
}

func (db TDB) Backup(w io.Writer) (int, error) {
	return 0, nil
}
