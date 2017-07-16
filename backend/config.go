package backend

import "time"

// Config saves important game configurations.
type Config struct {
	PlayerLimit int
	Timeout     time.Duration
}

// DefaultConfig returns the default config.
func DefaultConfig() Config {
	return Config{
		PlayerLimit: 4,
		Timeout:     60 * time.Second,
	}
}
