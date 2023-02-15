package core

import "time"

type Transaction struct {
	Id              string
	Payeer          int
	CheckExpiration int64
}

type TransactionInfo struct {
	Confirmations     int               `json:"confirmations"`
	Confirmed         bool              `json:"confirmed"`
	Timestamp         int64             `json:"timestamp"`
	TokenTransferInfo TokenTransferInfo `json:"tokenTransferInfo"`
}

type TokenTransferInfo struct {
	Name            string `json:"name"`
	Type            string `json:"type"`
	AmountStr       string `json:"amount_str"`
	ContractAddress string `json:"contract_address"`
	FromAddress     string `json:"from_address"`
	ToAddress       string `json:"to_address"`
}

func (t *Transaction) IsExpired() bool {
	if t.CheckExpiration < time.Now().Unix() {
		return true
	}
	return false
}
