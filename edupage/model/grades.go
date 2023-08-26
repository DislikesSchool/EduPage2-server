package model

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
