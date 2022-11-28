package structs

import "time"

type Customer struct {
	ID       uint32  `json:"id"`
	Username string  `json:"username"`
	Account  Account `json:"account"`
}

type CustomerBasicInfos struct {
	ID           uint32       `json:"id"`
	Username     string       `json:"username"`
	BasicAccount BasicAccount `json:"account"`
}

type AccountTransactions struct {
	Id                      uint32    `json:"id"`
	DebitedAccountUsername  string    `json:"debitedAccountUsername"`
	CreditedAccountUsername string    `json:"creditedAccountUsername"`
	DebitedAccountId        uint32    `json:"debitedAccountId"`
	CreditedAccountId       uint32    `json:"creditedAccountId"`
	Balance                 uint64    `json:"balance"`
	CreatedAt               time.Time `json:"createdAt"`
}
