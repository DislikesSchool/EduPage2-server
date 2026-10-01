package edupage

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/DislikesSchool/EduPage2-server/edupage/model"
)

func doAuthedGet(client *http.Client, rawurl, referer string) (*http.Response, []byte, error) {
	req, err := http.NewRequest("GET", rawurl, nil)
	if err != nil {
		return nil, nil, err
	}
	req.Header.Set("User-Agent", EduPageUserAgent)
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8")
	req.Header.Set("Accept-Language", "sk-SK,sk;q=0.9,cs;q=0.8,en;q=0.7")
	if referer != "" {
		req.Header.Set("Referer", referer)
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, nil, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(io.LimitReader(resp.Body, 20<<20))
	if err != nil {
		return resp, nil, err
	}
	return resp, body, nil
}

func doAuthedPostForm(client *http.Client, rawurl string, form map[string]string, referer string) (*http.Response, []byte, error) {
	payload := CreatePayload(form)
	// CreatePayload returns url.Values with eqap/eqacs/eqaz.
	// Encode with standard QueryEscape; server accepts both + and %20.
	// Use Encode() then POST as x-www-form-urlencoded with eqav params.
	u := rawurl
	if !strings.Contains(u, "eqav=") {
		sep := "?"
		if strings.Contains(u, "?") {
			sep = "&"
		}
		u = u + sep + "eqav=1&maxEqav=7"
	}
	req, err := http.NewRequest("POST", u, strings.NewReader(payload.Encode()))
	if err != nil {
		return nil, nil, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")
	req.Header.Set("User-Agent", EduPageUserAgent)
	req.Header.Set("Accept", "*/*")
	req.Header.Set("X-Requested-With", "XMLHttpRequest")
	if referer != "" {
		req.Header.Set("Referer", referer)
		req.Header.Set("Origin", referer)
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, nil, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(io.LimitReader(resp.Body, 20<<20))
	if err != nil {
		return resp, nil, err
	}
	return resp, body, nil
}

func (client *EdupageClient) fetchTimelineModel(datefrom, dateto time.Time) (model.Timeline, error) {
	if client.Credentials.httpClient == nil {
		return model.Timeline{}, errors.New("invalid credentials")
	}

	url := fmt.Sprintf("https://%s/timeline/?akcia=getData", client.Credentials.Server)

	_, body, err := doAuthedPostForm(client.Credentials.httpClient, url, map[string]string{
		"datefrom": datefrom.Format("2006-01-02"),
		"dateto":   dateto.Format("2006-01-02"),
	}, fmt.Sprintf("https://%s/timeline/", client.Credentials.Server))
	if err != nil {
		return model.Timeline{}, ErrorUnauthorized // most likely case
	}

	decodedBody, err := DecodeLegacyBody(body)
	if err != nil {
		return model.Timeline{}, fmt.Errorf("failed to decode response body: %s", err)
	}

	timeline, err := model.ParseTimeline(decodedBody)
	if err != nil {
		// New backend may return plain JSON with timelineItems (see edupage-api).
		// Try to be helpful: if body was eqwd, surface it.
		if strings.HasPrefix(strings.TrimSpace(string(body)), "eqwd:") {
			return model.Timeline{}, fmt.Errorf("server asked to retry (eqwd); failed to parse timeline json: %s", err)
		}
		return model.Timeline{}, fmt.Errorf("failed to parse timeline json: %s", err)
	}

	return timeline, nil
}

func (client *EdupageClient) fetchUserModel() (model.User, error) {
	if client.Credentials.httpClient == nil {
		return model.User{}, errors.New("invalid credentials")
	}
	u := fmt.Sprintf("https://%s/user/?", client.Credentials.Server)

	_, body, err := doAuthedGet(client.Credentials.httpClient, u, "")
	if err != nil {
		return model.User{}, ErrorUnauthorized // most likely case
	}

	hash, err := findGSCEHash(body)
	if err != nil {
		// If login page was returned (not logged in), surface unauthorized.
		if strings.Contains(string(body), "MainLogin") || strings.Contains(string(body), "csrftoken") {
			return model.User{}, ErrorUnauthorized
		}
		return model.User{}, fmt.Errorf("failed to parse user json: %s", err)
	}

	client.gsechash = hash

	js, err := findUserHome(body)
	if err != nil {
		return model.User{}, fmt.Errorf("failed to parse user json: %s", err)
	}

	var user model.User
	err = json.Unmarshal([]byte(js), &user)
	if err != nil {
		return model.User{}, fmt.Errorf("failed to parse user json: %s", err)
	}

	return user, nil
}

func (client *EdupageClient) fetchResultsModel(year, halfyear string) (model.Results, error) {
	if client.Credentials.httpClient == nil {
		return model.Results{}, errors.New("invalid credentials")
	}

	url := fmt.Sprintf("https://%s/znamky/?what=studentviewer&akcia=studentData&eqav=1&maxEqav=7", client.Credentials.Server)

	_, body, err := doAuthedPostForm(client.Credentials.httpClient, url, map[string]string{
		"pohlad":           "podladatumu",
		"znamky_yearid":    year,
		"znamky_yearid_ns": "1",
		"nadobdobie":       halfyear,
		"rokobdobie":       fmt.Sprintf("%s::%s", year, halfyear),
		"doRq":             "1",
		"what":             "studentviewer",
		"updateLastView":   "0",
	}, fmt.Sprintf("https://%s/znamky/", client.Credentials.Server))
	if err != nil {
		return model.Results{}, ErrorUnauthorized // most likely case
	}

	decodedBody, err := DecodeLegacyBody(body)
	if err != nil {
		return model.Results{}, fmt.Errorf("failed to decode response body: %s", err)
	}

	results, err := model.ParseResults(decodedBody)
	if err != nil {
		return model.Results{}, fmt.Errorf("failed to parse results: %s", err)
	}

	return results, nil
}

func (client *EdupageClient) fetchTimetableModel(datefrom, dateto time.Time) (model.Timetable, error) {
	if client.Credentials.httpClient == nil {
		return model.Timetable{}, errors.New("invalid credentials")
	}

	u := fmt.Sprintf("https://%s/timetable/server/currenttt.js?__func=curentttGetData", client.Credentials.Server)

	id, err := client.GetStudentID()
	if err == ErrorUnitialized {
		return model.Timetable{}, errors.New("failed to create request, user is not initialized")
	}

	year, currentMonth, _ := datefrom.Date()
	if currentMonth < 9 {
		year--
	}

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

	requestBody, err := json.Marshal(request)
	if err != nil {
		return model.Timetable{}, fmt.Errorf("failed to create request: %s", err)
	}

	req, err := http.NewRequest("POST", u, bytes.NewBuffer(requestBody))
	if err != nil {
		return model.Timetable{}, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", EduPageUserAgent)
	req.Header.Set("Accept", "application/json, text/javascript, */*; q=0.01")
	req.Header.Set("X-Requested-With", "XMLHttpRequest")
	req.Header.Set("Referer", fmt.Sprintf("https://%s/timetable/", client.Credentials.Server))

	response, err := client.Credentials.httpClient.Do(req)
	if err != nil {
		return model.Timetable{}, ErrorUnauthorized // most likely case
	}
	defer response.Body.Close()

	if response.StatusCode != 200 {
		return model.Timetable{}, fmt.Errorf("server returned code: %d", response.StatusCode)
	}

	body, err := io.ReadAll(io.LimitReader(response.Body, 20<<20))
	if err != nil {
		return model.Timetable{}, fmt.Errorf("failed to read response body: %s", err)
	}

	tt, err := model.ParseTimetable(body)
	if err != nil {
		return model.Timetable{}, fmt.Errorf("failed to parse timetable json: %s", err)
	}

	return tt, nil
}

func (client *EdupageClient) fetchCanteenModel(date time.Time) (model.Canteen, error) {
	if client.Credentials.httpClient == nil {
		return model.Canteen{}, errors.New("invalid credentials")
	}
	u := fmt.Sprintf("https://%s/menu/?date=%s", client.Credentials.Server, date.Format("20060102"))

	_, body, err := doAuthedGet(client.Credentials.httpClient, u, "")
	if err != nil {
		return model.Canteen{}, ErrorUnauthorized // most likely case
	}

	data, err := findEdupageData(body)
	if err != nil {
		return model.Canteen{}, fmt.Errorf("failed to parse json: %s", err)
	}

	canteen, err := model.ParseCanteen([]byte(data))
	if err != nil {
		return model.Canteen{}, fmt.Errorf("failed to parse json: %s", err)
	}

	return canteen, nil
}

func findGSCEHash(body []byte) (string, error) {
	// Current pages: ASC.gsechash="abc123";  Old: ASC.gsechash="..."
	// Be non-greedy and support single quotes.
	patterns := []string{
		`ASC\.gsechash\s*=\s*"([^"]+)"`,
		`ASC\.gsechash\s*=\s*'([^']+)'`,
		`gsechash\s*=\s*"([^"]+)"`,
		`ASC\.gsechash="(.*)";`,
	}
	s := string(body)
	for _, p := range patterns {
		rg, _ := regexp.Compile(p)
		matches := rg.FindAllStringSubmatch(s, -1)
		if len(matches) > 0 && len(matches[0]) > 1 && matches[0][1] != "" && matches[0][1] != "00000000" {
			return matches[0][1], nil
		}
	}
	// Fallback: allow 00000000 (login1 landing page) with error context.
	rg, _ := regexp.Compile(`ASC\.gsechash="([^"]*)"`)
	if m := rg.FindStringSubmatch(s); len(m) == 2 {
		if m[1] == "" || m[1] == "00000000" {
			return "", errors.New("gsechash not found in the document body (not logged in or login page returned)")
		}
		return m[1], nil
	}
	return "", errors.New("gsechash not found in the document body")
}

func findEdupageData(body []byte) (string, error) {
	rg, _ := regexp.Compile(`edupageData: (\{.*\}),`)
	matches := rg.FindAllStringSubmatch(string(body), -1)
	if len(matches) == 0 {
		return "", errors.New("edupageData not found in the document body")
	}

	return matches[0][1], nil
}

func findUserHome(body []byte) (string, error) {
	// Modern pages: userhome({...}); possibly with whitespace/newlines.
	// Python splits on "userhome(" and rsplit(");",2)[0] - replicate robustly.
	s := string(body)
	idx := strings.Index(s, "userhome(")
	if idx == -1 {
		// Try .userhome( (old regex)
		rg, _ := regexp.Compile(`\.userhome\((.*)\);`)
		matches := rg.FindAllStringSubmatch(s, -1)
		if len(matches) == 0 {
			return "", errors.New("userhome not found in the document body")
		}
		return matches[0][1], nil
	}
	rest := s[idx+len("userhome("):]
	// Find matching closing ");" - use last occurrence to handle nested?
	// Python: rsplit(");",2)[0] - take up to second-last ");".
	// We implement: find last ");" then second-last? Simpler: greedy to last ");".
	end := strings.LastIndex(rest, ");")
	if end == -1 {
		return "", errors.New("userhome not found in the document body")
	}
	candidate := rest[:end]
	// If there are multiple userhome calls, LastIndex may include extra;
	// try to JSON-validate, else fallback to first balanced parse.
	trimmed := strings.TrimSpace(candidate)
	if json.Valid([]byte(trimmed)) {
		return trimmed, nil
	}
	// Fallback to regex non-greedy first match.
	rg, _ := regexp.Compile(`userhome\((\{.*?\})\)`)
	if m := rg.FindStringSubmatch(s); len(m) == 2 {
		return m[1], nil
	}
	// Return candidate anyway; Unmarshal will error with context.
	return candidate, nil
}
