package edupage

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"regexp"
	"strconv"
	"time"

	"github.com/DislikesSchool/EduPage2-server/edupage/model"
)

// EdupageClient is used to access the edupage api.
type EdupageClient struct {
	Credentials Credentials
	gsechash    string

	user      *model.User
	timeline  *model.Timeline
	results   *model.Results
	timetable *model.Timetable
}

func CreateClient(credentials Credentials) (EdupageClient, error) {
	var client EdupageClient
	if credentials.httpClient == nil {
		return EdupageClient{}, errors.New("http client in credentials can not be nil")
	}
	client.Credentials = credentials

	user, err := client.loadUser()

	if err != nil {
		return EdupageClient{}, err
	}

	client.user = &user

	return client, nil
}

// GetUser retrieves the user from edupage or returns the stored data.
// If update is set to true, user data wil explicitly update
func (client *EdupageClient) GetUser(update bool) (model.User, error) {
	if client.user == nil || update {
		user, err := client.loadUser()
		if err != nil {
			return model.User{}, err
		}

		client.user = &user
		return user, nil
	} else {
		return *client.user, nil
	}
}

// GetTimeline retrieves last 30 days of timeline from edupage.
func (client *EdupageClient) GetTimeline() (model.Timeline, error) {
	timeline, err := client.loadTimeline(time.Now().AddDate(0, 0, -30), time.Now())
	if err != nil {
		return model.Timeline{}, err
	}

	return timeline, nil
}

// GetUser retrieves the timeline in a specified time interval from edupage.
func (client *EdupageClient) GetTimelineFrom(from, to time.Time) (model.Timeline, error) {
	tt, err := client.loadTimeline(from, to)
	if err != nil {
		return model.Timeline{}, err
	}
	return tt, nil
}

// GetRecentResults retrieves the results from the current year from edupage.
func (client *EdupageClient) GetRecentResults() (model.Results, error) {
	year := time.Now().Format("2006")
	halfyear := "RX" //TODO
	return client.fetchResults(year, halfyear)
}

// GetResults retrieves the results in a specified interval from edupage.
func (client *EdupageClient) GetResults(year, halfyear string) (model.Results, error) {
	results, err := client.fetchResults(year, halfyear)
	if err != nil {
		return model.Results{}, err
	}

	return results, nil
}

// GetResults retrieves this week's timetable from edupage.
func (client *EdupageClient) GetTimetable() (model.Timetable, error) {
	tt, err := client.fetchTimetable(time.Now().AddDate(0, 0, -7), time.Now())
	if err != nil {
		return model.Timetable{}, err
	}

	return tt, nil
}

// LoadTimeline loads the timeline data from the specified date range.
func (client *EdupageClient) loadTimeline(datefrom, dateto time.Time) (model.Timeline, error) {
	url := fmt.Sprintf("https://%s/timeline/?akcia=getData", client.Credentials.Server)

	form, err := CreatePayload(map[string]string{
		"datefrom": datefrom.Format("2006-01-02"),
		"dateto":   dateto.Format("2006-01-02"),
	})

	if err != nil {
		return model.Timeline{}, fmt.Errorf("failed to create payload: %s", err)
	}

	response, err := client.Credentials.httpClient.PostForm(url, form)
	if err != nil {
		return model.Timeline{}, ErrAuthorization // most likely case
	}

	if response.StatusCode != 200 {
		return model.Timeline{}, fmt.Errorf("server returned code:%d", response.StatusCode)
	}

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return model.Timeline{}, fmt.Errorf("failed to read response body: %s", err)
	}

	decoded_body := make([]byte, base64.StdEncoding.DecodedLen(len(body)-4))

	_, err = base64.StdEncoding.Decode(decoded_body, body[4:])
	if err != nil {
		return model.Timeline{}, fmt.Errorf("failed to decode response body: %s", err)
	}

	decoded_body = bytes.Trim(decoded_body, "\x00")
	timeline, err := model.ParseTimeline(decoded_body)
	if err != nil {
		return model.Timeline{}, fmt.Errorf("failed to parse timeline json into json object: %s", err)
	}

	return timeline, nil
}

// LoadUser loads the user data
func (client *EdupageClient) loadUser() (model.User, error) {
	u := fmt.Sprintf("https://%s/user/?", client.Credentials.Server)

	response, err := client.Credentials.httpClient.Get(u)
	if err != nil {
		return model.User{}, ErrAuthorization // most likely case
	}

	if response.StatusCode != 200 {
		return model.User{}, fmt.Errorf("server returned code: %d", response.StatusCode)
	}

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return model.User{}, fmt.Errorf("failed to read response body: %s", err)
	}

	hash, err := findGSCEhash(body)
	if err != nil {
		return model.User{}, fmt.Errorf("failed to parse user json into json object: %s", err)
	}

	client.gsechash = hash

	js, err := findUserHome(body)
	if err != nil {
		return model.User{}, fmt.Errorf("failed to parse user json into json object: %s", err)
	}

	var user model.User
	err = json.Unmarshal([]byte(js), &user)
	if err != nil {
		return model.User{}, fmt.Errorf("failed to parse user json into json object: %s", err)
	}

	return user, nil
}

