package repositories

import (
	"financial-api/src/modules/customers/services"
	"financial-api/src/utils"
	"financial-api/src/utils/database"

	"github.com/gofiber/fiber/v2"
)

func (r *Repository) CreateCustomer(ctx *fiber.Ctx) error {
	var customer database.CustomerDB

	ctx.BodyParser(&customer)

	customerExists, _ := r.FindByUsername(customer.Username)

	if customerExists.Username == customer.Username {
		ctx.Status(400)
		ctx.Write(utils.NewResponseMessage("This username is already taken", "error"))
		return nil
	}

	hashedPassword, _ := utils.HashPassword(customer.Password)

	customer.Password = hashedPassword

	r.SaveCustomer(&customer)

	customerAccount := r.CreateAccount(customer.ID)

	r.AssociateAccount(customer.ID, customerAccount.ID)

	ctx.Status(201)
	ctx.Write(services.NewResponseReturn(customer, customerAccount, []database.TransactionsDB{}, []database.TransactionsDB{}))

	return nil
}

func (r *Repository) AuthCustomer(ctx *fiber.Ctx) error {

	var customer database.CustomerDB

	ctx.BodyParser(&customer)

	customerExists, err := r.FindByUsername(customer.Username)
	if err != nil {
		ctx.Status(400)
		ctx.Write(utils.NewResponseMessage("User not found.", "error"))
		return nil
	}

	if !utils.CheckPasswordHash(customer.Password, customerExists.Password) {
		ctx.Status(400)
		ctx.Write(utils.NewResponseMessage("Wrong password.", "error"))
		return nil
	}
	account := r.FindByCustomerId(customerExists.ID)

	debitedTransactions, creditedTransactions := r.FindAllAccountTransactions(account.ID)

	ctx.Status(200)
	ctx.Write(services.NewResponseReturn(*customerExists, account, debitedTransactions, creditedTransactions))

	return nil
}

func (cr *Repository) FindByUsername(username string) (*database.CustomerDB, error) {

	var customer database.CustomerDB
	err := cr.DB.Where(&database.CustomerDB{Username: username}).Take(&customer).Error

	return &customer, err
}

func (cr *Repository) FindCustomerById(customerId uint32) database.CustomerDB {

	var customer database.CustomerDB
	_ = cr.DB.Where(&database.CustomerDB{ID: customerId}).Take(&customer).Error

	return customer
}

func (cr *Repository) SaveCustomer(customer *database.CustomerDB) database.CustomerDB {
	var insertedCustomer database.CustomerDB
	_ = cr.DB.Create(customer).Take(&insertedCustomer).Error

	return insertedCustomer
}

func (cr *Repository) AssociateAccount(customerId uint32, accountId uint32) {
	cr.DB.Table("customer_dbs").Where(database.CustomerDB{ID: customerId}).Update("account_id", accountId)
}
