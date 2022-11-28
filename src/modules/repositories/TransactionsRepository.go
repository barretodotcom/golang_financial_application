package repositories

import (
	"encoding/json"
	"financial-api/src/modules/customers/services"
	"financial-api/src/utils"
	"financial-api/src/utils/authorization"
	"financial-api/src/utils/broker"
	"financial-api/src/utils/database"
	"financial-api/src/utils/structs"
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/streadway/amqp"
)

func (tr *Repository) FindAllAccountTransactions(accountId uint32) ([]database.TransactionsDB, []database.TransactionsDB) {
	var debitedTransactions []database.TransactionsDB
	var creditedTransactions []database.TransactionsDB
	tr.DB.Find(&database.TransactionsDB{DebitedAccountId: accountId}).Take(&debitedTransactions)
	tr.DB.Find(&database.TransactionsDB{CreditedAccountId: accountId}).Take(&creditedTransactions)

	return debitedTransactions, creditedTransactions
}

func (tr *Repository) FindAllUserTransactions(c *fiber.Ctx) error {

	token := c.Get("Authorization")

	accountId, customerId := authorization.GetTokenSecret(token)

	customer := tr.FindCustomerById(customerId)
	account := tr.FindByAccountId(accountId)

	allUserTransactions := FindAllUserTransactionsRepository(accountId, tr)

	responseReturn := services.NewMainPageResponseReturn(customer, account, allUserTransactions)

	c.Write(responseReturn)
	return nil
}

func (tr *Repository) SendTransactionMessage(c *fiber.Ctx) error {
	token := c.Get("Authorization")

	var transactionData structs.CreateTransaction

	c.BodyParser(&transactionData)
	debitedAccountId, _ := authorization.GetTokenSecret(token)
	customer, _ := tr.FindByUsername(transactionData.CreditedAccountUsername)

	if customer.ID == debitedAccountId {
		c.Status(400)
		c.Write(utils.NewResponseMessage("You cant send a transaction to yourself", "error"))
		return nil
	}

	transactionData.Value = transactionData.Value * 100

	if transactionData.Value < 0 {
		c.Status(400)
		c.Write(utils.NewResponseMessage("Value need to be greater than zero.", "error"))
		return nil
	}

	if customer.Username == "" {
		c.Status(400)
		c.Write(utils.NewResponseMessage("User not found", "error"))
		return nil
	}

	account := tr.FindByCustomerId(debitedAccountId)

	if account.Balance < uint64(transactionData.Value) {
		c.Status(400)
		c.Write(utils.NewResponseMessage("You don't have balance.", "error"))
		return nil
	}

	channel := broker.GetChannel()

	message := amqp.Publishing{
		ContentType: "text/plain",
		Body:        NewTransactionBrokerMessage(debitedAccountId, customer.ID, transactionData.Value),
	}

	err := channel.Publish("", "transactions", false, false, message)

	if err != nil {
		fmt.Println(err.Error())
	}
	c.Status(201)
	c.Write(utils.NewResponseMessage("Transaction sent with success", "sucess"))
	return nil
}

func NewTransactionBrokerMessage(debitedAccountId, creditedAccountId uint32, value uint64) []byte {
	message := structs.TransactionBrokerMessage{
		DebitedAccountId:  debitedAccountId,
		CreditedAccountId: creditedAccountId,
		Value:             value,
	}

	stringfied, _ := json.Marshal(message)

	return stringfied
}

func FindAllUserTransactionsRepository(accountId uint32, tr *Repository) []structs.AccountTransactions {
	var allTransactions []structs.AccountTransactions
	tr.DB.Raw(allAccountTransactionsQuery, accountId, accountId).Take(&allTransactions)
	fmt.Println(allTransactions[0])
	return allTransactions
}

var allAccountTransactionsQuery string = `
	SELECT
		t.id,
		c.username as "DebitedAccountUsername",
		c2.username  as "CreditedAccountUsername",
		t.value as "Balance",
		t."created_at" as "CreatedAt"
	FROM
		account_dbs a
		INNER JOIN transactions_dbs t 
		ON
		a.id = t."debited_account_id"
	INNER JOIN customer_dbs c 
	ON
		c.id = t."debited_account_id" 
	INNER JOIN customer_dbs c2 
	ON
		c2.id = t."credited_account_id"
	WHERE 
		t."credited_account_id" = ? or t."debited_account_id" = ?
	`
