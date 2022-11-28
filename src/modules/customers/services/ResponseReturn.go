package services

import (
	"encoding/json"
	"financial-api/src/utils/authorization"
	"financial-api/src/utils/database"
	"financial-api/src/utils/structs"
	"fmt"
)

type ResponseReturn struct {
	Customer structs.Customer `json:"customer"`
	Token    string           `json:"token"`
}

type AllTransactionsResponseReturn struct {
	Customer               structs.CustomerBasicInfos    `json:"customer"`
	AllAccountTransactions []structs.AccountTransactions `json:"allAccountTransactions"`
}

func NewResponseReturn(customer database.CustomerDB, account database.AccountDB, debitedTransactions []database.TransactionsDB, creditedTransactions []database.TransactionsDB) []byte {

	parsedCustomer := structs.Customer{
		ID:       customer.ID,
		Username: customer.Username,
	}

	debitedTransactionsChan := make(chan []structs.Transaction)
	creditedTransactionsChan := make(chan []structs.Transaction)

	go ParseTransaction(debitedTransactionsChan, debitedTransactions)
	go ParseTransaction(creditedTransactionsChan, creditedTransactions)

	parsedCustomer.Account = structs.Account{
		ID:                   account.ID,
		Balance:              int64(account.Balance),
		DebitedTransactions:  <-debitedTransactionsChan,
		CreditedTransactions: <-creditedTransactionsChan,
	}

	token, err := authorization.CreateToken(parsedCustomer.ID, parsedCustomer.Account.ID)

	if err != nil {
		fmt.Printf("Cannot generate token, reason: %s", err.Error())
	}

	responseReturn := ResponseReturn{
		Customer: parsedCustomer,
		Token:    token,
	}

	responseBytes, _ := json.Marshal(responseReturn)

	return responseBytes

}

func NewMainPageResponseReturn(customer database.CustomerDB, account database.AccountDB, allAccountTransactions []structs.AccountTransactions) []byte {
	parsedAccount := structs.BasicAccount{
		ID:      account.ID,
		Balance: int64(account.Balance),
	}

	parsedCustomer := structs.CustomerBasicInfos{
		ID:           customer.ID,
		Username:     customer.Username,
		BasicAccount: parsedAccount,
	}

	responseReturn := AllTransactionsResponseReturn{
		Customer:               parsedCustomer,
		AllAccountTransactions: []structs.AccountTransactions{},
	}

	for _, v := range allAccountTransactions {
		responseReturn.AllAccountTransactions = append(responseReturn.AllAccountTransactions, v)
	}

	responseBytes, _ := json.Marshal(responseReturn)

	return responseBytes

}

func ParseTransaction(transactionsChan chan []structs.Transaction, transactions []database.TransactionsDB) {
	allTransactions := []structs.Transaction{}

	for _, v := range transactions {
		allTransactions = append(allTransactions, structs.Transaction{
			ID:        v.ID,
			Value:     int64(v.Value),
			CreatedAt: v.CreatedAt,
		})
	}
	transactionsChan <- allTransactions
}
