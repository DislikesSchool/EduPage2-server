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

func CreateMessage(receiver, text, attachments string) MessagePayload {
	return MessagePayload{
		SelectedUser: receiver,
		Text:         text,
		Attachments:  attachments,
		Typ:          "sprava",
	}
}

func CreatePayload(payload map[string]string) (url.Values, error) {
	data := url.Values{}
	for key, val := range payload {
		data.Add(key, val)
	}

	es := data.Encode()
	e := make([]byte, base64.URLEncoding.EncodedLen(len(es)))
	base64.URLEncoding.Encode(e, []byte(es))

	r := url.Values{}
	r.Add("eqap", string(e))

	hasher := sha1.New()
	hasher.Reset()
	hasher.Write(e)
	r.Add("eqacs", base64.URLEncoding.EncodeToString(hasher.Sum(nil)))

	r.Add("eqaz", "1")
	return r, nil
}
