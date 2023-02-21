package models

const (
	FoundChannelType = "found channel"
	SavingValidType  = "saving valid"
	ErrorNotifyType  = "error notify"
	CheckStatusType  = "check status"
)

type Update struct {
	Type    string      `json:"type"`
	Data    interface{} `json:"data"`
	Err     interface{} `json:"error"`
	Channel *chan Update
}

type CheckedAmount struct {
	Valid   int `json:"valid"`
	Errors  int `json:"errors"`
	Checked int `json:"checked"`
}

func (u *Update) Send() {
	*u.Channel <- *u
}
