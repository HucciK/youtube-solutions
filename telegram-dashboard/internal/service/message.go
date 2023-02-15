package service

import (
	"fmt"
	"math"
	"strconv"
	"time"
	"yt-solutions-telegram-dashboard/internal/core"
	"yt-solutions-telegram-dashboard/internal/resources"
)

const (
	usdtTrc20Contract = "TR7NHqjeKQxGTCi8q8ZY4pL8otSzgjLj6t"
	ownerAddress      = "TLHt5hKgtbe7hgc7yCbrttpoG77SMeQWYy"
)

type MessageService struct {
	Client
	StateMachine
	TransactionManager
}

func NewMessageService(cli Client, sm StateMachine, tm TransactionManager) *MessageService {
	return &MessageService{
		Client:             cli,
		StateMachine:       sm,
		TransactionManager: tm,
	}
}

func (m MessageService) ProcessMessage(userId, messageId int, username, message string) {
	user, ok := m.StateMachine.Get(userId)
	fmt.Println(user)
	if !ok {
		fmt.Println("CREATING NEW")
		newUser, err := m.Client.AuthUser(userId, username)
		if err != nil {
			//
		}
		fmt.Println("AUTH USER", newUser)
		user = newUser
		m.StateMachine.Set(user)
	}

	switch message {
	case resources.StartCmd:
		m.Start(user)
	}

	switch user.States.ChatState {
	case resources.ProcessingCryptoPayment:
		m.CheckCryptoPayment(user, messageId, message)
	}
}

func (m MessageService) Start(user *core.User) {
	if user.HasKey {
		msg, err := m.Client.SendMessage(resources.GreetWithKeyTemplate, user.ID, resources.StartMarkup)
		if err != nil {
			//
		}
		user.States.MessageId = msg
		return
	}

	msg, err := m.Client.SendMessage(resources.GreetTemplate, user.ID, resources.StartMarkup)
	if err != nil {
		//
	}
	user.States.MessageId = msg
}

func (m MessageService) CheckCryptoPayment(user *core.User, messageId int, message string) {
	//Проверка транзакции( Время создания, кол-во подтверждений, сумма )
	txInfo, err := m.Client.CheckCryptoPayment(message)
	if err != nil {
		//TODO
	}

	if txInfo.TokenTransferInfo.ToAddress != ownerAddress {
		//Уведомление о неправильном получателе
		if err := m.Client.EditMessage(fmt.Sprintf("%s\n\nСтатус: Адрес получателя не соответствует указанному или переведены неверные активы", resources.CryptoPaymentInfoTemplate), user.ID, user.States.MessageId, resources.ReturnMarkup); err != nil {
			//TODO
		}

		if err := m.Client.DeleteMessage(user.ID, messageId); err != nil {
			//TODO
		}
		return
	}

	if !txInfo.Confirmed {
		//Уведомление что транзакция ещё не доставлена
		if err := m.Client.EditMessage(fmt.Sprintf("%s\n\nСтатус: Траназкция ожидает подтверждения в сети", resources.CryptoPaymentInfoTemplate), user.ID, user.States.MessageId, resources.ReturnMarkup); err != nil {
			//TODO
		}

		if err := m.Client.DeleteMessage(user.ID, messageId); err != nil {
			//TODO
		}
		return
	}

	if txInfo.TokenTransferInfo.Type != "Transfer" || txInfo.TokenTransferInfo.ContractAddress != usdtTrc20Contract {
		//Уведомление что в транзакции что то не так
		if err := m.Client.EditMessage(fmt.Sprintf("%s\n\nСтатус: Перевод не в USDT или иные несоответсвия в транзации", resources.CryptoPaymentInfoTemplate), user.ID, user.States.MessageId, resources.ReturnMarkup); err != nil {
			//TODO
		}

		if err := m.Client.DeleteMessage(user.ID, messageId); err != nil {
			//TODO
		}
		return
	}

	if txInfo.Timestamp+604800000 < time.Now().UnixMilli() {
		//Уведомление о слишком старой транзакции
		if err := m.Client.EditMessage(fmt.Sprintf("%s\n\nСтатус: Слишком старая транзакция", resources.CryptoPaymentInfoTemplate), user.ID, user.States.MessageId, resources.ReturnMarkup); err != nil {
			//TODO
		}

		if err := m.Client.DeleteMessage(user.ID, messageId); err != nil {
			//TODO
		}
		return
	}

	tx, ok := m.TransactionManager.Get(message)
	if ok {
		//Уведомление о том что транзакция уже была
		if err := m.Client.EditMessage(fmt.Sprintf("%s\n\nСтатус: По данной транзакции уже было совершено зачисление", resources.CryptoPaymentInfoTemplate), user.ID, user.States.MessageId, resources.ReturnMarkup); err != nil {
			//TODO
		}

		if err := m.Client.DeleteMessage(user.ID, messageId); err != nil {
			//TODO
		}
		return
	}

	tx = &core.Transaction{
		Id:              message,
		Payeer:          user.ID,
		CheckExpiration: txInfo.Timestamp + 604800,
	}
	m.TransactionManager.Set(tx)

	amount, err := strconv.ParseFloat(txInfo.TokenTransferInfo.AmountStr, 64)
	amount = math.Ceil((amount / 1000000) * 70)

	//Получить актуальный курс usdt И пополнить баланс

	if err := m.Client.ChargeBalance(user.ID, int(amount)); err != nil {
		if err := m.Client.EditMessage("Что-то пошло не так, обратитесь в поддержку @tetrapachof", user.ID, user.States.MessageId, resources.ReturnMarkup); err != nil {
			//TODO
		}
		return
	}
	user.Balance = user.Balance + amount

	if err := m.Client.EditMessage(fmt.Sprintf("Успешно пополнено на %d", int(amount)), user.ID, user.States.MessageId, resources.ReturnMarkup); err != nil {
		//TODO
	}

	if err := m.Client.DeleteMessage(user.ID, messageId); err != nil {
		//TODO
	}
}
