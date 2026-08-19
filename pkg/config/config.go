package config

import (
	"os"
)

const defaultAddr = ":8080"

type Config struct {
	Addr string
}

func Load() Config {
	addr := os.Getenv("ADDR")
	if addr == "" {
		addr = defaultAddr
	}
	return Config{Addr: addr}
}
