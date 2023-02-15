package repository

import (
	"database/sql"
	"fmt"
	"yt-solutions-server/internal/core"
)

type OrderRepo struct {
	db *sql.DB
}

func NewOrderRepo(db *sql.DB) *OrderRepo {
	return &OrderRepo{
		db: db,
	}
}

func (o OrderRepo) ProcessOrder(userId, price int, key *core.Key) error {
	tx, err := o.db.Begin()
	if err != nil {
		return err
	}

	if _, err = tx.Exec("UPDATE users SET balance = balance - $1 WHERE id=$2", price, userId); err != nil {
		tx.Rollback()
		return err
	}

	if _, err = tx.Exec("INSERT INTO keys(key, type, expire, is_expired, owner) VALUES ($1, $2, $3, $4, $5)", key.Key, key.Type, key.Expire, key.IsExpired, key.Owner); err != nil {
		tx.Rollback()
		return err
	}

	if _, err = tx.Exec("UPDATE users SET has_key=$1 WHERE  id=$2", true, userId); err != nil {
		return err
	}

	if err = tx.Commit(); err != nil {
		tx.Rollback()
		return err
	}

	return nil
}

func (o OrderRepo) ProcessRenewal(ownerId, renewal int, key, expiration string) error {
	fmt.Println(ownerId, renewal, key, expiration)

	tx, err := o.db.Begin()
	if err != nil {
		return err
	}

	if _, err = tx.Exec("UPDATE users SET balance = balance - $1 WHERE id=$2", renewal, ownerId); err != nil {
		return err
	}

	if _, err = tx.Exec("UPDATE keys SET expire=$1 WHERE key=$2", expiration, key); err != nil {
		tx.Rollback()
		return err
	}

	if _, err = tx.Exec("UPDATE keys SET is_expired=$1 WHERE key=$2", false, key); err != nil {
		tx.Rollback()
		return err
	}

	if err = tx.Commit(); err != nil {
		tx.Rollback()
		return err
	}

	return nil
}

func (o OrderRepo) saveNewKey(key *core.Key) error {
	if _, err := o.db.Query("INSERT INTO keys(key, type, expire, is_expired, owner)", key.Key, key.Type, key.Expire, key.IsExpired, key.Owner); err != nil {
		return err
	}

	return nil
}
