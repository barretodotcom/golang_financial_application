package repositories

import (
	"financial-api/src/utils/database"

	"gorm.io/gorm"
)

type Repository struct {
	DB *gorm.DB
}

func (ar *Repository) FindByAccountId(accountId uint32) database.AccountDB {

	var account database.AccountDB
	ar.DB.First(&database.AccountDB{ID: accountId}).Take(&account)

	return account
}

func (ar *Repository) FindByCustomerId(customerId uint32) database.AccountDB {
	var account database.AccountDB
	ar.DB.Raw("SELECT * FROM account_dbs WHERE customer_id = ?", customerId).Take(&account)

	return account
}

func (ar *Repository) CreateAccount(customerId uint32) database.AccountDB {
	account := database.AccountDB{
		Balance:    10000,
		CustomerId: customerId,
	}

	ar.DB.Create(&account)

	return account
}
