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

	User      *model.User
	Timeline  *model.Timeline
	Results   *model.Results
	Timetable *model.Timetable
}

func CreateClient(credentials Credentials) (EdupageClient, error) {
	var client EdupageClient
	if credentials.httpClient == nil {
		return EdupageClient{}, errors.New("http client in credentials can not be nil")
	}
	client.Credentials = credentials

	if err := client.LoadUser(); err != nil {
		return EdupageClient{}, err
	}

	return client, nil
}

// LoadRecentTimeline loads the recent timeline data.
// That's from today, to 30 days in the past.
// Also updates the Timeline property in Edupage struct.
func (client *EdupageClient) LoadRecentTimeline() error {
	duration, err := time.ParseDuration("-720h") // 30 days
	if err != nil {
		return fmt.Errorf("failed to parse duration: %s", err)
	}

	start := time.Now().Add(duration)
	return client.LoadTimeline(start, time.Now())
}

// LoadTimeline loads the timeline data from the specified date range.
func (client *EdupageClient) LoadTimeline(datefrom, dateto time.Time) error {
	url := fmt.Sprintf("https://%s/timeline/?akcia=getData", client.Credentials.Server)

	form, err := CreatePayload(map[string]string{
		"datefrom": datefrom.Format("2006-01-02"),
		"dateto":   dateto.Format("2006-01-02"),
	})

	if err != nil {
		return fmt.Errorf("failed to create payload: %s", err)
	}

	response, err := client.Credentials.httpClient.PostForm(url, form)
	if err != nil {
		return ErrAuthorization // most likely case
	}

	if response.StatusCode != 200 {
		return fmt.Errorf("server returned code:%d", response.StatusCode)
	}

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body: %s", err)
	}

	decoded_body := make([]byte, base64.StdEncoding.DecodedLen(len(body)-4))

	_, err = base64.StdEncoding.Decode(decoded_body, body[4:])
	if err != nil {
		return fmt.Errorf("failed to decode response body: %s", err)
	}

	decoded_body = bytes.Trim(decoded_body, "\x00")
	timeline, err := model.ParseTimeline(decoded_body)
	if err != nil {
		return fmt.Errorf("failed to parse timeline json into json object: %s", err)
	}

	if client.Timeline == nil {
		client.Timeline = &timeline
	} else {
		client.Timeline.Merge(&timeline)
	}

	return nil
}

// LoadUser loads the user data
func (client *EdupageClient) LoadUser() error {
	u := fmt.Sprintf("https://%s/user/?", client.Credentials.Server)

	response, err := client.Credentials.httpClient.Get(u)
	if err != nil {
		return ErrAuthorization // most likely case
	}

	if response.StatusCode != 200 {
		return fmt.Errorf("server returned code: %d", response.StatusCode)
	}

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body: %s", err)
	}

	hash, err := findGSCEhash(body)
	if err != nil {
		return fmt.Errorf("failed to parse user json into json object: %s", err)
	}

	client.gsechash = hash

	js, err := findUserHome(body)
	if err != nil {
		return fmt.Errorf("failed to parse user json into json object: %s", err)
	}

	err = json.Unmarshal([]byte(js), &client.User)
	if err != nil {
		return fmt.Errorf("failed to parse user json into json object: %s", err)
	}

	return nil
}

// LoadResults loads the grade data from specified year and halfyear
// Halfyears types are: P1 (first halfyear), P2 (second halfyear), RX (whole year)
func (client *EdupageClient) LoadResults(year, halfyear string) error {
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
		return fmt.Errorf("failed to create payload: %s", err)
	}

	response, err := client.Credentials.httpClient.PostForm(url, form)
	if err != nil {
		return ErrAuthorization // most likely case
	}

	if response.StatusCode != 200 {
		return fmt.Errorf("server returned code:%d", response.StatusCode)
	}

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body: %s", err)
	}

	decoded_body := make([]byte, base64.StdEncoding.DecodedLen(len(body)-4))

	_, err = base64.StdEncoding.Decode(decoded_body, body[4:])
	if err != nil {
		return fmt.Errorf("failed to decode response body: %s", err)
	}

	decoded_body = bytes.Trim(decoded_body, "\x00")

	results, err := model.ParseResults(decoded_body)
	if err != nil {
		return fmt.Errorf("failed to parse results: %s", err)
	}

	if client.Results == nil {
		client.Results = &results
	} else {
		client.Results.Merge(&results)
	}

	return nil
}

func (client *EdupageClient) LoadTimetable(datefrom, dateto time.Time) error {
	u := fmt.Sprintf("https://%s/timetable/server/currenttt.js?__func=curentttGetData", client.Credentials.Server)

	id, err := client.GetStudentID()
	if err == ErrorUnitialized {
		return errors.New("failed to create request, user is not initialized")
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
		return fmt.Errorf("failed to create request: %s", err)
	}

	response, err := client.Credentials.httpClient.Post(u, "application/json", bytes.NewBuffer(request_body))
	if err != nil {
		return ErrAuthorization // most likely case
	}

	if response.StatusCode != 200 {
		return fmt.Errorf("server returned code: %d", response.StatusCode)
	}

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body: %s", err)
	}

	tt, err := model.ParseTimetable(body)
	if err != nil {
		return fmt.Errorf("failed to read response body: %s", err)
	}

	if client.Timetable == nil {
		client.Timetable = &tt
	} else {
		client.Timetable.Merge(&tt)
	}

	return nil
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
