package edupage

import (
	"encoding/json"
	"strconv"
	"strings"
	"time"
)

// TimelineItemType represents the timeline item's type, so it can be handled correctly according to type.
type TimelineItemType struct {
	uint8
}

var (
	TimelineMessage  = TimelineItemType{0}
	TimelineHomework = TimelineItemType{1}
	TimelineInvalid  = TimelineItemType{2}
)

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

type Data map[string]interface{}

// TimelineData contains raw timeline data
type TimelineData struct {
	Value Data
}

func (n *TimelineData) UnmarshalJSON(b []byte) error {
	r := string(b)
	if r == "[]" {
		n.Value = Data{}
		return nil
	}
	s, err := strconv.Unquote(r)

	if err != nil {
		fix := []byte("{\"data\":" + r + "}") // weird fix, TODO: fix
		var temp Data
		_ = json.Unmarshal(fix, &temp)
		_ = json.Unmarshal([]byte(temp["data"].(string)), &n.Value)
	} else {
		if err := json.Unmarshal([]byte(s), &n.Value); err != nil {
			n.Value = Data{}
		}
	}
	return nil
}

func (n *TimelineData) MarshalJSON() ([]byte, error) {
	return json.Marshal(n.Value)
}

// JSONNumber is a robust representation of an integer in json to help with parsing
type JSONNumber struct {
	int64
}

func (n *JSONNumber) UnmarshalJSON(b []byte) error {
	var s = string(b)
	s = strings.ReplaceAll(s, "\"", "")
	var err error
	n.int64, err = strconv.ParseInt(s, 10, 64)
	return err
}

func (n *JSONNumber) MarshalJSON() ([]byte, error) {
	return []byte(strconv.FormatInt(n.int64, 10)), nil
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
