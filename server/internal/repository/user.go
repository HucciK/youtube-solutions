package repository

import (
	"database/sql"
	"yt-solutions-server/internal/core"
)

type UserRepo struct {
	db *sql.DB
}

func NewUserRepo(db *sql.DB) *UserRepo {
	return &UserRepo{
		db: db,
	}
}

func (u UserRepo) CreateUser(userId int, username string) error {
	if _, err := u.db.Exec("INSERT  INTO users(id, name) VALUES($1, $2)", userId, username); err != nil {
		return err
	}

	return nil
}

func (u UserRepo) GetUserById(userId int) (*core.User, error) {
	var user core.User

	res, err := u.db.Query("SELECT * FROM users WHERE id=$1", userId)
	if err != nil {
		return nil, err
	}

	for res.Next() {
		if err = res.Scan(&user.ID, &user.Name, &user.Balance, &user.HasKey); err != nil {
			return nil, err
		}
	}

	return &user, nil
}

func (u UserRepo) ChargeBalance(userId, amount int) error {
	if _, err := u.db.Exec("UPDATE users SET balance = balance + $1 WHERE id=$2", amount, userId); err != nil {
		return err
	}

	return nil
}
