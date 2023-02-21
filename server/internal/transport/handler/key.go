package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"yt-solutions-server/internal/core"
)

type KeyService interface {
	GetKeyInfoByOwnerId(ownerId int) (*core.Key, error)
	GetKeyInfo(key, ip string) (*core.Key, error)
	UnbindAddress(userId int) error
}

func (h Handler) GetKeyInfo(w http.ResponseWriter, r *http.Request) {
	var key *core.Key

	method := r.URL.Query().Get("by")

	if method == "ownerId" {

		ownerId, err := strconv.Atoi(r.URL.Query().Get("owner"))
		if err != nil {
			//
		}

		key, err = h.KeyService.GetKeyInfoByOwnerId(ownerId)
		if err != nil {
			//
		}
	}

	if method == "key" {

	}

	if err := json.NewEncoder(w).Encode(&key); err != nil {
		//
	}
}

func (h Handler) AuthKey(w http.ResponseWriter, r *http.Request) {
	var valid bool

	key := r.Header.Get("Authorization")
	ip := r.Header.Get("X-Forwarded-For")

	if key == "" {
		resp := KeyAuthResponse{
			Valid:         false,
			LatestVersion: h.AppVersion,
		}

		if err := json.NewEncoder(w).Encode(&resp); err != nil {
			//
		}
		return
	}

	keyInfo, err := h.KeyService.GetKeyInfo(key, ip)
	if err != nil {
		fmt.Println(err)
		resp := KeyAuthResponse{
			Valid:         false,
			LatestVersion: h.AppVersion,
		}

		if err = json.NewEncoder(w).Encode(&resp); err != nil {
			//
		}
		return
	}

	if !keyInfo.IsExpired {
		valid = true
	}

	resp := KeyAuthResponse{
		Valid:         valid,
		LatestVersion: h.AppVersion,
		Key:           keyInfo,
	}

	if err = json.NewEncoder(w).Encode(&resp); err != nil {
		//
	}
}

func (h Handler) UnbindAddress(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil {
		return
	}

	if id == 0 {
		return
	}

	if err = h.KeyService.UnbindAddress(id); err != nil {
		fmt.Println(err)
		return
	}
}
