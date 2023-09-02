package edupage

import (
	"errors"
	"flag"
	"testing"
	"time"
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

	client, err := CreateClient(credentials)
	if err != nil {
		t.Error(err)
		return
	}

	if len(client.User.UserGroups) == 0 {
		t.Log("Recieved usergroup array is empty")
	}

	if len(client.User.DBI.Teachers) == 0 {
		t.Log("Recieved teacher map is empty")
	}

	if err := client.LoadRecentTimeline(); err != nil {
		t.Error(err)
		return
	}

	if len(client.Timeline.Items) == 0 {
		t.Log("Recieved timeline array is empty")
	}

	if err := client.LoadResults(time.Now().Format("2006"), "RX"); err != nil {
		t.Error(err)
		return
	}

	if len(client.Results.Grades) == 0 {
		t.Log("Recieved grade array is empty")
	}

	err = client.LoadTimetable(time.Now().Local().AddDate(0, 0, 3), time.Now().Local().AddDate(0, 0, 3))
	if err != nil {
		t.Error(err)
		return
	}

	println(len(client.Timetable.Days))
	for _, v := range client.Timetable.Days {
		for _, sobj := range v {
			subject, err := client.GetSubjectByID(sobj.SubjectID)
			if err != nil {
				println(err.Error())
			}
			println(subject.Name)
		}
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

func BenchmarkTimeline(t *testing.B) {
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

	client, err := CreateClient(credentials)
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

	credentials, err := Login(username, password, server)
	if err != nil {
		t.Error(err)
		return
	}

	client, err := CreateClient(credentials)
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

	credentials, err := Login(username, password, server)
	if err != nil {
		t.Error(err)
		return
	}

	client, err := CreateClient(credentials)
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
