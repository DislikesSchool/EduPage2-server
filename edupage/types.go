package edupage

import (
	"encoding/json"
	"strconv"
	"strings"
	"time"
)

type TimelineType struct {
	uint8
}

var (
	TimelineMessage  = TimelineType{0}
	TimelineHomework = TimelineType{1}
	TimelineInvalid  = TimelineType{2}
)

func (n *TimelineType) UnmarshalJSON(b []byte) error {
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

func (n *TimelineType) MarshalJSON() ([]byte, error) {
	return []byte{n.uint8}, nil
}

type Data map[string]interface{}

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
		fix := []byte("{\"data\":" + r + "}")
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

type Number struct {
	int64
}

func (n *Number) UnmarshalJSON(b []byte) error {
	var s = string(b)
	s = strings.ReplaceAll(s, "\"", "")
	var err error
	n.int64, err = strconv.ParseInt(s, 10, 64)
	return err
}

func (n *Number) MarshalJSON() ([]byte, error) {
	return []byte(strconv.FormatInt(n.int64, 10)), nil
}

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
