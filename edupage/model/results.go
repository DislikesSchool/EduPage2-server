package model

import (
	"encoding/json"

	"golang.org/x/exp/maps"
)

type Event struct {
	Provider     string      `json:"provider"`
	ID           string      `json:"znamkaid"`
	StudentID    string      `json:"studentid"`
	SubjectID    string      `json:"predmetid"`
	EventID      string      `json:"udalostID"`
	Month        string      `json:"mesiac"`
	Data         string      `json:"data"`
	Date         string      `json:"datum"`
	TeacherID    string      `json:"ucitelid"`
	Signed       string      `json:"podpisane"`
	SignedAdult  string      `json:"podpisane_rodic"`
	Timestamp    string      `json:"timestamp"`
	State        string      `json:"stav"`
	Color        string      `json:"p_farba"`
	EventName    string      `json:"p_meno"`
	FirstAverage string      `json:"p_najskor_priemer"`
	EventType    interface{} `json:"p_typ_udalosti"`
	Weight       interface{} `json:"p_vaha"`
	ClassID      string      `json:"TriedaID"`
	PlanID       string      `json:"planid"`
	GradeCount   interface{} `json:"p_pocet_znamok"`
	MoreData     interface{} `json:"moredata"`
	Average      string      `json:"priemer"`
}

type Note struct {
	ID        string `json:"VcelickaID"`
	Date      string `json:"p_datum"`
	Text      string `json:"p_text"`
	Type      string `json:"p_typ"`
	SubjectID string `json:"PredmetID"`
}

type Grade struct {
	Provider    string `json:"provider"`
	ID          string `json:"udalostid"`
	GradeID     string `json:"znamkaid"`
	StudentID   string `json:"studentid"`
	SubjectID   string `json:"predmetid"`
	Month       string `json:"mesiac"`
	Data        string `json:"data"`
	Date        string `json:"datum"`
	TeacherID   string `json:"ucitelid"`
	Signed      string `json:"podpisane"`
	SignedAdult string `json:"podpisane_rodic"`
	Timestamp   string `json:"timestamp"`
	State       string `json:"stav"`
}

type Results struct {
	Events map[string]Event
	Notes  map[string]Note
}

func (dst *Results) Merge(src *Results) {
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

	results.Notes = make(map[string]Note, len(rgrades.Data.Notes))

	for _, v := range rgrades.Data.Notes {
		results.Notes[v.ID] = v
	}

	results.Events = make(map[string]Event, len(rgrades.Data.Events["edupage"]))

	for k, v := range rgrades.Data.Events["edupage"] {
		results.Events[k] = v
	}

	for _, v := range rgrades.Data.Grades {
		if event, ok := results.Events[v.ID]; ok {
			event.ID = v.ID
			event.StudentID = v.StudentID
			event.Data = v.Data
			event.Date = v.Date
			event.TeacherID = v.TeacherID
			event.Signed = v.Signed
			event.SignedAdult = v.SignedAdult
			event.Timestamp = v.Timestamp
			event.State = v.State
			results.Events[v.ID] = event
		}
	}

	return results, nil
}
