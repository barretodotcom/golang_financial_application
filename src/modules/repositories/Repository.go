package repositories

import (
	"financial-api/src/utils/authorization"

	"github.com/gofiber/fiber/v2"
)

func (r *Repository) SetupRoutes(app *fiber.App) {

	customersRoutes := app.Group("/customers")

	customersRoutes.Post("/create", r.CreateCustomer)
	customersRoutes.Post("/session", r.AuthCustomer)

	transactionsRoutes := app.Group("/transactions")
	transactionsRoutes.Get("/find-all-user-transactions", authorization.VerifyToken(), r.FindAllUserTransactions)
	transactionsRoutes.Post("/send-transaction-message", authorization.VerifyToken(), r.SendTransactionMessage)
}
