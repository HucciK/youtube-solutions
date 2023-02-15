package models

const (
	FoundChannelType = "found channel"
	SavingValidType  = "saving valid"
	ErrorNotifyType  = "error notify"
	CheckStatusType  = "check status"
)

type Update struct {
	Type string      `json:"type"`
	Data interface{} `json:"data"`
	Err  error       `json:"error"`
}

type CheckedAmount struct {
	Valid   int `json:"valid"`
	Checked int `json:"checked"`
}
