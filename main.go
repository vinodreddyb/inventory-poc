package main

import (
	"github.com/gofiber/fiber/v2"
	"mongo-rest/configs"
	"mongo-rest/routes"
)

//https://dev.to/hackmamba/build-a-rest-api-with-golang-and-mongodb-fiber-version-4la0
func main() {

	app := fiber.New()

	configs.ConnectDB()

	routes.UserRoute(app)

	app.Listen(":6000")
}
