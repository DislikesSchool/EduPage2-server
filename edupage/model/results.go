package model

import (
	"encoding/json"

	"golang.org/x/exp/maps"
)

type Grade struct {
	Provider    string `json:"provider"`
	ID          string `json:"znamkaid"`
	StudentID   string `json:"studentid"`
	SubjectID   string `json:"predmetid"`
	EventID     string `json:"udalostID"`
	Month       string `json:"mesiac"`
	Data        string `json:"data"`
	Date        Time   `json:"datum"`
	TeacherID   string `json:"ucitelid"`
	Signed      string `json:"podpisane"`
	SignedAdult string `json:"podpisane_rodic"`
	Timestamp   Time   `json:"timestamp"`
	State       string `json:"stav"`
}

type Note struct {
	ID        string `json:"VcelickaID"`
	Date      string `json:"p_datum"`
	Text      string `json:"p_text"`
	Type      string `json:"p_typ"`
	SubjectID string `json:"PredmetID"`
}

type Event struct {
	ID string `json:"UdalostID"`
}

type Results struct {
	Grades map[string]Grade
	Events map[string]Event
	Notes  map[string]Note
}

func (dst *Results) Merge(src *Results) {
	maps.Copy(dst.Grades, src.Grades)
	maps.Copy(dst.Events, src.Events)
	maps.Copy(dst.Notes, src.Notes)
}

func ParseResults(jsondata []byte) (Results, error) {
	type RawGradesData struct {
		Grades []Grade                     `json:"vsetkyZnamky"`
		Events map[string]map[string]Event `json:"vsetkyUdalosti"`
		Notes  []Note                      `json:"vsetkyVcelicky"`
	}

	type RawGrades struct {
		Status string        `json:"status"`
		Data   RawGradesData `json:"data"`
	}

	var rgrades RawGrades
	var results Results
	err := json.Unmarshal(jsondata, &rgrades)
	if err != nil {
		return Results{}, err
	}

	results.Grades = make(map[string]Grade, len(rgrades.Data.Grades))

	for _, v := range rgrades.Data.Grades {
		results.Grades[v.ID] = v
	}

	results.Notes = make(map[string]Note, len(rgrades.Data.Notes))

	for _, v := range rgrades.Data.Notes {
		results.Notes[v.ID] = v
	}

	results.Events = make(map[string]Event, len(rgrades.Data.Events))

	for _, v := range rgrades.Data.Events["edupage"] {
		results.Events[v.ID] = v
	}

	return results, nil
}
