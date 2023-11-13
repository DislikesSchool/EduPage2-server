package edupage

import (
	"errors"
	"flag"
	"fmt"
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

	canteen, err := client.GetRecentCanteen()
	if err != nil {
		t.Error(fmt.Errorf("failed to recieve canteen: %s", err))
		//return
	}

	if len(canteen.Days) == 0 {
		t.Log("Recieved canteen is empty")
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
