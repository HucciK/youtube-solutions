package handler

type Update struct {
	Message  Message  `json:"message"`
	Callback Callback `json:"callback_query"`
}

type Message struct {
	ID   int    `json:"message_id"`
	Chat Chat   `json:"chat"`
	Text string `json:"text"`
}

type Chat struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
}

type Callback struct {
	ID      string  `json:"id"`
	Message Message `json:"message"`
	Data    string  `json:"data"`
}
