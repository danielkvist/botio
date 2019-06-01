package db

import (
	"fmt"
	"time"

	"github.com/danielkvist/botio/models"

	bolt "go.etcd.io/bbolt"
)

type DB struct {
	db *bolt.DB
}

func Open(path string) (*DB, error) {
	db, err := bolt.Open(path, 0600, &bolt.Options{Timeout: 2 * time.Second})
	if err != nil {
		return nil, fmt.Errorf("while opening to DB on %q: %v", path, err)
	}

	return &DB{db}, nil
}

func (db *DB) Set(collection string, element string, value string) (*models.Command, error) {
	err := db.db.Update(func(tx *bolt.Tx) error {
		bucket := []byte(collection)
		if _, err := tx.CreateBucketIfNotExists(bucket); err != nil {
			return err
		}

		b := tx.Bucket(bucket)
		return b.Put([]byte(element), []byte(value))
	})

	if err != nil {
		return nil, err
	}

	command := &models.Command{
		Cmd:      element,
		Response: value,
	}

	return command, nil
}

func (db *DB) Get(collection string, element string) (*models.Command, error) {
	var val []byte
	err := db.db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(collection))
		val = bucket.Get([]byte(element))
		return nil
	})

	if err != nil {
		return nil, err
	}

	command := &models.Command{
		Cmd:      element,
		Response: string(val),
	}

	return command, nil
}

func (db *DB) GetAll(collection string) ([]*models.Command, error) {
	var commands []*models.Command
	err := db.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(collection))
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

func (db *DB) Remove(collection string, element string) error {
	return db.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(collection))
		return b.Delete([]byte(element))
	})
}

func (db *DB) Update(collection string, element string, value string) (*models.Command, error) {
	return db.Set(collection, element, value)
}
