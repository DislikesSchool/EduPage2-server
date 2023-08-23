package edupage

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"path"
	"reflect"
	"regexp"
	"sort"
)

var (
	timelinePath               = "timeline/"
	TimeFormat                 = "2006-01-02 15:04:05"
	ErrUnobtainableAttachments = errors.New("couldn't obtain attachments")
)

// Timeline contains all timeline information
type Timeline struct {
	Raw           map[string]interface{} // unparsed json mapping
	Homeworks     []Homework             `json:"homeworks"`
	TimelineItems []TimelineItem         `json:"timelineItems"`
}

type TimelineItem struct {
	Timeline        *Timeline
	TimelineID      string           `json:"timelineid"`
	Timestamp       Time             `json:"timestamp"`
	ReactionTo      string           `json:"reakcia_na"`
	Type            TimelineItemType `json:"typ"`
	User            string           `json:"user"`
	TargetUser      string           `json:"target_user"`
	UserName        string           `json:"user_meno"`
	OtherID         string           `json:"ineid"`
	Text            string           `json:"text"`
	TimeAdded       Time             `json:"cas_pridania"`
	TimeEvent       Time             `json:"cas_udalosti"`
	Data            TimelineData     `json:"data"`
	Owner           string           `json:"vlastnik"`
	OwnerName       string           `json:"vlastnik_meno"`
	ReactionCount   int              `json:"poct_reakcii"`
	LastReaction    string           `json:"posledna_reakcia"`
	PomocnyZaznam   string           `json:"pomocny_zaznam"`
	Removed         JSONNumber       `json:"removed"`
	TimeAddedBTC    Time             `json:"cas_pridania_btc"`
	LastReactionBTC Time             `json:"cas_udalosti_btc"`
}

type Homework struct {
	HomeworkID        string           `json:"homeworkid"`
	ESuperID          string           `json:"e_superid"`
	UserID            string           `json:"userid"`
	LessonID          JSONNumber       `json:"predmetid"`
	PlanID            string           `json:"planid"`
	Name              string           `json:"name"`
	Details           string           `json:"details"`
	DateTo            string           `json:"dateto"`
	DateFrom          string           `json:"datefrom"`
	DatetimeTo        string           `json:"datetimeto"`
	DatetimeFrom      string           `json:"datetimefrom"`
	DateCreated       string           `json:"datecreated"`
	Period            interface{}      `json:"period"`
	Timestamp         string           `json:"timestamp"`
	TestID            string           `json:"testid"`
	Type              TimelineItemType `json:"typ"`
	LikeCount         JSONNumber       `json:"pocet_like"`
	ReactionCount     JSONNumber       `json:"pocet_reakcii"`
	DoneCount         JSONNumber       `json:"pocet_done"`
	State             string           `json:"stav"`
	LastResult        string           `json:"posledny_vysledok"`
	Groups            []string         `json:"skupiny"`
	HWKID             string           `json:"hwkid"`
	ETestCards        int              `json:"etestCards"`
	ETestAnswerCards  int              `json:"etestAnswerCards"`
	StudyTopics       bool             `json:"studyTopics"`
	GradeEventID      interface{}      `json:"znamky_udalostid"`
	StudentsHidden    string           `json:"students_hidden"`
	Data              TimelineData     `json:"data"`
	EvaluationStatus  string           `json:"stavhodnotenia"`
	Result            interface{}      `json:"vysledok"`
	ResultsInfo       string           `json:"resultsInfo"`
	AssignmentID      string           `json:"pridelenieid"`
	Ended             interface{}      `json:"skoncil"`
	MissingNextLesson bool             `json:"missingNextLesson"`
	Attachments       interface{}      `json:"attachements"`
	AuthorName        string           `json:"autor_meno"`
	LessonName        string           `json:"predmet_meno"`
}

// Function region

// GetTimeline obtains the recent timeline data from the specified handle.
func (h *Handle) GetTimeline() (Timeline, error) {
	url := fmt.Sprintf("https://%s/%s", h.server, timelinePath)
	rs, err := h.hc.Get(url)
	if err != nil {
		return Timeline{}, fmt.Errorf("failed to fetch timeline: %s", err)
	}

	if rs.StatusCode == 302 {
		// edupage is trying to redirect us, that means an authorization error
		return Timeline{}, ErrAuthorization
	}

	if rs.StatusCode != 200 {
		return Timeline{}, fmt.Errorf("server returned code:%d", rs.StatusCode)
	}

	body, _ := io.ReadAll(rs.Body)
	text := string(body)

	rg, _ := regexp.Compile(`\.homeworklist\((.*)\);`)
	matches := rg.FindAllStringSubmatch(text, -1)
	if len(matches) == 0 {
		return Timeline{}, errors.New("homework list not found in the document body")
	}

	js := matches[0][1]
	var r map[string]interface{}
	err = json.Unmarshal([]byte(js), &r)
	if err != nil {
		return Timeline{}, fmt.Errorf("failed to parse timeline json: %s", err.Error())
	}

	var data Timeline
	data.Raw = r
	err = json.Unmarshal([]byte(js), &data)
	if err != nil {
		return Timeline{}, fmt.Errorf("failed to parse timeline json: %s", err.Error())
	}
	return data, nil
}

func (t *Timeline) SortedTimelineItems(predicate func(TimelineItem) bool) []TimelineItem {
	var a []TimelineItem
	if predicate != nil {
		for _, item := range t.TimelineItems {
			if predicate(item) {
				item.Timeline = t
				a = append(a, item)
			}
		}
	} else {
		a = t.TimelineItems
	}

	sort.Slice(a, func(i, j int) bool {
		return a[i].TimeAdded.Time.After(a[j].TimeAdded.Time)
	})
	return a
}
func (t *Timeline) FindHomework(superid string) (Homework, error) {
	for _, homework := range t.Homeworks {
		if homework.ESuperID == superid {
			return homework, nil
		}
	}
	return Homework{}, errors.New("homework not found")
}

func (i *TimelineItem) IsHomeworkWithAttachments() bool {
	if i.Type == TimelineHomework {
		if superid, ok := i.Data.Value["superid"]; ok && superid != nil && reflect.TypeOf(superid).Kind() == reflect.String {
			if etc, ok := i.Data.Value["etestCards"]; ok && etc != nil && etc.(float64) == 1 {
				return true
			}
		}
	}

	return false
}

func (i *TimelineItem) ToHomework() (Homework, error) {
	if i.Type == TimelineHomework {
		if superid, ok := i.Data.Value["superid"]; ok && superid != nil && reflect.TypeOf(superid).Kind() == reflect.String {
			if etc, ok := i.Data.Value["etestCards"]; ok && etc != nil && etc.(float64) == 1 {
				return i.Timeline.FindHomework(superid.(string))
			}
		}
	}

	return Homework{}, errors.New("not a homework")
}

func (i *TimelineItem) GetAttachments() (map[string]string, error) {
	if i.Type == TimelineMessage {
		var attachments = make(map[string]string)
		data := i.Data
		if val, ok := data.Value["attachements"]; ok { // It's misspelled in the JSON payload
			if reflect.TypeOf(val).Kind() == reflect.Map {
				a := val.(map[string]interface{})
				for k, v := range a {
					attachments[v.(string)] = k
				}
			}
		}
		return attachments, nil
	}
	return nil, ErrUnobtainableAttachments
}

func (h *Handle) FetchHomeworkAttachments(i *Homework) (map[string]string, error) {
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

	resp, err := h.hc.PostForm(
		"https://"+path.Join(h.server, "elearning", "?cmd=MaterialPlayer&akcia=getETestData"),
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
