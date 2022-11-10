package routes

import (
	"github.com/gofiber/fiber/v2"
	"mongo-rest/controllers"
)

func CivilRoutes(app *fiber.App) {
	app.Post("/user", controllers.CreateUser)
	app.Get("/users", controllers.GetAllUsers)
	app.Get("/civil", controllers.GetAllCivil)
	app.Get("/civil/fields", controllers.GetAllCivilFields)
	app.Post("/civil/node", controllers.AddNewNode)
	app.Put("/civil/node/values", controllers.UpdateNodeValues)
}
