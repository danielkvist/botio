package db

import (
	"fmt"
	"time"

	"github.com/danielkvist/botio/proto"
	bolt "go.etcd.io/bbolt"
)

// Bolt wraps a bolt.DB client and
// satisfies the DB interface.
type Bolt struct {
	Path string
	Col  string
	db   *bolt.DB
}

// Connect treis to connect to a BoltDB database. If it fails
// it returns a non-nil error.
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

// Add receives a *proto.BotCommand and adds it to the bucket
// designated. If something goes wrong it returns a non-nil error.
func (bdb *Bolt) Add(cmd *proto.BotCommand) error {
	el := cmd.GetCmd().GetCommand()
	val := cmd.GetResp().GetResponse()
	err := bdb.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bdb.Col))
		return b.Put([]byte(el), []byte(val))
	})

	if err != nil {
		return fmt.Errorf("while adding command %q: %v", el, err)
	}

	return nil
}

// Get receives a *proto.Command and returns the respective *proto.BotCommand
// if exists in the designated bucket. If not it returns a non-nil error.
func (bdb *Bolt) Get(cmd *proto.Command) (*proto.BotCommand, error) {
	el := cmd.GetCommand()

	var val []byte
	err := bdb.db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(bdb.Col))
		val = bucket.Get([]byte(el))

		if len(val) == 0 {
			return fmt.Errorf("command %q not found", el)
		}

		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("while getting command %q: %v", el, err)
	}

	return &proto.BotCommand{
		Cmd: &proto.Command{
			Command: el,
		},
		Resp: &proto.Response{
			Response: string(val),
		},
	}, nil
}

// GetAll ranges over all the entries of the designated bucket
// and returns a *proto.BotCommands with all the *proto.BotCommand
// found. If something goes wrong it returns a non-nil error.
func (bdb *Bolt) GetAll() (*proto.BotCommands, error) {
	var commands []*proto.BotCommand

	err := bdb.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bdb.Col))
		c := b.Cursor()

		for k, v := c.First(); k != nil; k, v = c.Next() {
			commands = append(commands, &proto.BotCommand{
				Cmd: &proto.Command{
					Command: string(k),
				},
				Resp: &proto.Response{
					Response: string(v),
				},
			})
		}

		return nil
	})

	return &proto.BotCommands{
		Commands: commands,
	}, err
}

// Remove removes a *proto.BotCommand from the designated bucket.
// It returns a non-nil error if something goes wrong.
func (bdb *Bolt) Remove(cmd *proto.Command) error {
	el := cmd.GetCommand()

	return bdb.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bdb.Col))
		err := b.Delete([]byte(el))

		if err != nil {
			return fmt.Errorf("while removing command %q: %v", el, err)
		}

		return nil
	})
}

// Update updates the Response of an existing *proto.BotCommand
// with he Response of the received *proto.BotCommand. If the
// *proto.BotCommand didn't exists it adds it to the bucket
// due to how BoltDB databases work. If something goes wrong
// it returns a non-nil error.
func (bdb *Bolt) Update(cmd *proto.BotCommand) error {
	if err := bdb.Add(cmd); err != nil {
		return fmt.Errorf("while updating command %q: %v", cmd.GetCmd().GetCommand(), err)
	}

	return nil
}

// Close tries to close the connection to the BoltDB database.
// If fails it returns a non-nil error.
func (bdb *Bolt) Close() error {
	if err := bdb.db.Close(); err != nil {
		return fmt.Errorf("while closing connection to DB: %v", err)
	}

	return nil
}
