package resources

import "fmt"

type InlineKeyboardMarkup struct {
	Keyboard [][]InlineKeyboardButton `json:"inline_keyboard"`
}

type InlineKeyboardButton struct {
	Text     string `json:"text"`
	Url      string `json:"url"`
	Callback string `json:"callback_data"`
}

var (
	StartMarkup = InlineKeyboardMarkup{
		Keyboard: [][]InlineKeyboardButton{
			[]InlineKeyboardButton{
				{
					Text:     "О подписке",
					Callback: "licenseInfo",
				},
				{
					Text:     "Профиль",
					Callback: "profile",
				},
			},
			[]InlineKeyboardButton{
				{
					Text:     "Пополнить баланс",
					Callback: "chargeBalance",
				},
			},
		},
	}

	AboutLicneseMarkup = InlineKeyboardMarkup{
		Keyboard: [][]InlineKeyboardButton{
			[]InlineKeyboardButton{
				{
					Text:     "Купить подписку",
					Callback: "buyLicense",
				},
				{
					Text:     "Назад",
					Callback: "return",
				},
			},
		},
	}

	OwnerMarkup = InlineKeyboardMarkup{
		Keyboard: [][]InlineKeyboardButton{
			[]InlineKeyboardButton{
				{
					Text:     "Отвязать IP",
					Callback: "unbindIp",
				},
				{
					Text:     "Назад",
					Callback: "return",
				},
			},
		},
	}

	RenewMarkup = InlineKeyboardMarkup{
		Keyboard: [][]InlineKeyboardButton{
			[]InlineKeyboardButton{
				{
					Text:     "Продлить подписку",
					Callback: "renewLicense",
				},
				{
					Text:     "Назад",
					Callback: "return",
				},
			},
		},
	}

	SelectPaymentTypeMarkup = InlineKeyboardMarkup{
		Keyboard: [][]InlineKeyboardButton{
			[]InlineKeyboardButton{
				{
					Text:     "Lolz",
					Callback: "fiatDeposit",
				},
			},
			[]InlineKeyboardButton{
				{
					Text:     "Назад",
					Callback: "return",
				},
			},
		},
	}

	ReturnMarkup = InlineKeyboardMarkup{
		Keyboard: [][]InlineKeyboardButton{
			[]InlineKeyboardButton{
				{
					Text:     "Назад",
					Callback: "return",
				},
			},
		},
	}
)

func FiatPaymentMarkup(url, invoiceId string) *InlineKeyboardMarkup {
	return &InlineKeyboardMarkup{
		[][]InlineKeyboardButton{
			[]InlineKeyboardButton{
				{
					Text: "Оплатить",
					Url:  url,
				},
			},
			[]InlineKeyboardButton{
				{
					Text:     "Проверить оплату",
					Callback: fmt.Sprintf("checkFiatPayment %s", invoiceId),
				},
			},
			[]InlineKeyboardButton{
				{
					Text:     "Назад",
					Callback: "return",
				},
			},
		},
	}
}
