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
	id                     uint32
	debitedAccountUsername string
	creditedAccountUsernam string
	debitedAccountId       uint32
	creditedAccountId      uint32
	balance                uint64
	createdAt              time.Time
}
