package config

import (
	"testing"
)

func TestLoadDefaultAddr(t *testing.T) {
	t.Setenv("ADDR", "")

	cfg := Load()
	if cfg.Addr != defaultAddr {
		t.Fatalf("Addr = %q, want %q", cfg.Addr, defaultAddr)
	}
}

func TestLoadCustomAddr(t *testing.T) {
	t.Setenv("ADDR", ":3000")

	cfg := Load()
	if cfg.Addr != ":3000" {
		t.Fatalf("Addr = %q, want %q", cfg.Addr, ":3000")
	}
}
