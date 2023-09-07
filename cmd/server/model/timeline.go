package main

import (
	"encoding/json"

	"github.com/DislikesSchool/EduPage2-server/edupage/model"
)

type RecentTimelineSuccessResponse struct {
	Success  bool       `json:"success" example:"true"`
	Error    string     `json:"error" example:""`
	Timeline []Timeline `json:"timeline"`
}

type RecentTimelineUnauthorizedResponse struct {
	Success bool   `json:"success" example:"false"`
	Error   string `json:"error" example:"Unauthorized"`
}

type RecentTimelineInternalErrorResponse struct {
	Success bool   `json:"success" example:"false"`
	Error   string `json:"error" example:"failed to create payload"`
}

type TimelineItemReduced struct {
	ID            string                 `json:"timelineid"`
	Timestamp     model.Time             `json:"timestamp"`
	ReactionTo    string                 `json:"reakcia_na"`
	Type          model.TimelineItemType `json:"typ"`
	User          string                 `json:"user"`
	TargetUser    string                 `json:"target_user"`
	Text          string                 `json:"text"`
	Data          model.TimelineItemData `json:"data"`
	Owner         string                 `json:"vlastnik"`
	ReactionCount int                    `json:"poct_reakcii"`
	Removed       json.Number            `json:"removed"`
}

type HomeworkReduced struct {
	ID               string                 `json:"hwkid"`
	HomeworkID       string                 `json:"homeworkid"`
	UserID           string                 `json:"userid"`
	LessonID         json.Number            `json:"predmetid"`
	Name             string                 `json:"name"`
	Details          string                 `json:"details"`
	DateCreated      string                 `json:"datecreated"`
	Period           interface{}            `json:"period"`
	Timestamp        string                 `json:"timestamp"`
	TestID           string                 `json:"testid"`
	Type             model.TimelineItemType `json:"typ"`
	LikeCount        json.Number            `json:"pocet_like"`
	ReactionCount    json.Number            `json:"pocet_reakcii"`
	DoneCount        json.Number            `json:"pocet_done"`
	State            string                 `json:"stav"`
	LastResult       string                 `json:"posledny_vysledok"`
	GradeEventID     interface{}            `json:"znamky_udalostid"`
	StudentsHidden   string                 `json:"students_hidden"`
	Data             model.TimelineItemData `json:"data"`
	EvaluationStatus string                 `json:"stavhodnotetimelinePathd"`
	Attachments      interface{}            `json:"attachements"`
}

type Timeline struct {
	Homeworks map[string]HomeworkReduced
	Items     map[string]TimelineItemReduced
}

type TimelineRequest struct {
	From string `json:"from" example:"2022-01-01T00:00:00Z2022-01-01T00:00:00Z"`
	To   string `json:"to" example:"2022-01-01T00:00:00Z" default:"time.Now()"`
}

type TimelineSuccessResponse struct {
	Success  bool       `json:"success" example:"true"`
	Error    string     `json:"error" example:""`
	Timeline []Timeline `json:"timeline"`
	From     string     `json:"from"`
	To       string     `json:"to"`
}

type TimelineUnauthorizedResponse struct {
	Success bool   `json:"success" example:"false"`
	Error   string `json:"error" example:"Unauthorized"`
}

type TimelineInternalErrorResponse struct {
	Success bool   `json:"success" example:"false"`
	Error   string `json:"error" example:"failed to create payload"`
}
