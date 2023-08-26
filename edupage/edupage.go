package edupage

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"time"

	"github.com/DislikesSchool/EduPage2-server/edupage/model"
)

// EdupageClient is used to access the edupage api.
type EdupageClient struct {
	hc     *http.Client //TODO: remove, save only edid, hsid, phpsessid tokens
	server string

	User     *model.User
	Timeline *model.Timeline
	Results  *model.Results
}

// Fetch loads all possible data into the object
func (client *EdupageClient) Fetch() error {
	err := client.LoadUser()
	if err != nil {
		return err
	}

	err = client.LoadRecentTimeline()
	if err != nil {
		return err
	}

	err = client.LoadResults("2022", "RX")
	if err != nil {
		return err
	}

	return nil
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
	url := fmt.Sprintf("https://%s/timeline/?akcia=getData", client.server)

	form, err := CreatePayload(map[string]string{
		"datefrom": datefrom.Format("2006-01-02"),
		"dateto":   dateto.Format("2006-01-02"),
	})

	if err != nil {
		return fmt.Errorf("failed to create payload: %s", err)
	}

	response, err := client.hc.PostForm(url, form)
	if err != nil {
		return fmt.Errorf("failed to fetch timeline: %s", err)
	}

	if response.StatusCode == 302 {
		// edupage is trying to redirect us, that means an authorization error
		return ErrAuthorization
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
	url := fmt.Sprintf("https://%s/user/", client.server)
	response, err := client.hc.Get(url)
	if err != nil {
		return fmt.Errorf("failed to fetch timeline: %s", err)
	}

	if response.StatusCode == 302 {
		// edupage is trying to redirect us, that means an authorization error
		return ErrAuthorization
	}

	if response.StatusCode != 200 {
		return fmt.Errorf("server returned code: %d", response.StatusCode)
	}

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body: %s", err)
	}

	text := string(body)

	rg, _ := regexp.Compile(`\.userhome\((.*)\);`)
	matches := rg.FindAllStringSubmatch(text, -1)
	if len(matches) == 0 {
		return errors.New("userhome not found in the document body")
	}

	js := matches[0][1]
	err = json.Unmarshal([]byte(js), &client.User)
	if err != nil {
		return fmt.Errorf("failed to parse user json into json object: %s", err)
	}

	return nil
}

// LoadResults loads the grade data from specified year and halfyear
// Halfyears types are: P1 (first halfyear), P2 (second halfyear), RX (whole year)
func (client *EdupageClient) LoadResults(year, halfyear string) error {
	url := fmt.Sprintf("https://%s/znamky/?what=studentviewer&akcia=studentData&eqav=1&maxEqav=7", client.server)

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

	response, err := client.hc.PostForm(url, form)
	if err != nil {
		return fmt.Errorf("failed to fetch grades: %s", err)
	}

	if response.StatusCode == 302 {
		// edupage is trying to redirect us, that means an authorization error
		return ErrAuthorization
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
