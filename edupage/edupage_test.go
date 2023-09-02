package edupage

import (
	"errors"
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

func checkCredentials() error {
	if len(username) == 0 {
		return errors.New("username parameter missing, (-username=?)")
	}

	if len(password) == 0 {
		return errors.New("password parameter missing, (-password=?)")
	}

	if len(server) == 0 {
		return errors.New("password parameter missing, (-server=?)")
	}
	return nil
}

func TestLoginAuto(t *testing.T) {
	credentials, err := LoginAuto(username, password)
	if err != nil {
		t.Error(err)
		return
	}

	_, err = CreateClient(credentials)
	if err != nil {
		t.Error(err)
		return
	}
}

func TestEdupage(t *testing.T) {
	err := checkCredentials()
	if err != nil {
		t.Error(err)
		return
	}

	credentials, err := Login(username, password, server)
	if err != nil {
		t.Error(err)
		return
	}

	_, err = CreateClient(credentials)
	if err != nil {
		t.Error(err)
		return
	}

}

func BenchmarkLogin(t *testing.B) {
	err := checkCredentials()
	if err != nil {
		t.Error(err)
		return
	}
	t.ResetTimer()

	_, err = Login(username, server, password)

	if err != nil {
		t.Error(err)
		return
	}
}
