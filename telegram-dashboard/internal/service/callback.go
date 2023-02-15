package service

import (
	"fmt"
	"strings"
	"yt-solutions-telegram-dashboard/internal/core"
	"yt-solutions-telegram-dashboard/internal/resources"
)

type CallbackService struct {
	Client
	StateMachine
	DownloadLink string
	GuideLink    string
}

func NewCallbackService(cli Client, sm StateMachine, download, guide string) *CallbackService {
	return &CallbackService{
		Client:       cli,
		StateMachine: sm,
		DownloadLink: download,
		GuideLink:    guide,
	}
}

func (c CallbackService) ProcessCallback(userId, msgId int, username, callback, callbackId string) {

	user, ok := c.StateMachine.Get(userId)
	fmt.Println(user)
	if !ok {
		newUser, err := c.Client.AuthUser(userId, username)
		if err != nil {
			//
		}
		user = newUser
		user.States.MessageId = msgId
		c.StateMachine.Set(user)
	}

	switch callback {

	case resources.ProfileCmd:
		c.Profile(user)
	case resources.LicenseInfoCmd:
		c.ShowLicenseInfo(user)

	case resources.ChargeBalanceCmd:
		c.ChargeBalance(user)
	case resources.FiatChargingCmd:
		c.FiatBalanceCharging(user)
	case resources.CryptoChargingCmd:
		c.CryptoBalanceCharging(user)

	case resources.BuyLicenseCmd:
		c.ProcessOrder(user)
	case resources.RenewLicenseCmd:
		c.ProcessRenewal(user)

	case resources.UnbindIpCmd:
		c.UnbindAddress(user)

	case resources.ReturnCmd:
		c.MainMenu(user)

	default:
		cmd := strings.Split(callback, " ")[0]
		if cmd == resources.CheckFiatPaymentCmd {
			invoiceId := strings.Split(callback, " ")[1]
			c.CheckFiatPayment(user, callbackId, invoiceId)
		}
	}

}

func (c CallbackService) Profile(user *core.User) {
	if err := c.Client.EditMessage(user.Info(), user.ID, user.States.MessageId, resources.ReturnMarkup); err != nil {
		//TODO
		//
	}
}

func (c CallbackService) ShowLicenseInfo(user *core.User) {
	if user.HasKey {
		key, err := c.Client.GetKeyInfo(user.ID)
		if err != nil {
			//TODO
			//
		}

		if key.IsExpired {
			if err = c.Client.EditMessage(key.Info(), user.ID, user.States.MessageId, resources.RenewMarkup); err != nil {
				//TODO
				//
			}
			return
		}

		if err = c.Client.EditMessage(fmt.Sprintf("%s\n\nСсылка на скачивание:\n%s\n\nГайд:\n%s", key.Info(), c.DownloadLink, c.GuideLink), user.ID, user.States.MessageId, resources.OwnerMarkup); err != nil {
			//TODO
			//
		}
		return
	}

	free, lifetime, price, err := c.Client.GetAvailability()
	if err != nil {
		//TODO
		//
	}

	//TODO
	if free <= 0 && lifetime <= 0 {
		if err = c.Client.EditMessage(fmt.Sprintf("%s\n\nСейчас доступно:\n\nМесячная подписка на YouTube Solutions\n\nЦена: %d", resources.AboutBotTemplate, price), user.ID, user.States.MessageId, resources.AboutLicneseMarkup); err != nil {
			//TODO
			//
		}
		return
	}

	//Если фри и лайфтайм 0 - добавлять уведомление о том что будет продана renew версия
	if err = c.Client.EditMessage(fmt.Sprintf("%s\n\nСейчас доступно:\n\nFree копий: %d\n\nLifetime копий: %d\n\nЦена: %d", resources.AboutBotTemplate, free, lifetime, price), user.ID, user.States.MessageId, resources.AboutLicneseMarkup); err != nil {
		//TODO
		//
	}
}

func (c CallbackService) ProcessOrder(user *core.User) {

	if user.HasKey {
		//TODO
		// уведомление о наличие ключа
	}

	_, _, price, err := c.Client.GetAvailability()
	if err != nil {
		//TODO
		//
	}

	if user.Balance < float64(price) {
		if err = c.Client.EditMessage("Недостаточный баланс", user.ID, user.States.MessageId, resources.ReturnMarkup); err != nil {
			//TODO
			//
		}
		return
	}

	success, key, err := c.Client.ProcessOrder(user.ID)
	if err != nil {
		//TODO
		//
	}

	if !success {
		// уведомление о том что что-то пошло не так
		return
	}

	if err = c.Client.EditMessage(fmt.Sprintf("%s\n\nСсылка на скачивание:\n%s", key.Info(), c.DownloadLink), user.ID, user.States.MessageId, resources.ReturnMarkup); err != nil {
		//TODO
		//
	}

	user.Balance = user.Balance - float64(price)
	user.HasKey = true
}

