// Package cache exports a Cache interface to manage in-memory caches.
package cache

import (
	"github.com/danielkvist/botio/proto"

	"github.com/dgraph-io/ristretto"
	"github.com/pkg/errors"
)

// Cache represents an in-memory cache client.
type Cache interface {
	Add(cmd *proto.BotCommand) error
	Get(cmd *proto.Command) (*proto.BotCommand, error)
	Remove(cmd *proto.Command) error
}

type cache struct {
	cache *ristretto.Cache
}

// New receives a set of values that cannot be zero and returns a new
// Cache client or an error if somethign went wrong.
func New(counters, cost, bufferItems int64) (Cache, error) {
	c, err := ristretto.NewCache(&ristretto.Config{
		NumCounters: counters,
		MaxCost:     cost,
		BufferItems: bufferItems,
	})

	if err != nil {
		return nil, errors.Wrap(err, "while creating a new Cache")
	}

	cache := &cache{cache: c}
	return cache, nil
}

// Add adds to the cache a new *proto.BotCommand. It returns a non-nill error
// if the received *proto.BotCommand has a Command or a Response empty or if
// something went wrong while adding the command to the cache itself.
func (c *cache) Add(cmd *proto.BotCommand) error {
	command := cmd.GetCmd().GetCommand()
	resp := cmd.GetResp().GetResponse()

	switch {
	case command == "":
		return errors.Errorf("command cannot be an empty string")
	case resp == "":
		return errors.Errorf("command's response cannot be an empty string")
	}

	ok := c.cache.Set(command, resp, 1)
	if !ok {
		return errors.Errorf("error while adding command %q with response %q to cache", command, resp)
	}

	return nil
}

// Get receives a *proto.Command and returns the respective *proto.BotCommand
// if exists. It returns a non-nil error if the command was not found of if
// there is any error while getting it.
func (c *cache) Get(cmd *proto.Command) (*proto.BotCommand, error) {
	el := cmd.GetCommand()

	val, ok := c.cache.Get(el)
	if !ok {
		return nil, errors.Errorf("command %q not found on cache", el)
	}

	resp, ok := val.(string)
	if !ok {
		return nil, errors.Errorf("while converting received value as response to command %q from cache to string", el)
	}

	return &proto.BotCommand{
		Cmd: &proto.Command{
			Command: el,
		},
		Resp: &proto.Response{
			Response: resp,
		},
	}, nil
}

// Remove deletes a *proto.BotCommand from the cache. It never
// returns a non-nil error.
func (c *cache) Remove(cmd *proto.Command) error {
	c.cache.Del(cmd.GetCommand())
	return nil
}
