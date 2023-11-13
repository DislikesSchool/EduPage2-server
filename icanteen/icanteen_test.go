package icanteen

import (
	"flag"
	"testing"
)

var (
	server   = flag.String("server", "", "Server URL")
	username = flag.String("username", "", "Username")
	password = flag.String("password", "", "Password")
)

func TestGetLunches(t *testing.T) {
	lunches, err := LoadLunches(*username, *password, *server)
	if err != nil {
		t.Fatalf("Failed to get lunches: %s", err)
	}

	if len(lunches) == 0 {
		t.Fatalf("No lunches found")
	}
}
