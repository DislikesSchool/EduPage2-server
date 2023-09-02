package model

import "encoding/json"

type Timetable struct {
	Items map[string]TimetableItem
}

type TimetableItem struct {
	Type         string   `json:"type"`
	Date         string   `json:"date"`
	Period       string   `json:"uniperiod"`
	StartTime    string   `json:"starttime"`
	EndTime      string   `json:"endtime"`
	SubjectID    string   `json:"subjectid"`
	ClassIDs     []string `json:"classids"`
	GroupNames   []string `json:"groupnames"`
	IGroupID     string   `json:"igroupid"`
	TeacherIDs   []string `json:"teacherids"`
	ClassroomIDs []string `json:"classroomids"`
	StudentIDs   []string `json:"studentids"`
	Colors       []string `json:"colors"`
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

	t.Items = make(map[string]TimetableItem, len(r.Response.Items))
	for _, v := range r.Response.Items {
		t.Items[v.Date] = v
	}

	return t, nil
}
