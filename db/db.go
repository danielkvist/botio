// Package db exports a simple bolt.DB wrapper to manage bot commands.
package db

import (
	"fmt"
	"io"
	"time"

	"github.com/danielkvist/botio/models"

	bolt "go.etcd.io/bbolt"
)

// Bolter is a simple interface for BDB with testing purposes.
type Bolter interface {
	Set(col string, el string, val string) (*models.Command, error)
	Get(col string, el string) (*models.Command, error)
	GetAll(col string) ([]*models.Command, error)
	Remove(col string, el string) error
	Update(col string, el string, val string) (*models.Command, error)
	Backup(w io.Writer) (int, error)
}

// BDB is a simple wrapper around a bolt ddatabase.
type BDB struct {
	db *bolt.DB
}

// Connect tries to get the database on the received path and create the specified
// collection (bucket in the case of a bolt database) if it does not already exist.
// If all the previous operations have succeeded it returns the database as a *BDB.
// If something goes wrong while opening the database o while creating the collection
// it returns the error.
func Connect(path string, col string) (*BDB, error) {
	db, err := bolt.Open(path, 0600, &bolt.Options{Timeout: 2 * time.Second})
	if err != nil {
		return nil, fmt.Errorf("while opening to DB on %q: %v", path, err)
	}

	err = db.Update(func(tx *bolt.Tx) error {
		bucket := []byte(col)
		if _, err := tx.CreateBucketIfNotExists(bucket); err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("while initializing collection %q on database %q: %v", col, path, err)
	}

	return &BDB{db}, nil
}

// Set tries to add to the specified collection an element with a value.
// It returns the added Command if the operation was successful
// or an error if something went wrong.
func (bdb *BDB) Set(col string, el string, val string) (*models.Command, error) {
	err := bdb.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(col))
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

// Get returns the element specified on a collection.
// If the element doesn't exist, it returns an error.
func (bdb *BDB) Get(col string, el string) (*models.Command, error) {
	var val []byte
	err := bdb.db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(col))
		val = bucket.Get([]byte(el))

		if len(val) == 0 {
			return fmt.Errorf("element %q not found", el)
		}

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

// GetAll returns all the elements found on the specified collection.
// It never returns an error.
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

// Remove tries to delete an specified element of a collection.
// It returns an error if there is any problem while deleting the element.
func (bdb *BDB) Remove(col string, el string) error {
	return bdb.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(col))
		return b.Delete([]byte(el))
	})
}

// Update is basically a wrapper for the Set method since
// updating is the same as posting the same element in a bolt database.
func (bdb *BDB) Update(col string, el string, val string) (*models.Command, error) {
	return bdb.Set(col, el, val)
}

// Backup writes to the received io.Writer the bolt database.
// It returns the size in bytes of the database or an error if there is
// any problem.
func (bdb *BDB) Backup(w io.Writer) (int, error) {
	var length int
	err := bdb.db.View(func(tx *bolt.Tx) error {
		length = int(tx.Size())
		_, err := tx.WriteTo(w)
		return err
	})

	return length, err
}
