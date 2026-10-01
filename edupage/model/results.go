package model

import (
	"bytes"
	"encoding/json"
	"fmt"

	"golang.org/x/exp/maps"
)

// flexList accepts either a JSON array or a JSON object keyed by ID
// (EduPage switched several collections from arrays to keyed objects).
type flexList[T any] []T

func (l *flexList[T]) UnmarshalJSON(data []byte) error {
	t := bytes.TrimSpace(data)
	if len(t) == 0 || string(t) == "null" {
		*l = nil
		return nil
	}
	if t[0] == '[' {
		var arr []T
		if err := json.Unmarshal(data, &arr); err != nil {
			return err
		}
		*l = arr
		return nil
	}
	if t[0] == '{' {
		if len(t) == 2 { // {}
			*l = nil
			return nil
		}
		var m map[string]T
		if err := json.Unmarshal(data, &m); err != nil {
			return err
		}
		arr := make([]T, 0, len(m))
		for _, v := range m {
			arr = append(arr, v)
		}
		*l = arr
		return nil
	}
	return fmt.Errorf("expected array or object, got %.30s", string(t))
}

type Event struct {
	Provider     FlexString  `json:"provider"`
	ID           FlexString  `json:"znamkaid"`
	StudentID    FlexString  `json:"studentid"`
	SubjectID    FlexString  `json:"predmetid"`
	EventID      FlexString  `json:"udalostID"`
	Month        FlexString  `json:"mesiac"`
	Data         FlexString  `json:"data"`
	Date         FlexString  `json:"datum"`
	TeacherID    FlexString  `json:"ucitelid"`
	Signed       FlexString  `json:"podpisane"`
	SignedAdult  FlexString  `json:"podpisane_rodic"`
	Timestamp    FlexString  `json:"timestamp"`
	State        FlexString  `json:"stav"`
	Color        FlexString  `json:"p_farba"`
	EventName    FlexString  `json:"p_meno"`
	FirstAverage FlexString  `json:"p_najskor_priemer"`
	EventType    interface{} `json:"p_typ_udalosti"`
	Weight       interface{} `json:"p_vaha"`
	ClassID      FlexString  `json:"TriedaID"`
	PlanID       FlexString  `json:"planid"`
	GradeCount   interface{} `json:"p_pocet_znamok"`
	MoreData     interface{} `json:"moredata"`
	Average      FlexString  `json:"priemer"`
}

type Note struct {
	ID        FlexString `json:"VcelickaID"`
	Date      FlexString `json:"p_datum"`
	Text      FlexString `json:"p_text"`
	Type      FlexString `json:"p_typ"`
	SubjectID FlexString `json:"PredmetID"`
}

type Grade struct {
	Provider    FlexString `json:"provider"`
	ID          FlexString `json:"udalostid"`
	GradeID     FlexString `json:"znamkaid"`
	StudentID   FlexString `json:"studentid"`
	SubjectID   FlexString `json:"predmetid"`
	Month       FlexString `json:"mesiac"`
	Data        FlexString `json:"data"`
	Date        FlexString `json:"datum"`
	TeacherID   FlexString `json:"ucitelid"`
	Signed      FlexString `json:"podpisane"`
	SignedAdult FlexString `json:"podpisane_rodic"`
	Timestamp   FlexString `json:"timestamp"`
	State       FlexString `json:"stav"`
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
		Grades flexList[Grade]            `json:"vsetkyZnamky"`
		Events map[string]flexList[Event] `json:"vsetkyUdalosti"`
		Notes  flexList[Note]             `json:"vsetkyVcelicky"`
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
		results.Notes[string(v.ID)] = v
	}

	results.Events = make(map[string]Event, len(rgrades.Data.Events["edupage"]))

	for _, v := range rgrades.Data.Events["edupage"] {
		results.Events[string(v.EventID)] = v
	}

	for _, v := range rgrades.Data.Grades {
		if event, ok := results.Events[string(v.ID)]; ok {
			event.ID = v.ID
			event.StudentID = v.StudentID
			event.Data = v.Data
			event.Date = v.Date
			event.TeacherID = v.TeacherID
			event.Signed = v.Signed
			event.SignedAdult = v.SignedAdult
			event.Timestamp = v.Timestamp
			event.State = v.State
			results.Events[string(v.ID)] = event
		}
	}

	return results, nil
}
