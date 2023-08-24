package edupage

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"path"
	"reflect"
	"time"
)

// Edupage is used to access the Edupage API.
type Edupage struct {
	hc       *http.Client
	server   string
	Timeline Timeline
}

// LoadRecentTimeline loads the recent timeline data.
// That's from today, to 30 days in the past.
// Also updates the Timeline property in Edupage struct.
func (edupage *Edupage) LoadRecentTimeline() (Timeline, error) {
	duration, err := time.ParseDuration("-720h") // 30 days
	if err != nil {
		return Timeline{}, fmt.Errorf("failed to parse duration: %s", err)
	}

	start := time.Now().Add(duration)
	return edupage.LoadTimeline(start, time.Now())
}

// LoadTimeline loads the timeline data from the specified date range.
// Also updates the Timeline property in Edupage struct.
func (h *Edupage) LoadTimeline(datefrom, dateto time.Time) (Timeline, error) {
	url := fmt.Sprintf("https://%s/timeline/?akcia=getData", h.server)

	form, err := CreatePayload(map[string]string{
		"datefrom": datefrom.Format("2006-01-02"),
		"dateto":   dateto.Format("2006-01-02"),
	})

	if err != nil {
		return Timeline{}, fmt.Errorf("failed to create payload: %s", err)
	}

	response, err := h.hc.PostForm(url, form)
	if err != nil {
		return Timeline{}, fmt.Errorf("failed to fetch timeline: %s", err)
	}

	if response.StatusCode == 302 {
		// edupage is trying to redirect us, that means an authorization error
		return Timeline{}, ErrAuthorization
	}

	if response.StatusCode != 200 {
		return Timeline{}, fmt.Errorf("server returned code:%d", response.StatusCode)
	}

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return Timeline{}, fmt.Errorf("failed to read response body: %s", err)
	}

	decoded_body := make([]byte, base64.StdEncoding.DecodedLen(len(body)-4))

	_, err = base64.StdEncoding.Decode(decoded_body, body[4:])

	if err != nil {
		return Timeline{}, fmt.Errorf("failed to decode response body: %s", err)
	}

	var raw map[string]interface{}
	err = json.Unmarshal(decoded_body[0:len(decoded_body)-1], &raw) // omitting null character at end
	if err != nil {
		return Timeline{}, fmt.Errorf("failed to parse timeline json into raw object: %s", err.Error())
	}

	var data Timeline
	err = json.Unmarshal(decoded_body[0:len(decoded_body)-1], &data)
	if err != nil {
		return Timeline{}, fmt.Errorf("failed to parse timeline json: %s", err.Error())
	}

	data.Raw = raw

	h.Timeline = data
	return data, nil
}

// FetchHomeworkAttachmens obtains the homework attchments for the specified homework.
// Returns ErrUnobtainableAttachments in case the attachments are not present.
// Retruns map, key is the resource name and value is the resource link
func (edupage *Edupage) FetchHomeworkAttachments(i *Homework) (map[string]string, error) {
	if i.ESuperID == "" || i.TestID == "" {
		return nil, errors.New("required fields superid and testid not set")
	}

	data := map[string]string{
		"testid":  i.TestID,
		"superid": i.ESuperID,
	}

	payload, err := CreatePayload(data)
	if err != nil {
		return nil, fmt.Errorf("failed to create payload: %w", err)
	}

	resp, err := edupage.hc.PostForm(
		"https://"+path.Join(edupage.server, "elearning", "?cmd=MaterialPlayer&akcia=getETestData"),
		payload,
	)
	if err != nil {
		return nil, fmt.Errorf("homework request failed: %w", err)
	}

	response, err := io.ReadAll(resp.Body)

	if len(response) < 5 {
		return nil, fmt.Errorf("homework request failed, bad response: %w", err)
	}

	response = response[4:]

	decoded := make([]byte, base64.StdEncoding.DecodedLen(len(response)))
	_, err = base64.StdEncoding.Decode(decoded, response)
	if err != nil {
		return nil, fmt.Errorf("homework request failed, bad response: %w", err)
	}

	decoded = bytes.Trim(decoded, "\x00")
	var object map[string]interface{}
	err = json.Unmarshal(decoded, &object)
	if err != nil {
		return nil, fmt.Errorf("homework request failed, bad response: %w", err)
	}

	attachments := make(map[string]string)

	// God help those who may try to debug this.
	if object["materialData"] == nil ||
		(reflect.TypeOf(object["materialData"]).Kind() != reflect.Map ||
			reflect.TypeOf(object["materialData"]).Elem().Kind() != reflect.Interface) {
		return nil, ErrUnobtainableAttachments
	}
	materialData := object["materialData"].(map[string]interface{})

	if materialData["cardsData"] == nil ||
		(reflect.TypeOf(materialData["cardsData"]).Kind() != reflect.Map ||
			reflect.TypeOf(materialData["cardsData"]).Elem().Kind() != reflect.Interface) {
		return nil, ErrUnobtainableAttachments
	}
	cardsData := materialData["cardsData"].(map[string]interface{})

	for _, entry := range cardsData {
		if entry == nil ||
			(reflect.TypeOf(entry).Kind() != reflect.Map ||
				reflect.TypeOf(entry).Elem().Kind() != reflect.Interface) {
			return nil, ErrUnobtainableAttachments
		}

		if e, ok := entry.(map[string]interface{})["content"]; !ok && reflect.TypeOf(e).Kind() != reflect.String {
			return nil, ErrUnobtainableAttachments
		}

		var content map[string]interface{}
		contentJson := entry.(map[string]interface{})["content"].(string)
		err = json.Unmarshal([]byte(contentJson), &content)
		if err != nil {
			return nil, err
		}

		if content["widgets"] == nil ||
			(reflect.TypeOf(content["widgets"]).Kind() != reflect.Slice ||
				reflect.TypeOf(content["widgets"]).Elem().Kind() != reflect.Interface) {
			return nil, ErrUnobtainableAttachments
		}

		widgets := content["widgets"].([]interface{})
		for _, widget := range widgets {
			if widget == nil ||
				(reflect.TypeOf(widget).Kind() != reflect.Map ||
					reflect.TypeOf(widget).Elem().Kind() != reflect.Interface) {
				return nil, ErrUnobtainableAttachments
			}
			if widget.(map[string]interface{})["props"] == nil ||
				(reflect.TypeOf(widget.(map[string]interface{})["props"]).Kind() != reflect.Map ||
					reflect.TypeOf(widget.(map[string]interface{})["props"]).Elem().Kind() != reflect.Interface) {
				return nil, ErrUnobtainableAttachments
			}
			props := widget.(map[string]interface{})["props"].(map[string]interface{})
			if files, ok := props["files"]; ok {
				for _, file := range files.([]interface{}) {
					if file == nil ||
						(reflect.TypeOf(file).Kind() != reflect.Map ||
							reflect.TypeOf(file).Elem().Kind() != reflect.Interface) {
						return nil, ErrUnobtainableAttachments
					}
					attachments[file.(map[string]interface{})["name"].(string)] = file.(map[string]interface{})["src"].(string)
				}
			}
		}
		if err != nil {
			continue
		}
		continue
	}

	return attachments, nil
}
