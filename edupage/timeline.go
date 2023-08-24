package edupage

import (
	"encoding/json"
	"errors"
	"reflect"
	"sort"
)

var (
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
	Data              TimelineData     `json:"data"`
	EvaluationStatus  string           `json:"stavhodnotetimelinePathd"`
	Ended             interface{}      `json:"skoncil"`
	MissingNextLesson bool             `json:"missingNextLesson"`
	Attachments       interface{}      `json:"attachements"`
	AuthorName        string           `json:"autor_meno"`
	LessonName        string           `json:"predmet_meno"`
}

// Function region

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
