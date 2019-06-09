package db

import (
	"fmt"
	"time"

	"github.com/danielkvist/botio/models"

	bolt "go.etcd.io/bbolt"
)

type Bolter interface {
	Set(col string, el string, val string) (*models.Command, error)
	Get(col string, el string) (*models.Command, error)
	GetAll(col string) ([]*models.Command, error)
	Remove(col string, el string) error
	Update(col string, el string, val string) (*models.Command, error)
}

type BDB struct {
	db *bolt.DB
}

func Open(path string) (*BDB, error) {
	db, err := bolt.Open(path, 0600, &bolt.Options{Timeout: 2 * time.Second})
	if err != nil {
		return nil, fmt.Errorf("while opening to DB on %q: %v", path, err)
	}

	return &BDB{db}, nil
}

func (bdb *BDB) Set(col string, el string, val string) (*models.Command, error) {
	err := bdb.db.Update(func(tx *bolt.Tx) error {
		bucket := []byte(col)
		if _, err := tx.CreateBucketIfNotExists(bucket); err != nil {
			return err
		}

		b := tx.Bucket(bucket)
		return b.Put([]byte(el), []byte(val))
	})

	if err != nil {
		return nil, err
	}

	command := &models.Command{
		Cmd:      el,
		Response: val,
	}

	return command, nil
}

func (bdb *BDB) Get(col string, el string) (*models.Command, error) {
	var val []byte
	err := bdb.db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(col))
		val = bucket.Get([]byte(el))
		return nil
	})

	if err != nil {
		return nil, err
	}

	command := &models.Command{
		Cmd:      el,
		Response: string(val),
	}

	return command, nil
}

func (bdb *BDB) GetAll(col string) ([]*models.Command, error) {
	var commands []*models.Command
	err := bdb.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(col))
		c := b.Cursor()

		for k, v := c.First(); k != nil; k, v = c.Next() {
			commands = append(commands, &models.Command{
				Cmd:      string(k),
				Response: string(v),
			})
		}

		return nil
	})

	return commands, err
}

func (bdb *BDB) Remove(col string, el string) error {
	return bdb.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(col))
		return b.Delete([]byte(el))
	})
}

func (bdb *BDB) Update(col string, el string, val string) (*models.Command, error) {
	return bdb.Set(col, el, val)
}