func (c CallbackService) ProcessRenewal(user *core.User) {
	if !user.HasKey {
		return
	}

	key, err := c.Client.GetKeyInfo(user.ID)
	if err != nil {
		//TODO
		//
	}

	if !key.IsExpired {
		if err = c.Client.EditMessage("Срок действия ключа ещё не истёк", user.ID, user.States.MessageId, resources.ReturnMarkup); err != nil {
			//TODO
			//
		}
		return
	}

	renewal, err := c.Client.GetRenewalInfo()
	if err != nil {
		//TODO
		//
	}

	if user.Balance < float64(renewal) {
		if err = c.Client.EditMessage("Недостаточный баланс", user.ID, user.States.MessageId, resources.ReturnMarkup); err != nil {
			//TODO
			//
		}
		return
	}

	success, err := c.Client.ProcessRenewal(user.ID)
	if err != nil {
		//TODO
		//
	}

	if !success {
		// уведомление о том что что-то пошло не так
		return
	}

	if err = c.Client.EditMessage("Подписка успешно продлена", user.ID, user.States.MessageId, resources.ReturnMarkup); err != nil {
		//TODO
		//
	}
	user.Balance = user.Balance - float64(renewal)
}
func (c CallbackService) ChargeBalance(user *core.User) {

	_, _, price, err := c.Client.GetAvailability()
	if err != nil {
		//TODO
	}

	invoiceId, paymentUrl, err := c.Client.CreateInvoice(user.PaymentUrl(price))
	if err != nil {
		//TODO
		//
	}

	if err = c.Client.EditMessage(fmt.Sprintf("ID платежа %s", invoiceId), user.ID, user.States.MessageId, resources.FiatPaymentMarkup(paymentUrl, invoiceId)); err != nil {
		//TODO
		//
	}

}

// FiatBalanceCharging Currently not used
func (c CallbackService) FiatBalanceCharging(user *core.User) {

	_, _, price, err := c.Client.GetAvailability()
	if err != nil {
		//TODO
	}

	invoiceId, paymentUrl, err := c.Client.CreateInvoice(user.PaymentUrl(price))
	if err != nil {
		//TODO
		//
	}

	if err = c.Client.EditMessage(fmt.Sprintf("ID платежа %s", invoiceId), user.ID, user.States.MessageId, resources.FiatPaymentMarkup(paymentUrl, invoiceId)); err != nil {
		//TODO
		//
	}
}

// CryptoBalanceCharging Currently not used
func (c CallbackService) CryptoBalanceCharging(user *core.User) {
	if err := c.Client.EditMessage(resources.CryptoPaymentInfoTemplate, user.ID, user.States.MessageId, resources.ReturnMarkup); err != nil {
		//TODO
	}
	user.States.ChatState = resources.ProcessingCryptoPayment
}

func (c CallbackService) CheckFiatPayment(user *core.User, callbackId, invoiceId string) {
	amount, state, err := c.Client.CheckFiatPayment(invoiceId)
	if err != nil {
		//TODO
		//
	}

	if state == "payed" {

		if err := c.Client.EditMessage("Оплата успешно произведена", user.ID, user.States.MessageId, resources.ReturnMarkup); err != nil {
			return
		}

		user.Balance = user.Balance + float64(amount)
		return
	}

	if err := c.Client.AnswerCallback("Ожидает оплаты", callbackId); err != nil {
		//
	}
}

func (c CallbackService) UnbindAddress(user *core.User) {
	if err := c.Client.UnbindAddress(user.ID); err != nil {
		//TODO
		// Уведомление о то что что-то пошло не так
		return
	}

	if err := c.Client.EditMessage("IP успешно отвязан", user.ID, user.States.MessageId, resources.ReturnMarkup); err != nil {
		//
	}
}

func (c CallbackService) MainMenu(user *core.User) {
	if user.HasKey {
		if err := c.Client.EditMessage(resources.GreetWithKeyTemplate, user.ID, user.States.MessageId, resources.StartMarkup); err != nil {
			//TODO
			//
		}
		user.States.ChatState = resources.DefaultChatState
		return
	}

	if err := c.Client.EditMessage(resources.GreetTemplate, user.ID, user.States.MessageId, resources.StartMarkup); err != nil {
		//TODO
		//
	}
	user.States.ChatState = resources.DefaultChatState
}
