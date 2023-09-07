package model

import (
	"encoding/json"
	"strconv"
)

type User struct {
	Edubar           map[string]interface{} `json:"_edubar"`
	Timeline         []TimelineItem         `json:"items"` // Only recent timeline, see EdupageClient.Timeline for more
	DBI              DBI                    `json:"dbi"`
	UserRow          UserRow                `json:"userrow"`
	EventTypes       []EventType            `json:"eventtypes"`
	UserGroups       []string               `json:"usergroups"`
	DayPlan          DayPlan                `json:"dp"`
	NamedayToday     string                 `json:"meninyDnes"`
	NamedayTommorrow string                 `json:"meninyZajtra"`
}

type HeaderItem struct {
	SubjectID string      `json:"subjectid"`
	Changes   interface{} `json:"changes"`
}

type Plan struct {
	Period       string       `json:"uniperiod"`
	Type         string       `json:"type"`
	Header       []HeaderItem `json:"header"`
	SubjectID    string       `json:"subjectid"`
	ClassIDs     []string     `json:"classids"`
	GroupNames   []string     `json:"groupnames"`
	LID          string       `json:"lid"`
	TeacherIDs   []string     `json:"teacherids"`
	ClassroomIDs []string     `json:"classroomids"`
	StudentIDs   []string     `json:"studentids"`
	StartTime    string       `json:"starttime"`
	EndTime      string       `json:"endtime"`
}

type Date struct {
	Plans          []Plan        `json:"plan"`
	StudentAbsents []interface{} `json:"student_absents"` //TODO: check type
	Number         json.Number   `json:"tt_num"`
	Day            json.Number   `json:"tt_day"`
	Week           json.Number   `json:"tt_week"`
	Term           json.Number   `json:"tt_term"`
}

type DayPlan struct {
	Dates map[string]Date `json:"dates"`
}

type UserRow struct {
	UserID    string `json:"UserID"`
	StudentID string `json:"StudentID"`
	Firstname string `json:"p_meno"`
	Lastname  string `json:"p_priezvisko"`
	Email     string `json:"p_mail"`
	ClassID   string `json:"TriedaID"`
}

type EventType struct {
	ID            string   `json:"id"`
	C             string   `json:"c"`
	Name          string   `json:"name"`
	TTCancel      bool     `json:"ttcancel"`
	CTCan         bool     `json:"ctcan"`
	ClassRequired bool     `json:"classrequired"`
	Publish       string   `json:"publish"`
	NoCustomTime  bool     `json:"nocustomtime"`
	HideFields    []string `json:"hidefields"`
	DP            bool     `json:"dp"`
	Lesson        bool     `json:"lesson"`
	Attendance    bool     `json:"attendance"`
	DPrivacy      string   `json:"d_privacy"`
	CTEvent       bool     `json:"ctevent"`
	TemplateID    string   `json:"templateid"`
	CategoryID    string   `json:"categoryid"`
	SubID         string   `json:"subId"`
}

type DBI struct {
	Teachers           map[string]Teacher           `json:"teachers"`
	Classes            map[string]Class             `json:"classes"`
	Subjects           map[string]Subject           `json:"subjects"`
	Classrooms         map[string]Classroom         `json:"classrooms"`
	Students           map[string]Students          `json:"students"`
	Parents            map[string]Parents           `json:"parents"`
	Periods            map[json.Number]Period       `json:"periods"`
	DayParts           map[string]DayParts          `json:"dayparts"`
	AbsentTypes        map[string]AbsentType        `json:"absenttypes"`
	SubstitutionTypes  map[string]SubstitionType    `json:"substitutiontypes"`
	StudentAbsentTypes map[string]StudentAbsentType `json:"studentabsenttypes"`
	EventTypes         map[string]UserEventType     `json:"eventtypes"`
	ProcessTypes       map[string]ProcessType       `json:"processtypes"`
	ProcessStates      map[string]ProcessState      `json:"processstates"`
	IsStudentAdult     bool                         `json:"isstudentadult"`
}

type Teacher struct {
	ID          string `json:"id"`
	Firstname   string `json:"firstname"`
	Lastname    string `json:"lastname"`
	Short       string `json:"short"`
	Gender      string `json:"gender"`
	ClassroomID string `json:"classroomid"`
	DateFrom    string `json:"datefrom"`
	DateTo      string `json:"dateto"`
	IsOut       bool   `json:"isout"`
}

type Class struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Short       string `json:"short"`
	Grade       string `json:"grade"`
	TeacherID   string `json:"teacherid"`
	Teacher2ID  string `json:"teacher2id"`
	ClassroomID string `json:"classroomid"`
}

type Subject struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Short    string `json:"short"`
	CBHidden bool   `json:"cbhidden"`
}

