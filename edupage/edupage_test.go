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
		t.Log("Password parameter missing, (-password=?)")
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
		t.Log("Password parameter missing, (-password=?)")
		return
	}

	if len(server) == 0 {
		t.Log("Server parameter missing, (-server=?)")
		return
	}

	e, err := Login(server, username, password)
	if err != nil {
		t.Error(err)
		return
	}

	err = e.Fetch()

	if err != nil {
		t.Error(err)
		return
	}

	if len(e.EdupageData.Timeline.Items) == 0 {
		t.Error("Recieved timeline items length is zero")
		return
	}

	if len(e.EdupageData.User.UserGroups) == 0 {
		t.Error("Recieved usergroup length is zero")
		return
	}

	if len(e.EdupageData.User.DBI.Teachers) == 0 {
		t.Error("Recieved teacher map length is zero")
		return
	}

}
