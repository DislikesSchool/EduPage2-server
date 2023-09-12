package edupage

import (
	"errors"
	"flag"
	"fmt"
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

	if len(username) == 0 {
		t.Error(errors.New("username parameter missing, (-username=?)"))
	}

	if len(password) == 0 {
		t.Error(errors.New("password parameter missing, (-password=?)"))
	}

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
		//return
	}

	timeline, err := client.GetRecentTimeline()
	if err != nil {
		t.Error(fmt.Errorf("failed to recieve timeline: %s", err))
		//return
	}

	if len(timeline.Items) == 0 {
		t.Log("Recieved timeline is empty")
	}

	results, err := client.GetRecentResults()
	if err != nil {
		t.Error(fmt.Errorf("failed to recieve results: %s", err))
		//return
	}

	if len(results.Grades) == 0 {
		t.Log("Recieved grades are empty")
	}

	timetable, err := client.GetRecentTimetable()
	if err != nil {
		t.Error(fmt.Errorf("failed to recieve timetable: %s", err))
		//return
	}

	if len(timetable.Days) == 0 {
		t.Log("Recieved timetable is empty")
	}
	for k := range timetable.Days {
		println(k)
	}

	for _, v := range timetable.Days["2023-09-04"] {
		subject, _ := client.GetSubjectByID(v.SubjectID)
		println(subject.Name)
	}

}

func BenchmarkLogin(t *testing.B) {
	err := checkCredentials()
	if err != nil {
		t.Error(err)
		return
	}
	t.ResetTimer()

	_, err = Login(username, password, server)

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

	timeline, err := client.GetRecentTimeline()
	if err != nil {
		t.Error(fmt.Errorf("failed to recieve timeline: %s", err))
		return
	}

	if len(timeline.Items) == 0 {
		t.Log("Recieved timeline is empty")
	}
}

func BenchmarkResults(t *testing.B) {
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

	results, err := client.GetRecentResults()
	if err != nil {
		t.Error(fmt.Errorf("failed to recieve results: %s", err))
		return
	}

	if len(results.Grades) == 0 {
		t.Log("Recieved grades are empty")
	}
}

func BenchmarkTimetable(t *testing.B) {
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

	timetable, err := client.GetRecentTimetable()
	if err != nil {
		t.Error(fmt.Errorf("failed to recieve timetable: %s", err))
		return
	}

	if len(timetable.Days) == 0 {
		t.Log("Recieved timetable is empty")
	}
}