type Classroom struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Short string `json:"short"`
}

type Students struct {
	ID            string `json:"id"`
	ClassroomID   string `json:"classroomid"`
	Firstname     string `json:"firstname"`
	Lastname      string `json:"lastname"`
	Parent1ID     string `json:"parent1id"`
	Parent2ID     string `json:"parent2id"`
	Parent3ID     string `json:"parent3id"`
	Gender        string `json:"gender"`
	DateFrom      string `json:"datefrom"`
	DateTo        string `json:"dateto"`
	NumberInClass string `json:"numberinclass"`
	IsOut         bool   `json:"isout"`
	Number        string `json:"number"`
	DataCopy      string `json:"kopiadata"`
}

type Parents struct {
	ID        string `json:"id"`
	Firstname string `json:"firstname"`
	Lastname  string `json:"lastname"`
	Gender    string `json:"gender"`
}

type Period struct {
	ID        string `json:"id"`
	StartTime string `json:"starttime"`
	EndTime   string `json:"endtime"`
	Name      string `json:"name"`
	Short     string `json:"short"`
}

type DayParts struct {
	ID        string `json:"id"`
	StartTime string `json:"starttime"`
	EndTime   string `json:"endtime"`
	Name      string `json:"name"`
	Short     string `json:"short"`
}

type AbsentType struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Short string `json:"short"`
}

type SubstitionType struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Short string `json:"short"`
}

type StudentAbsentType struct {
	ID         string `json:"id"`
	Name       string `json:"name"`
	Short      string `json:"short"`
	Color      string `json:"color"`
	ExcuseType string `json:"excusetype"`
}

type UserEventType struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Color string `json:"color"`
}

type ProcessType struct {
	ID           string   `json:"id"`
	User         string   `json:"user"`
	Name         string   `json:"name"`
	Workflow     string   `json:"workflow"`
	Enabled      bool     `json:"enabled"`
	TextOptional bool     `json:"textoptional"`
	DataColumns  []string `json:"datacolumns"`
}

type ProcessState struct {
	ID      string            `json:"id"`
	Name    string            `json:"name"`
	Icon    string            `json:"icon"`
	Color   string            `json:"color"`
	Next    map[string]string `json:"next"`
	Changes map[string]string `json:"changes"`
}

func (dbi *DBI) UnmarshalJSON(data []byte) error {
	type FakeDBI DBI
	var fdbi FakeDBI

	err := json.Unmarshal(data, &fdbi)
	if err != nil {
		type AltDBI struct {
			FakeDBI
			Periods []Period `json:"periods"`
		}

		var adbi AltDBI
		err := json.Unmarshal(data, &adbi)
		if err != nil {
			return err
		}

		dbi.AbsentTypes = adbi.AbsentTypes
		dbi.Classes = adbi.Classes
		dbi.Classrooms = adbi.Classrooms
		dbi.DayParts = adbi.DayParts
		dbi.EventTypes = adbi.EventTypes
		dbi.IsStudentAdult = adbi.IsStudentAdult

		dbi.Periods = make(map[json.Number]Period, len(adbi.Periods))
		for index, v := range adbi.Periods {
			dbi.Periods[json.Number(strconv.Itoa(index))] = v
		}

		dbi.Parents = adbi.Parents
		dbi.ProcessStates = adbi.ProcessStates
		dbi.ProcessTypes = adbi.ProcessTypes
		dbi.StudentAbsentTypes = adbi.StudentAbsentTypes
		dbi.Students = adbi.Students
		dbi.Subjects = adbi.Subjects
		dbi.SubstitutionTypes = adbi.SubstitutionTypes
		dbi.Teachers = adbi.Teachers
		dbi.IsStudentAdult = adbi.IsStudentAdult

	}

	dbi.AbsentTypes = fdbi.AbsentTypes
	dbi.Classes = fdbi.Classes
	dbi.Classrooms = fdbi.Classrooms
	dbi.DayParts = fdbi.DayParts
	dbi.EventTypes = fdbi.EventTypes
	dbi.IsStudentAdult = fdbi.IsStudentAdult
	dbi.Periods = fdbi.Periods
	dbi.Parents = fdbi.Parents
	dbi.ProcessStates = fdbi.ProcessStates
	dbi.ProcessTypes = fdbi.ProcessTypes
	dbi.StudentAbsentTypes = fdbi.StudentAbsentTypes
	dbi.Students = fdbi.Students
	dbi.Subjects = fdbi.Subjects
	dbi.SubstitutionTypes = fdbi.SubstitutionTypes
	dbi.Teachers = fdbi.Teachers
	dbi.IsStudentAdult = fdbi.IsStudentAdult
	return nil
}
