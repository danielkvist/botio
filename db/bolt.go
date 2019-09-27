package db

import (
	"fmt"
	"time"

	"github.com/danielkvist/botio/models"

	bolt "go.etcd.io/bbolt"
)

// Bolt wraps a bolt.DB database and satifies the DB interface.
type Bolt struct {
	Path string
	Col  string
	db   *bolt.DB
}

func (bdb *Bolt) Connect() error {
	db, err := bolt.Open(bdb.Path, 0600, &bolt.Options{Timeout: 2 * time.Second})
	if err != nil {
		return fmt.Errorf("while opening DB on %q: %v", bdb.Path, err)
	}

	err = db.Update(func(tx *bolt.Tx) error {
		bucket := []byte(bdb.Col)
		if _, err := tx.CreateBucketIfNotExists(bucket); err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return fmt.Errorf("while initializing collection %q on DB %q: %v", bdb.Col, bdb.Path, err)
	}

	bdb.db = db
	return nil
}

func (bdb *Bolt) Add(el, val string) (*models.Command, error) {
	err := bdb.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bdb.Col))
		return b.Put([]byte(el), []byte(val))
	})

	if err != nil {
		return nil, fmt.Errorf("while adding command %q: %v", el, err)
	}

	command := &models.Command{
		Cmd:      el,
		Response: val,
	}

	return command, nil
}

func (bdb *Bolt) Get(el string) (*models.Command, error) {
	var val []byte
	err := bdb.db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(bdb.Col))
		val = bucket.Get([]byte(el))

		if len(val) == 0 {
			return fmt.Errorf("element %q not found", el)
		}

		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("while getting command %q: %v", el, err)
	}

	command := &models.Command{
		Cmd:      el,
		Response: string(val),
	}

	return command, nil
}

func (bdb *Bolt) GetAll() ([]*models.Command, error) {
	var commands []*models.Command
	err := bdb.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bdb.Col))
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

func (bdb *Bolt) Remove(el string) error {
	return bdb.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bdb.Col))
		return b.Delete([]byte(el))
	})
}

func (bdb *Bolt) Update(el, val string) (*models.Command, error) {
	command, err := bdb.Add(el, val)
	if err != nil {
		return nil, fmt.Errorf("while updating command %q: %v", el, err)
	}

	return command, nil
}

func (bdb *Bolt) Close() error {
	if err := bdb.db.Close(); err != nil {
		return fmt.Errorf("while closing DB: %v", err)
	}

	return nil
}
