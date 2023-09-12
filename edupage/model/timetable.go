package model

import (
	"encoding/json"
	"maps"
)

type Timetable struct {
	Days map[string][]TimetableItem // key format is YYYY-MM-dd or 2006-01-02
}

type TimetableItem struct {
	Type         string      `json:"type"`
	Date         string      `json:"date"`
	Period       string      `json:"uniperiod"`
	StartTime    string      `json:"starttime"`
	EndTime      string      `json:"endtime"`
	SubjectID    string      `json:"subjectid"`
	ClassIDs     []string    `json:"classids"`
	GroupNames   []string    `json:"groupnames"`
	IGroupID     string      `json:"igroupid"`
	TeacherIDs   []string    `json:"teacherids"`
	ClassroomIDs []string    `json:"classroomids"`
	StudentIDs   []string    `json:"studentids"`
	Colors       []string    `json:"colors"`
	Duration     json.Number `json:"durationperiods"`
}

func (t *Timetable) Merge(src *Timetable) {
	maps.Copy(t.Days, src.Days)
}

func ParseTimetable(data []byte) (Timetable, error) {
	type Response struct {
		Items []TimetableItem `json:"ttitems"`
	}

	type RawTimetable struct {
		Response Response `json:"r"`
	}

	var r RawTimetable

	err := json.Unmarshal(data, &r)
	if err != nil {
		return Timetable{}, err
	}

	var t Timetable

	t.Days = make(map[string][]TimetableItem, len(r.Response.Items))
	for _, new := range r.Response.Items {
		if original, ok := t.Days[new.Date]; ok {
			t.Days[new.Date] = append(original, new)
		} else {
			t.Days[new.Date] = []TimetableItem{new}
		}

	}

	return t, nil
}
