package app

import (
	"errors"
	"flag"
	"strings"
)

const defaultListenAddr = ":50051"

var ErrNameRequired = errors.New("-name is required")

type Config struct {
	Name       string
	PeerAddr   string
	ListenAddr string
}

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

func (c Config) Mode() string {
	if c.PeerAddr == "" {
		return "server"
	}

	return "client"
}
