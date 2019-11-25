// Package cache exports a Cache interface to manage in-memory caches.
package cache

import (
	"github.com/danielkvist/botio/proto"
)

// Cache represents a cache with basic methods to manage
// the items in the cache itself.
type Cache interface {
	Init(cap int) error
	Add(cmd *proto.BotCommand) error
	Get(cmd *proto.Command) (*proto.BotCommand, error)
	Remove(cmd *proto.Command) error
}

// Create follows the Factory patterns to return a Cache
// system depending on the received platform parameter.
func Create(platform string) Cache {
	switch platform {
	case "ristretto":
		return &ristrettoCache{}
	default:
		return nil
	}
}
