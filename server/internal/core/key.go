package core

import (
	"crypto/sha256"
	"fmt"
	"time"
)

const (
	TypeDev      = "dev"
	TypeFree     = "free"
	TypeLifetime = "lifetime"
	TypeRenew    = "renewal"
)

type Key struct {
	Key       string `json:"key"`
	Type      string `json:"type"`
	Ip        string `json:"ip"`
	Expire    string `json:"expire"`
	IsExpired bool   `json:"is_expired"`
	Owner     int    `json:"owner"`
}

func (k *Key) GenerateKey(user *User) {
	hash := sha256.New()
	hash.Write([]byte(fmt.Sprintf("%s_%d_%d", user.Name, user.ID, time.Now().UnixNano())))

	k.Key = fmt.Sprintf("%x", hash.Sum(nil))
}

func (k *Key) SelectKeyType(free, lifetime int) {
	if free > 0 {
		k.Type = TypeFree
		return
	}

	if lifetime > 0 {
		k.Type = TypeLifetime
		return
	}

	k.Type = TypeRenew
}

func (k *Key) SetExpireDate() {
	if k.Type == TypeFree || k.Type == TypeLifetime || k.Type == TypeDev {
		k.Expire = "lifetime"
		return
	}

	k.Expire = time.Now().AddDate(0, 1, 0).Format("02-01-2006")
}

func (k *Key) SetOwner(owner int) {
	k.Owner = owner
}

func (k *Key) CheckExpiration() error {
	if k.Type == TypeFree || k.Type == TypeLifetime || k.Type == TypeDev {
		return nil
	}

	expireParsed, err := time.Parse("02-01-2006", k.Expire)
	if err != nil {
		return err
	}
	fmt.Println(expireParsed, time.Now())

	isExpired := expireParsed.Before(time.Now())
	if isExpired {
		k.IsExpired = isExpired
	}
	return nil
}

func (k *Key) UpdateExpiration() {
	k.Expire = time.Now().AddDate(0, 1, 0).Format("02-01-2006")
}
