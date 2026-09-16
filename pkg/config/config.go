package config

import (
	"os"
)

const defaultAddr = ":8080"

type Config struct {
	Addr string

	// AdminToken authorizes requests to the admin endpoint. It is read from
	// the environment rather than compiled in. When empty, the admin
	// endpoint is not registered.
	AdminToken string
}

func Load() Config {
	addr := os.Getenv("ADDR")
	if addr == "" {
		addr = defaultAddr
	}
	return Config{
		Addr:       addr,
		AdminToken: os.Getenv("ADMIN_TOKEN"),
	}
}
