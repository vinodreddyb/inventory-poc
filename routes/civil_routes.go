package routes

import (
	"github.com/gofiber/fiber/v2"
	"inventory-poc/controllers"
)

func CivilRoutes(app *fiber.App) {
	app.Post("/user", controllers.CreateUser)
	app.Get("/users", controllers.GetAllUsers)
	app.Get("/civil", controllers.GetAllCivil)
	app.Get("/civil/fields", controllers.GetAllCivilFields)
	app.Post("/civil/node/:nodeId", controllers.AddNewNode)
	app.Put("/civil/node/values", controllers.UpdateNodeValues)
	app.Post("/civil/status", controllers.AddStatus)
	app.Get("/civil/status/:nodeId", controllers.GetCivilProgress)
	app.Get("/contract/scurve", controllers.GetContractProgress)
	app.Get("/contract/schedule", controllers.GetContractScheduleProgress)
	app.Get("/contract/piechart", controllers.GetContractScheduleProgressPieChart)
}
