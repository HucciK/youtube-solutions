package services

import (
	"errors"
	"strings"
	"yt-solutions-server/internal/core"
)

type KeyRepo interface {
	GetKey(key string) (*core.Key, error)
	GetKeyByOwnerId(ownerId int) (*core.Key, error)
	GetAllFree() ([]*core.Key, error)
	GetAllLifetime() ([]*core.Key, error)

	SetKeyExpired(key string, expired bool) error
	UpdateIP(key, ip string) error
	UnbindAddress(userId int) error
}

type KeyService struct {
	KeyRepo
}

func NewKeyService(k KeyRepo) *KeyService {
	return &KeyService{
		KeyRepo: k,
	}
}

func (k KeyService) GetKeyInfoByOwnerId(ownerId int) (*core.Key, error) {
	key, err := k.KeyRepo.GetKeyByOwnerId(ownerId)
	if err != nil {
		return nil, err
	}

	if err = key.CheckExpiration(); err != nil {
		return nil, err
	}

	if key.IsExpired {
		if err = k.KeyRepo.SetKeyExpired(key.Key, true); err != nil {
			return nil, err
		}
	}

	return key, nil
}

func (k KeyService) GetKeyInfo(key, ip string) (*core.Key, error) {

	keyInfo, err := k.KeyRepo.GetKey(key)
	if err != nil {
		return nil, err
	}

	if keyInfo.Ip == "" {
		keyInfo.Ip = ip

		if err = k.KeyRepo.UpdateIP(key, ip); err != nil {
			return nil, err
		}
	}

	if k.GetTwoOctets(keyInfo.Ip) != k.GetTwoOctets(ip) {
		return nil, errors.New("invalid ip")
	}

	if err = keyInfo.CheckExpiration(); err != nil {
		return nil, err
	}

	if keyInfo.IsExpired {
		if err = k.KeyRepo.SetKeyExpired(keyInfo.Key, true); err != nil {
			return nil, err
		}
	}

	return keyInfo, nil
}

func (k KeyService) GetTwoOctets(ip string) string {
	return strings.Join(strings.Split(ip, ".")[:2], ".")
}

func (k KeyService) UnbindAddress(userId int) error {
	return k.KeyRepo.UnbindAddress(userId)
}
