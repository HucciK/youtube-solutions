package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"yt-solutions-server/internal/core"
)

type UserService interface {
	AuthUser(userId int, username string) (*core.User, error)
	GetUserInfo(userId int) (*core.User, error)
	GetUserPhoto(userId int) ([]byte, error)
	ChargeBalance(userId, amount int) error
}

func (h Handler) AuthUser(w http.ResponseWriter, r *http.Request) {
	var newUser core.User

	if err := json.NewDecoder(r.Body).Decode(&newUser); err != nil {
		//
	}

	if newUser.ID == 0 || newUser.Name == "" {
		return
	}

	user, err := h.UserService.AuthUser(newUser.ID, newUser.Name)
	if err != nil {
		//
	}

	if err = json.NewEncoder(w).Encode(user); err != nil {
		//
	}
}

func (h Handler) GetUserInfo(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil {
		return
	}

	if id == 0 {
		return
	}

	user, err := h.UserService.GetUserInfo(id)
	if err != nil {
		return
	}

	if err = json.NewEncoder(w).Encode(user); err != nil {
		return
	}
}

func (h Handler) GetUserPhoto(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil {
		return
	}

	if id == 0 {
		return
	}

	photo, err := h.UserService.GetUserPhoto(id)
	if err != nil {
		return
	}

	if _, err = w.Write(photo); err != nil {
		//
	}
}

func (h Handler) ChargeBalance(w http.ResponseWriter, r *http.Request) {
	//TODO
	//authorization := r.Header.Get("Authorization")

	fmt.Println()

	id, err := strconv.Atoi(r.URL.Query().Get("EXTRA"))
	if err != nil {
		//
	}

	if id == 0 {
		//
		return
	}

	amount, err := strconv.Atoi(r.URL.Query().Get("AMOUNT"))
	if err != nil {
		//
	}

	fmt.Println(id, amount)

	if err := h.UserService.ChargeBalance(id, amount); err != nil {
		fmt.Println(err)
		//
	}
}
