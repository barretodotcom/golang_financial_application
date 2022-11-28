package main

import (
	"log"

	"financial-api/src/modules/repositories"
	"financial-api/src/utils/database"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/joho/godotenv"
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}

}

func main() {

	db, err := database.Connect()

	if err != nil {
		log.Fatal(err.Error())
	}

	r := repositories.Repository{
		DB: db,
	}

	app := fiber.New()
	app.Use(cors.New())
	r.SetupRoutes(app)
	app.Listen(":3000")
}
