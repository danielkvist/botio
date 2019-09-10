// Package bot exports a Bot interface to manage bots
// for differents platforms easily.
package bot

// Bot is an interface to manage bots for differentes platforms.
type Bot interface {
	Connect(token string, cap int) error
	Start() error
	Listen(url, key string) error
	Stop() error
}

// Response represents a bot response.
type Response struct {
	id   string
	text string
}

// Create returns a bot that satisfies the Bot interface
// depending on the received platform.
func Create(platform string) Bot {
	switch platform {
	case "telegram":
		return &Telegram{}
	case "discord":
		return &Discord{}
	default:
		return nil
	}
}
