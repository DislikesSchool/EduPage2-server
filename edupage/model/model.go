package model

import (
	"encoding/json"
	"strconv"
)

var (
	TimeFormatYearMonthDay = "2006-01-02"
)

type Mergeable interface {
	Merge(src *Mergeable)
}

// StringJsonObject contains stringed json object
type StringJsonObject struct {
	Value map[string]interface{}
}

func (n *StringJsonObject) UnmarshalJSON(b []byte) error {
	// Handle the case where the value was marshaled by the marshal method
	err := json.Unmarshal(b, &n.Value)
	if err == nil {
		return nil
	}

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

func UnmarshalNestedString(src string, dst any) error {
	r := string(src)
	if r == "[]" {
		return nil
	}
	s, err := strconv.Unquote(r)

	if err != nil {
		fix := []byte("{\"data\":" + r + "}") // weird fix, TODO: not do it this way?
		var temp map[string]interface{}
		json.Unmarshal(fix, &temp)
		json.Unmarshal([]byte(temp["data"].(string)), &dst)
	} else {
		if err := json.Unmarshal([]byte(s), &dst); err != nil {
			return err
		}
	}
	return nil
}

func (n *StringJsonObject) MarshalJSON() ([]byte, error) {
	return json.Marshal(n.Value)
}
