package model

import (
	"strings"
	"time"
)

var (
	TimeFormat = "2006-01-02 15:04:05"
)

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
