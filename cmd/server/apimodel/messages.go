package apimodel

import "github.com/DislikesSchool/EduPage2-server/edupage"

type Recipient struct {
	ID   string `json:"id"`
	Type string `json:"type"`
	Name string `json:"name"`
}

type SendMessageRequest struct {
	Recipient string                 `json:"recipient"`
	Message   edupage.MessageOptions `json:"message"`
}
