package resources

type SendMessageRequest struct {
	ChatId   int         `json:"chat_id"`
	Text     string      `json:"text"`
	Keyboard interface{} `json:"reply_markup,omitempty"`
}

type SendMessageResult struct {
	Success bool    `json:"ok"`
	Result  Message `json:"result"`
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

type EditMessageRequest struct {
	ChatId    int         `json:"chat_id"`
	MessageId int         `json:"message_id"`
	Text      string      `json:"text"`
	Keyboard  interface{} `json:"reply_markup"`
}

type DeleteMessageRequest struct {
	ChatId    int `json:"chat_id"`
	MessageId int `json:"message_id"`
}

type AnswerCallbackRequest struct {
	ID        string `json:"callback_query_id"`
	Text      string `json:"text"`
	ShowAlert bool   `json:"show_alert"`
	CacheTime int    `json:"cache_time"`
}
