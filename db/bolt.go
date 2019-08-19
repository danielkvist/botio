package db

import (
	"fmt"
	"io"
	"time"

	"github.com/danielkvist/botio/models"

	bolt "go.etcd.io/bbolt"
)

type Bolt struct {
	db *bolt.DB
}

func (bdb *Bolt) Open(path, col string) error {
	db, err := bolt.Open(path, 0600, &bolt.Options{Timeout: 2 * time.Second})
	if err != nil {
		return fmt.Errorf("while opening to DB on %q: %v", path, err)
	}

	err = db.Update(func(tx *bolt.Tx) error {
		bucket := []byte(col)
		if _, err := tx.CreateBucketIfNotExists(bucket); err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return fmt.Errorf("while initializing collection %q on database %q: %v", col, path, err)
	}

	bdb.db = db
	return nil
}

func (bdb *Bolt) Set(col, el, val string) (*models.Command, error) {
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

func (bdb *Bolt) Get(col, el string) (*models.Command, error) {
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

func (bdb *Bolt) GetAll(col string) ([]*models.Command, error) {
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

func (bdb *Bolt) Remove(col, el string) error {
	return bdb.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(col))
		return b.Delete([]byte(el))
	})
}

func (bdb *Bolt) Update(col, el, val string) (*models.Command, error) {
	return bdb.Set(col, el, val)
}

func (bdb *Bolt) Backup(w io.Writer) (int, error) {
	var length int
	err := bdb.db.View(func(tx *bolt.Tx) error {
		length = int(tx.Size())
		_, err := tx.WriteTo(w)
		return err
	})

	return length, err
}