package handler

import (
	"encoding/json"
	"net/http"
)

type MessageService interface {
	ProcessMessage(userId, messageId int, username, message string)
}

type CallbackService interface {
	ProcessCallback(userId, msgId int, username, callback, callbackId string)
}

type Handler struct {
	MessageService
	CallbackService
}

func NewHandler(c CallbackService, m MessageService) *Handler {
	return &Handler{
		CallbackService: c,
		MessageService:  m,
	}
}

func (h Handler) InitRoutes() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("/dashboardUpdates", h.ProcessUpdate)

	return mux
}

func (h Handler) ProcessUpdate(w http.ResponseWriter, r *http.Request) {

	var update Update

	if err := json.NewDecoder(r.Body).Decode(&update); err != nil {
		//
	}

	if update.Message.Text != "" {
		h.MessageService.ProcessMessage(update.Message.Chat.ID, update.Message.ID, update.Message.Chat.Username, update.Message.Text)
	}

	if update.Callback.Data != "" {
		h.CallbackService.ProcessCallback(update.Callback.Message.Chat.ID, update.Callback.Message.ID, update.Callback.Message.Chat.Username, update.Callback.Data, update.Callback.ID)
	}

}
