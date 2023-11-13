package icanteen

import (
	"flag"
	"testing"
)

var (
	username    string
	password    string
	server      string
	name        string
	ic_username string
	ic_password string
	ic_server   string
)

func init() {
	flag.StringVar(&username, "username", "", "Edupage username")
	flag.StringVar(&password, "password", "", "Edupage password")
	flag.StringVar(&server, "server", "", "Edupage server")
	flag.StringVar(&name, "name", "", "Name of the user (firstname lastname)")
	flag.StringVar(&ic_username, "ic_username", "", "iCanteen username")
	flag.StringVar(&ic_password, "ic_password", "", "iCanteen password")
	flag.StringVar(&ic_server, "ic_server", "", "iCanteen server")
}

func TestGetLunchesMissingHttps(t *testing.T) {
	lunches, err := LoadLunches(ic_username, ic_password, ic_server)
	if err != nil {
		t.Fatalf("Failed to get lunches: %s", err)
	}

	if len(lunches) == 0 {
		t.Fatalf("No lunches found")
	}
}

func TestGetLunchesTrailingLogin(t *testing.T) {
	lunches, err := LoadLunches(ic_username, ic_password, ic_server+"/login")
	if err != nil {
		t.Fatalf("Failed to get lunches: %s", err)
	}

	if len(lunches) == 0 {
		t.Fatalf("No lunches found")
	}
}

func TestGetLunchesTrailingSlash(t *testing.T) {
	lunches, err := LoadLunches(ic_username, ic_password, ic_server+"/")
	if err != nil {
		t.Fatalf("Failed to get lunches: %s", err)
	}

	if len(lunches) == 0 {
		t.Fatalf("No lunches found")
	}
}

func TestGetLunchesInvalidURL(t *testing.T) {
	_, err := LoadLunches(ic_username, ic_password, "https:////there@shoudl..be@@an//error")
	if err == nil {
		t.Fatal("Expected error, got nil")
	}
}
