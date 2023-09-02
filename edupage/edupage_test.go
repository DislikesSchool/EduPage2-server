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
		return errors.New("server parameter missing, (-server=?)")
	}
	return nil
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
	err := checkCredentials()
	if err != nil {
		t.Error(err)
		return
	}

	client, err := Login(server, username, password)
	if err != nil {
		t.Error(err)
		return
	}

	err = client.Fetch()

	if err != nil {
		t.Error(err)
		return
	}

	if len(client.Timeline.Items) == 0 {
		t.Error("Recieved timeline array is empty")
	}

	if len(client.User.UserGroups) == 0 {
		t.Error("Recieved usergroup array is empty")
	}

	if len(client.User.DBI.Teachers) == 0 {
		t.Error("Recieved teacher map is empty")
	}

	if len(client.Results.Grades) == 0 {
		t.Error("Recieved grade array is empty")
	}
}

func BenchmarkLogin(t *testing.B) {
	err := checkCredentials()
	if err != nil {
		t.Error(err)
		return
	}
	t.ResetTimer()

	_, err = Login(server, username, password)
	if err != nil {
		t.Error(err)
		return
	}
}

func BenchmarkTimeline(t *testing.B) {
	err := checkCredentials()
	if err != nil {
		t.Error(err)
		return
	}

	client, err := Login(server, username, password)
	if err != nil {
		t.Error(err)
		return
	}

	t.ResetTimer()

	err = client.LoadRecentTimeline()
	if err != nil {
		t.Error(err)
		return
	}

	t.StopTimer()

	if len(client.Timeline.Items) == 0 {
		t.Error("Recieved timeline array is empty")
	}
}

func BenchmarkUser(t *testing.B) {
	err := checkCredentials()
	if err != nil {
		t.Error(err)
		return
	}

	client, err := Login(server, username, password)
	if err != nil {
		t.Error(err)
		return
	}

	t.ResetTimer()

	err = client.LoadUser()
	if err != nil {
		t.Error(err)
		return
	}

	t.StopTimer()

	if len(client.User.UserGroups) == 0 {
		t.Error("Recieved user group array is empty")
	}
}

func BenchmarkGrades(t *testing.B) {
	err := checkCredentials()
	if err != nil {
		t.Error(err)
		return
	}

	client, err := Login(server, username, password)
	if err != nil {
		t.Error(err)
		return
	}

	t.ResetTimer()

	err = client.LoadResults("2022", "RX")
	if err != nil {
		t.Error(err)
		return
	}

	t.StopTimer()

	if len(client.Results.Grades) == 0 {
		t.Error("Recieved grades array is empty")
	}
}
