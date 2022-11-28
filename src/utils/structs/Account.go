package structs

type Account struct {
	ID                   uint32        `json:"id"`
	Balance              int64         `json:"balance"`
	DebitedTransactions  []Transaction `json:"debitedTransactions"`
	CreditedTransactions []Transaction `json:"creditedTransactions"`
}

type BasicAccount struct {
	ID      uint32 `json:"id"`
	Balance int64  `json:"balance"`
}
