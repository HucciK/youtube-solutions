package handler

import "net/http"

type Handler struct {
	UserService
	KeyService
	OrderService
	AppVersion string
}

func NewHandler(u UserService, k KeyService, o OrderService, v string) *Handler {
	return &Handler{
		UserService:  u,
		KeyService:   k,
		OrderService: o,
		AppVersion:   v,
	}
}

func (h Handler) InitRoutes() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("/api/authUser", h.AuthUser)
	mux.HandleFunc("/api/getUserInfo", h.GetUserInfo)
	mux.HandleFunc("/api/getUserPhoto", h.GetUserPhoto)

	mux.HandleFunc("/api/getKeyInfo", h.GetKeyInfo)
	mux.HandleFunc("/api/authKey", h.AuthKey)
	mux.HandleFunc("/api/unbindAddress", h.UnbindAddress)

	mux.HandleFunc("/api/getAvailability", h.GetAvailability)
	mux.HandleFunc("/api/getRenewalInfo", h.GetRenewalInfo)

	mux.HandleFunc("/api/processOrder", h.ProcessOrder)
	mux.HandleFunc("/api/processRenewal", h.ProcessRenewal)
	mux.HandleFunc("/api/chargeBalance", h.ChargeBalance)

	mux.HandleFunc("/clientUpdate/", h.ClientUpdate)
	mux.HandleFunc("/devUpdate/", h.DevUpdate)

	return mux
}
