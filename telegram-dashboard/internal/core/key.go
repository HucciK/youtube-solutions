package core

import (
	"fmt"
	"strings"
)

const (
	TypeFree = ""
)

type Key struct {
	Key       string `json:"key"`
	Type      string `json:"type"`
	Ip        string `json:"ip"`
	Expire    string `json:"expire"`
	IsExpired bool   `json:"is_expired"`
	Owner     int    `json:"owner"`
}

func (k *Key) Info() string {
	if k.Type == "free" || k.Type == "lifetime" {
		return fmt.Sprintf("Ваш ключ: %s\n\nТип ключа: %s", k.Key, strings.ToTitle(k.Type))
	}

	active := "Нет"
	if !k.IsExpired {
		active = "Да"
	}

	return fmt.Sprintf("Ваш ключ: %s\n\nТип ключа: %s\nИстекает: %s\nАктивен: %s", k.Key, k.Type, k.Expire, active)
}
