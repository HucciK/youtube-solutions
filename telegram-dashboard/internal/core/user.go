package core

import "fmt"

type User struct {
	ID      int     `json:"id"`
	Name    string  `json:"name"`
	Balance float64 `json:"balance"`
	HasKey  bool    `json:"hasKey"`
	States  InternalStates
}

type InternalStates struct {
	ChatState int
	MessageId int
}

func (u *User) Info() string {
	owner := "Нет"

	if u.HasKey {
		owner = "Да"
	}

	return fmt.Sprintf("Ваш ID: %d\nВаш Никнейм: %s\nВаш Баланс: %.2f\nВладеете ключом: %s", u.ID, u.Name, u.Balance, owner)
}

func (u *User) PaymentUrl(amount int) string {
	return fmt.Sprintf("https://api.crystalpay.ru/v1/?s=17592c19723a2f86dc96e22b6d28f18897c0ea35&n=youtubemarket&o=invoice-create&amount=%d&lifetime=60&callback=http://youtubesolutions.ru/api/chargeBalance&extra=%d", amount, u.ID)
}
