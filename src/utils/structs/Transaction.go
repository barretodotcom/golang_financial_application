package structs

import (
	"time"
)

type Transaction struct {
	ID                      uint32    `json:"id"`
	Value                   int64     `json:"value"`
	DebitedAccountUsername  string    `json:"debitedAccountUsername"`
	CreditedAccountUsername string    `json:"creditedAccountUsername"`
	CreatedAt               time.Time `json:"createdAt"`
}

type CreateTransaction struct {
	CreditedAccountUsername string `json:"creditedAccountUsername"`
	Value                   uint64 `json:"value"`
}

type TransactionBrokerMessage struct {
	DebitedAccountId  uint32 `json:"debitedAccountId"`
	CreditedAccountId uint32 `json:"creditedAccountId"`
	Value             uint64 `json:"value"`
}
