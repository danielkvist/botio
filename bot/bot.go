// Package bot exports a Bot interface to manage bots
// for differents platforms easily.
package bot

import "fmt"

// Bot is an interface to manage bots for differentes platforms.
type Bot interface {
	Connect(addr string, cert string, token string, cap int) error
	Start() error
	Listen() error
	Stop() error
}

// Response represents a bot response.
type Response struct {
	id   string
	text string
}

// Create returns a bot that satisfies the Bot interface
// depending on the received platform. If the platform is not supported
// it returns an error.
func Create(platform string) (Bot, error) {
	switch platform {
	case "telegram":
		return &Telegram{}, nil
	case "discord":
		return &Discord{}, nil
	default:
		return nil, fmt.Errorf("platform %q not supported", platform)
	}
}
