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

func (client *EdupageClient) fetchTimeline(datefrom, dateto time.Time) (model.Timeline, error) {
	if client.Credentials.httpClient == nil {
		return model.Timeline{}, errors.New("invalid credentials")
	}

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
		return model.Timeline{}, ErrorUnauthorized // most likely case
	}

	if response.StatusCode != 200 {
		return model.Timeline{}, fmt.Errorf("server returned code: %d", response.StatusCode)
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
		return model.Timeline{}, fmt.Errorf("failed to parse timeline json: %s", err)
	}

	return timeline, nil
}

func (client *EdupageClient) fetchUser() (model.User, error) {
	if client.Credentials.httpClient == nil {
		return model.User{}, errors.New("invalid credentials")
	}
	u := fmt.Sprintf("https://%s/user/?", client.Credentials.Server)

	response, err := client.Credentials.httpClient.Get(u)
	if err != nil {
		return model.User{}, ErrorUnauthorized // most likely case
	}

	if response.StatusCode != 200 {
		return model.User{}, fmt.Errorf("server returned code: %d", response.StatusCode)
	}

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return model.User{}, fmt.Errorf("failed to read response body: %s", err)
	}

	hash, err := findGSCEHash(body)
	if err != nil {
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

func (client *EdupageClient) fetchResults(year, halfyear string) (model.Results, error) {
	if client.Credentials.httpClient == nil {
		return model.Results{}, errors.New("invalid credentials")
	}

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
		return model.Results{}, ErrorUnauthorized // most likely case
	}

	if response.StatusCode != 200 {
		return model.Results{}, fmt.Errorf("server returned code: %d", response.StatusCode)
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
	if client.Credentials.httpClient == nil {
		return model.Timetable{}, errors.New("invalid credentials")
	}

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
		return model.Timetable{}, ErrorUnauthorized // most likely case
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
		return model.Timetable{}, fmt.Errorf("failed to parse timetable json: %s", err)
	}

	return tt, nil
}

func (client *EdupageClient) fetchCanteen(date time.Time) (model.Canteen, error) {
	if client.Credentials.httpClient == nil {
		return model.Canteen{}, errors.New("invalid credentials")
	}
	u := fmt.Sprintf("https://%s/menu/?date=%s", client.Credentials.Server, date.Format("20060102"))

	response, err := client.Credentials.httpClient.Get(u)
	if err != nil {
		return model.Canteen{}, ErrorUnauthorized // most likely case
	}

	if response.StatusCode != 200 {
		return model.Canteen{}, fmt.Errorf("server returned code: %d", response.StatusCode)
	}

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return model.Canteen{}, fmt.Errorf("failed to read response body: %s", err)
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
	rg, _ := regexp.Compile(`ASC\.gsechash="(.*)";`)
	matches := rg.FindAllStringSubmatch(string(body), -1)
	if len(matches) == 0 {
		return "", errors.New("gsechash not found in the document body")
	}

	return matches[0][1], nil
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
	rg, _ := regexp.Compile(`\.userhome\((.*)\);`)
	matches := rg.FindAllStringSubmatch(string(body), -1)
	if len(matches) == 0 {
		return "", errors.New("userhome not found in the document body")
	}

	return matches[0][1], nil
}
