// Package cache exports a Cache interface to manage in-memory caches.
package cache

import (
	"github.com/danielkvist/botio/proto"
	"github.com/pkg/errors"

	"github.com/dgraph-io/ristretto"
)

// Cache represents an in-memory cache client.
type Cache interface {
	Add(cmd *proto.BotCommand) error
	Get(cmd *proto.Command) (*proto.BotCommand, error)
	Remove(cmd *proto.Command) error
}

// New receives a set of values that cannot be zero and returns a new
// Cache client or an error if somethign went wrong.
func New(cap int64) (Cache, error) {
	c, err := ristretto.NewCache(&ristretto.Config{
		NumCounters: 1e7,
		MaxCost:     cap,
		BufferItems: 64,
	})

	if err != nil {
		return nil, errors.Wrap(err, "while creating a new Cache")
	}

	cache := &cache{cache: c}
	return cache, nil
}
