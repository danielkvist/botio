package cache

import (
	"github.com/danielkvist/botio/proto"

	"github.com/dgraph-io/ristretto"
	"github.com/pkg/errors"
)

type ristrettoCache struct {
	cache *ristretto.Cache
}

// Init initializes a Cache based on ristretto with the
// received capacity.
func (r *ristrettoCache) Init(cap int) error {
	c, err := ristretto.NewCache(&ristretto.Config{
		NumCounters: 1e7,
		MaxCost:     int64(cap),
		BufferItems: 64,
	})
	if err != nil {
		return errors.Wrap(err, "while creating a new Cache based on ristretto")
	}

	r.cache = c
	return nil
}

// Add adds to the cache a new *proto.BotCommand. It returns a non-nill error
// if the received *proto.BotCommand has a Command or a Response empty or if
// something went wrong while adding the command to the cache itself.
func (r *ristrettoCache) Add(cmd *proto.BotCommand) error {
	command := cmd.GetCmd().GetCommand()
	resp := cmd.GetResp().GetResponse()

	switch {
	case command == "":
		return errors.Errorf("command cannot be an empty string")
	case resp == "":
		return errors.Errorf("command's response cannot be an empty string")
	}

	ok := r.cache.Set(command, resp, 1)
	if !ok {
		return errors.Errorf("error while adding command %q with response %q to cache", command, resp)
	}

	return nil
}

// Get receives a *proto.Command and returns the respective *proto.BotCommand
// if exists. It returns a non-nil error if the command was not found of if
// there is any error while getting it.
func (r *ristrettoCache) Get(cmd *proto.Command) (*proto.BotCommand, error) {
	el := cmd.GetCommand()

	val, ok := r.cache.Get(el)
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
func (r *ristrettoCache) Remove(cmd *proto.Command) error {
	r.cache.Del(cmd.GetCommand())
	return nil
}
