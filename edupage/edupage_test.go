package edupage

import (
	"flag"
	"testing"
)

var (
	username string
	password string
	server   string
)

func init() {
	flag.StringVar(&username, "username", "", "Edupage user name")
	flag.StringVar(&password, "password", "", "Edupage user password")
	flag.StringVar(&server, "server", "", "Edupage user server")
}

func TestAutoLogin(t *testing.T) {
	if len(username) == 0 {
		t.Log("Username parameter missing, (-username=?)")
		return
	}

	if len(password) == 0 {
		t.Log("Password parameter missing, (-pasword=?)")
		return
	}

	_, err := LoginAuto(username, password)
	if err != nil {
		t.Error(err)
		return
	}

}

func TestEdupage(t *testing.T) {
	if len(username) == 0 {
		t.Log("Username parameter missing, (-username=?)")
		return
	}

	if len(password) == 0 {
		t.Log("Password parameter missing, (-pasword=?)")
		return
	}

	if len(server) == 0 {
		t.Log("Server parameter missing, (-server=?)")
		return
	}

	h, err := Login(server, username, password)
	if err != nil {
		t.Error(err)
		return
	}

	timeline, err := h.GetTimeline()
	if err != nil {
		t.Error(err)
		return
	}

	if timeline.Raw["module"] != "messages" {
		t.Error("Invalid timeline raw data")
		return
	}
}
