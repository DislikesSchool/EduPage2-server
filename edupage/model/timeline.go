package model

import (
	"encoding/json"
	"errors"
	"reflect"
	"strconv"
	"strings"
	"time"
)

var (
	ItemTypeMessage  = TimelineItemType{0}
	ItemTypeHomework = TimelineItemType{1}
	TYPE_INVALID     = TimelineItemType{2}
)

var (
	ErrUnobtainableAttachments = errors.New("couldn't obtain attachments")
	TimeFormat                 = "2006-01-02 15:04:05"
)

// Timeline contains all timeline information
type Timeline struct {
	Homeworks []Homework     `json:"homeworks"`
	Items     []TimelineItem `json:"timelineItems"`
}

type TimelineItemType struct {
	uint8
}

// TimelineItemData contains raw timeline data
type TimelineItemData struct {
	Value map[string]interface{}
}

// Time is a representation of time instances to help with parsing.
type Time struct {
	time.Time
}

func (n *Time) UnmarshalJSON(b []byte) error {
	s := string(b)
	s = strings.ReplaceAll(s, "\"", "")
	n.Time, _ = time.Parse(TimeFormat, s)
	return nil
}

func (n *Time) MarshalJSON() ([]byte, error) {
	return []byte(n.Time.Format(TimeFormat)), nil
}

type TimelineItem struct {
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
	Data            TimelineItemData `json:"data"`
	Owner           string           `json:"vlastnik"`
	OwnerName       string           `json:"vlastnik_meno"`
	ReactionCount   int              `json:"poct_reakcii"`
	LastReaction    string           `json:"posledna_reakcia"`
	PomocnyZaznam   string           `json:"pomocny_zaznam"`
	Removed         json.Number      `json:"removed"`
	TimeAddedBTC    Time             `json:"cas_pridania_btc"`
	LastReactionBTC Time             `json:"cas_udalosti_btc"`
}

type Homework struct {
	HomeworkID        string           `json:"homeworkid"`
	ESuperID          string           `json:"e_superid"`
	UserID            string           `json:"userid"`
	LessonID          json.Number      `json:"predmetid"`
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
	LikeCount         json.Number      `json:"pocet_like"`
	ReactionCount     json.Number      `json:"pocet_reakcii"`
	DoneCount         json.Number      `json:"pocet_done"`
	State             string           `json:"stav"`
	LastResult        string           `json:"posledny_vysledok"`
	Groups            []string         `json:"skupiny"`
	HWKID             string           `json:"hwkid"`
	ETestCards        int              `json:"etestCards"`
	ETestAnswerCards  int              `json:"etestAnswerCards"`
	StudyTopics       bool             `json:"studyTopics"`
	GradeEventID      interface{}      `json:"znamky_udalostid"`
	StudentsHidden    string           `json:"students_hidden"`
	Data              TimelineItemData `json:"data"`
	EvaluationStatus  string           `json:"stavhodnotetimelinePathd"`
	Ended             interface{}      `json:"skoncil"`
	MissingNextLesson bool             `json:"missingNextLesson"`
	Attachments       interface{}      `json:"attachements"`
	AuthorName        string           `json:"autor_meno"`
	LessonName        string           `json:"predmet_meno"`
}

func (t *Timeline) GetHomework(superid string) (Homework, error) {
	for _, homework := range t.Homeworks {
		if homework.ESuperID == superid {
			return homework, nil
		}
	}
	return Homework{}, errors.New("homework not found")
}

func (i *TimelineItem) IsHomeworkWithAttachments() bool {
	if i.Type == ItemTypeHomework {
		if superid, ok := i.Data.Value["superid"]; ok && superid != nil && reflect.TypeOf(superid).Kind() == reflect.String {
			if etc, ok := i.Data.Value["etestCards"]; ok && etc != nil && etc.(float64) == 1 {
				return true
			}
		}
	}

	return false
}

func (i *TimelineItem) GetAttachments() (map[string]string, error) {
	if i.Type == ItemTypeMessage {
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

func (n *TimelineItemType) UnmarshalJSON(b []byte) error {
	s := string(b)
	if s == "\"sprava\"" {
		n.uint8 = 0
	} else if s == "\"homework\"" {
		n.uint8 = 1
	} else {
		n.uint8 = 2
	}
	return nil
}

func (n *TimelineItemType) MarshalJSON() ([]byte, error) {
	return []byte{n.uint8}, nil
}

func (n *TimelineItemData) UnmarshalJSON(b []byte) error {
	r := string(b)
	if r == "[]" {
		n.Value = make(map[string]interface{})
		return nil
	}
	s, err := strconv.Unquote(r)

	if err != nil {
		fix := []byte("{\"data\":" + r + "}") // weird fix, TODO: not do it this way?
		var temp map[string]interface{}
		_ = json.Unmarshal(fix, &temp)
		_ = json.Unmarshal([]byte(temp["data"].(string)), &n.Value)
	} else {
		if err := json.Unmarshal([]byte(s), &n.Value); err != nil {
			n.Value = make(map[string]interface{})
		}
	}
	return nil
}

func (n *TimelineItemData) MarshalJSON() ([]byte, error) {
	return json.Marshal(n.Value)
}