// LoadResults loads the grade data from specified year and halfyear
// Halfyears types are: P1 (first halfyear), P2 (second halfyear), RX (whole year)
func (client *EdupageClient) fetchResults(year, halfyear string) (model.Results, error) {
	url := fmt.Sprintf("https://%s/znamky/?what=studentviewer&akcia=studentData&eqav=1&maxEqav=7", client.Credentials.Server)

	form, err := CreatePayload(map[string]string{
		"pohlad":           "podladatumu",
		"znamky_yearid":    year,
		"znamky_yearid_ns": "1",
		"nadobdobie":       halfyear,
		"rokobdobie":       fmt.Sprintf("%s::%s", year, halfyear),
		"doRq":             "1",
		"what":             "studentviewer",
		"updateLastView":   "0",
	})

	if err != nil {
		return model.Results{}, fmt.Errorf("failed to create payload: %s", err)
	}

	response, err := client.Credentials.httpClient.PostForm(url, form)
	if err != nil {
		return model.Results{}, ErrAuthorization // most likely case
	}

	if response.StatusCode != 200 {
		return model.Results{}, fmt.Errorf("server returned code:%d", response.StatusCode)
	}

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return model.Results{}, fmt.Errorf("failed to read response body: %s", err)
	}

	decoded_body := make([]byte, base64.StdEncoding.DecodedLen(len(body)-4))

	_, err = base64.StdEncoding.Decode(decoded_body, body[4:])
	if err != nil {
		return model.Results{}, fmt.Errorf("failed to decode response body: %s", err)
	}

	decoded_body = bytes.Trim(decoded_body, "\x00")

	results, err := model.ParseResults(decoded_body)
	if err != nil {
		return model.Results{}, fmt.Errorf("failed to parse results: %s", err)
	}

	return results, nil
}

func (client *EdupageClient) fetchTimetable(datefrom, dateto time.Time) (model.Timetable, error) {
	u := fmt.Sprintf("https://%s/timetable/server/currenttt.js?__func=curentttGetData", client.Credentials.Server)

	id, err := client.GetStudentID()
	if err == ErrorUnitialized {
		return model.Timetable{}, errors.New("failed to create request, user is not initialized")
	}

	year, _ := strconv.Atoi(datefrom.Format("2006"))

	request := map[string]interface{}{
		"__args": []map[string]interface{}{
			nil,
			{
				"year":                 year,
				"datefrom":             datefrom.Format(model.TimeFormatYearMonthDay),
				"dateto":               dateto.Format(model.TimeFormatYearMonthDay),
				"table":                "students",
				"id":                   id,
				"showColors":           false,
				"showOrig":             true,
				"showIgroupsInClasses": false,
				"log_module":           "CurrentTTView",
			},
		},
		"__gsh": client.gsechash,
	}

	request_body, err := json.Marshal(request)
	if err != nil {
		return model.Timetable{}, fmt.Errorf("failed to create request: %s", err)
	}

	response, err := client.Credentials.httpClient.Post(u, "application/json", bytes.NewBuffer(request_body))
	if err != nil {
		return model.Timetable{}, ErrAuthorization // most likely case
	}

	if response.StatusCode != 200 {
		return model.Timetable{}, fmt.Errorf("server returned code: %d", response.StatusCode)
	}

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return model.Timetable{}, fmt.Errorf("failed to read response body: %s", err)
	}

	tt, err := model.ParseTimetable(body)
	if err != nil {
		return model.Timetable{}, fmt.Errorf("failed to read response body: %s", err)
	}

	return tt, nil
}

func findGSCEhash(body []byte) (string, error) {
	rg, _ := regexp.Compile(`ASC\.gsechash="(.*)";`)
	matches := rg.FindAllStringSubmatch(string(body), -1)
	if len(matches) == 0 {
		return "", errors.New("gsechash not found in the document body")
	}

	return matches[0][1], nil
}

func findUserHome(body []byte) (string, error) {
	rg, _ := regexp.Compile(`\.userhome\((.*)\);`)
	matches := rg.FindAllStringSubmatch(string(body), -1)
	if len(matches) == 0 {
		return "", errors.New("userhome not found in the document body")
	}

	return matches[0][1], nil
}
