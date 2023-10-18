package apimodel

import (
	"encoding/json"

	"github.com/DislikesSchool/EduPage2-server/edupage/model"
)

type TimelineItemReduced struct {
	ID            string                 `json:"timelineid"`
	Timestamp     model.Time             `json:"timestamp"`
	ReactionTo    string                 `json:"reakcia_na"`
	Type          string                 `json:"typ"`
	User          string                 `json:"user"`
	TargetUser    string                 `json:"target_user"`
	Text          string                 `json:"text"`
	Data          model.StringJsonObject `json:"data"`
	Owner         string                 `json:"vlastnik"`
	ReactionCount int                    `json:"poct_reakcii"`
	Removed       json.Number            `json:"removed"`
}

type TimelineItemWithOrigin struct {
	ID              string                 `json:"timelineid"`
	Timestamp       model.Time             `json:"timestamp"`
	ReactionTo      string                 `json:"reakcia_na"`
	Type            string                 `json:"typ"`
	User            string                 `json:"user"`
	TargetUser      string                 `json:"target_user"`
	UserName        string                 `json:"user_meno"`
	OtherID         string                 `json:"ineid"`
	Text            string                 `json:"text"`
	TimeAdded       model.Time             `json:"cas_pridania"`
	TimeEvent       model.Time             `json:"cas_udalosti"`
	Data            model.StringJsonObject `json:"data"`
	Owner           string                 `json:"vlastnik"`
	OwnerName       string                 `json:"vlastnik_meno"`
	ReactionCount   int                    `json:"poct_reakcii"`
	LastReaction    string                 `json:"posledna_reakcia"`
	PomocnyZaznam   string                 `json:"pomocny_zaznam"`
	Removed         json.Number            `json:"removed"`
	TimeAddedBTC    model.Time             `json:"cas_pridania_btc"`
	LastReactionBTC model.Time             `json:"cas_udalosti_btc"`
	OriginServer    string                 `json:"origin_server"`
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
	Type             string                 `json:"typ"`
	LikeCount        json.Number            `json:"pocet_like"`
	ReactionCount    json.Number            `json:"pocet_reakcii"`
	DoneCount        json.Number            `json:"pocet_done"`
	State            string                 `json:"stav"`
	LastResult       string                 `json:"posledny_vysledok"`
	GradeEventID     interface{}            `json:"znamky_udalostid"`
	StudentsHidden   string                 `json:"students_hidden"`
	Data             model.StringJsonObject `json:"data"`
	EvaluationStatus string                 `json:"stavhodnotetimelinePathd"`
	Attachments      interface{}            `json:"attachements"`
}

type Timeline struct {
	Homeworks map[string]HomeworkReduced
	Items     map[string]TimelineItemReduced
}

type TimelineRequest struct {
	From string `json:"from" example:"2022-01-01T00:00:00Z"`
	To   string `json:"to" example:"2022-01-01T00:00:00Z" default:"time.Now()"`
}
