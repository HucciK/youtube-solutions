package repository

import (
	"database/sql"
	"errors"
	"yt-solutions-server/internal/core"
)

type KeyRepo struct {
	db *sql.DB
}

func NewKeyRepo(db *sql.DB) *KeyRepo {
	return &KeyRepo{
		db: db,
	}
}

func (k KeyRepo) GetKey(key string) (*core.Key, error) {
	var keyInfo core.Key

	res, err := k.db.Query("SELECT * FROM keys WHERE key=$1", key)
	if err != nil {
		return nil, err
	}

	for res.Next() {
		if err = res.Scan(&keyInfo.Key, &keyInfo.Type, &keyInfo.Ip, &keyInfo.Expire, &keyInfo.IsExpired, &keyInfo.Owner); err != nil {
			return nil, err
		}
	}

	if keyInfo.Key == "" {
		return nil, errors.New("key doesnt exists")
	}

	return &keyInfo, nil
}

func (k KeyRepo) GetKeyByOwnerId(ownerId int) (*core.Key, error) {
	var key core.Key

	res, err := k.db.Query("SELECT * FROM keys WHERE owner=$1", ownerId)
	if err != nil {
		return nil, err
	}

	for res.Next() {
		if err = res.Scan(&key.Key, &key.Type, &key.Ip, &key.Expire, &key.IsExpired, &key.Owner); err != nil {
			return nil, err
		}
	}

	return &key, nil
}

func (k KeyRepo) GetAllFree() ([]*core.Key, error) {
	var keys []*core.Key

	res, err := k.db.Query("SELECT * FROM keys WHERE type=$1", "free")
	if err != nil {
		return nil, err
	}

	for res.Next() {
		var key core.Key
		if err = res.Scan(&key.Key, &key.Type, &key.Ip, &key.Expire, &key.IsExpired, &key.Owner); err != nil {
			return nil, err
		}
		keys = append(keys, &key)
	}

	return keys, nil
}

func (k KeyRepo) GetAllLifetime() ([]*core.Key, error) {
	var keys []*core.Key

	res, err := k.db.Query("SELECT * FROM keys WHERE type=$1", "lifetime")
	if err != nil {
		return nil, err
	}

	for res.Next() {
		var key core.Key
		if err = res.Scan(&key.Key, &key.Type, &key.Ip, &key.Expire, &key.IsExpired, &key.Owner); err != nil {
			return nil, err
		}
		keys = append(keys, &key)
	}

	return keys, nil
}

func (k KeyRepo) SetKeyExpired(key string, expired bool) error {
	if _, err := k.db.Exec("UPDATE keys SET is_expired=$1 WHERE key=$2", expired, key); err != nil {
		return err
	}

	return nil
}

func (k KeyRepo) UpdateIP(key, ip string) error {
	if _, err := k.db.Exec("UPDATE keys SET ip=$1 WHERE key=$2", ip, key); err != nil {
		return err
	}

	return nil
}

func (k KeyRepo) UnbindAddress(userId int) error {
	if _, err := k.db.Exec("UPDATE keys SET ip='' WHERE owner=$1", userId); err != nil {
		return err
	}

	return nil
}
