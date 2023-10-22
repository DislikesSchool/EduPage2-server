package edupage

import (
	"crypto/sha1"
	"encoding/base64"
	"net/url"
)

type MessagePayload struct {
	SelectedUser string
	Text         string
	Attachments  string
	Typ          string
}

type CanteenPayload struct {
	BoarderID   string            `json:"stravnikid"`
	Edupage     string            `json:"edupage"`
	BoarderUser string            `json:"stravnikUser,omitempty"`
	Date        string            `json:"mysqlDate"` // YYYY-mm-dd
	FIDS        map[string]string `json:"jids"`
	View        string            `json:"view"`
	Permission  string            `json:"pravo"`
	Action      string            `json:"akcia"`
}

func CreateMessage(receiver, text, attachments string) MessagePayload {
	return MessagePayload{
		SelectedUser: receiver,
		Text:         text,
		Attachments:  attachments,
		Typ:          "sprava",
	}
}

func CreatePayload(data map[string]string) url.Values {
	payload_values := url.Values{}
	for key, val := range data {
		payload_values.Add(key, val)
	}

	payload := payload_values.Encode()
	encoded := make([]byte, base64.URLEncoding.EncodedLen(len(payload)))
	base64.URLEncoding.Encode(encoded, []byte(payload))

	values := url.Values{}
	values.Add("eqap", string(encoded))

	hasher := sha1.New()
	hasher.Reset()
	hasher.Write(encoded)
	values.Add("eqacs", base64.URLEncoding.EncodeToString(hasher.Sum(nil)))

	values.Add("eqaz", "1")
	return values
}
