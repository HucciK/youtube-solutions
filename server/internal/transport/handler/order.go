package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"yt-solutions-server/internal/core"
)

type OrderService interface {
	GetAvailability() (int, int, int, error)
	GetRenewalInfo() int
	ProcessOrder(userId int) (bool, *core.Key, error)
	ProcessRenewal(ownerId int) (bool, error)
}

func (h Handler) GetAvailability(w http.ResponseWriter, r *http.Request) {
	free, lifetime, price, err := h.OrderService.GetAvailability()
	if err != nil {
		//
	}

	resp := AvailabilityResponse{
		Free:     free,
		Lifetime: lifetime,
		Price:    price,
	}

	if err = json.NewEncoder(w).Encode(&resp); err != nil {
		//
	}
}

func (h Handler) GetRenewalInfo(w http.ResponseWriter, r *http.Request) {
	renewal := h.OrderService.GetRenewalInfo()
	resp := RenewalInfoResponse{
		Renewal: renewal,
	}

	if err := json.NewEncoder(w).Encode(&resp); err != nil {
		//
	}
}

func (h Handler) ProcessOrder(w http.ResponseWriter, r *http.Request) {

	userId, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil {
		return
	}

	if userId == 0 {
		//
		fmt.Println(userId)
		return
	}

	success, key, err := h.OrderService.ProcessOrder(userId)
	if err != nil {
		fmt.Println(err)
		return
	}

	if !success {
		resp := ProcessOrderResponse{
			Success: success,
			Key:     key,
		}

		if err = json.NewEncoder(w).Encode(&resp); err != nil {
			return
		}

		return
	}

	resp := ProcessOrderResponse{
		Success: success,
		Key:     key,
	}

	if err = json.NewEncoder(w).Encode(&resp); err != nil {
		return
	}
}

func (h Handler) ProcessRenewal(w http.ResponseWriter, r *http.Request) {
	ownerId, err := strconv.Atoi(r.URL.Query().Get("owner"))
	if err != nil {
		//
	}

	success, err := h.OrderService.ProcessRenewal(ownerId)
	if err != nil {
		//
	}

	resp := ProcessRenewalResponse{
		Success: success,
	}

	if err = json.NewEncoder(w).Encode(&resp); err != nil {
		//
	}
}
