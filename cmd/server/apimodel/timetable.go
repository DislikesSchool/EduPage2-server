package apimodel

import "github.com/DislikesSchool/EduPage2-server/edupage/model"

type TimetableRequest struct {
	From string `json:"from" example:"2022-01-01T00:00:00Z"`
	To   string `json:"to" example:"2022-01-01T00:00:00Z" default:"time.Now()"`
}

type CompleteTimetable struct {
	Days map[string][]CompleteTimetableItem // key format is YYYY-MM-dd or 2006-01-02
}

type CompleteTimetableItem struct {
	Type       string            `json:"type"`
	Date       string            `json:"date"`
	Period     string            `json:"uniperiod"`
	StartTime  string            `json:"starttime"`
	EndTime    string            `json:"endtime"`
	Subject    model.Subject     `json:"subject"`
	Classes    []model.Class     `json:"classes"`
	GroupNames []string          `json:"groupnames"`
	IGroupID   string            `json:"igroupid"`
	Teachers   []model.Teacher   `json:"teachers"`
	Classrooms []model.Classroom `json:"classrooms"`
	StudentIDs []string          `json:"studentids"`
	Colors     []string          `json:"colors"`
}
