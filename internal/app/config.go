package app

import (
	"errors"
	"flag"
	"strings"
)

const defaultListenAddr = ":50051"

// ErrNameRequired is returned when the required -name flag is missing.
var ErrNameRequired = errors.New("-name is required")

// Config stores runtime settings parsed from CLI flags.
type Config struct {
	Name       string
	PeerAddr   string
	ListenAddr string
}

// ParseArgs parses CLI flags and returns normalized application config.
func ParseArgs(args []string) (Config, error) {
	fs := flag.NewFlagSet("chat", flag.ContinueOnError)

	var cfg Config
	fs.StringVar(&cfg.Name, "name", "", "username (required)")
	fs.StringVar(&cfg.PeerAddr, "peer", "", "peer address host:port to connect to (optional)")
	fs.StringVar(&cfg.ListenAddr, "listen", defaultListenAddr, "listen address host:port")

	if err := fs.Parse(args); err != nil {
		return Config{}, err
	}

	cfg.Name = strings.TrimSpace(cfg.Name)

	if cfg.Name == "" {
		return Config{}, ErrNameRequired
	}

	return cfg, nil
}

// Mode reports whether the app should wait for an incoming session or dial a peer.
func (c Config) Mode() string {
	if c.PeerAddr == "" {
		return "server"
	}

	return "client"
}
