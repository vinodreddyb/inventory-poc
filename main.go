package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	logs "github.com/sirupsen/logrus"
	"io"
	"mongo-rest/configs"
	"mongo-rest/routes"
	"os"
)

//https://dev.to/hackmamba/build-a-rest-api-with-golang-and-mongodb-fiber-version-4la0
func main() {

	fileAccessLog, fileAppLog, err := logConfig()
	accessLogMW := io.MultiWriter(os.Stdout, fileAccessLog)
	appLogMW := io.MultiWriter(os.Stdout, fileAppLog)
	logs.SetOutput(appLogMW)

	configs.ConnectDB()

	app := fiber.New()
	app.Use(logger.New(logger.Config{Output: accessLogMW}), cors.New(cors.Config{
		AllowHeaders:     "Origin,Content-Type,Accept,Content-Length,Accept-Language,Accept-Encoding,Connection,Access-Control-Allow-Origin",
		AllowOrigins:     "*",
		AllowCredentials: true,
		AllowMethods:     "GET,POST,HEAD,PUT,DELETE,PATCH,OPTIONS",
	}))
	routes.CivilRoutes(app)

	err = app.Listen(":8080")
	if err != nil {
		logs.Fatal("Error while starting app ", err)
	}

	defer func() {
		ferr := fileAccessLog.Close()
		if ferr != nil {
			logs.Fatal("Error closing access logs", ferr)
			return
		}

		ferr1 := fileAppLog.Close()
		if ferr1 != nil {
			logs.Fatal("Error closing access logs", ferr1)
			return
		}
	}()
}

func logConfig() (*os.File, *os.File, error) {
	file, err := os.OpenFile("./inventory_access.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)

	if err != nil {
		logs.Fatal(err)
	}
	filel, errl := os.OpenFile("inventory.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if errl != nil {
		logs.Fatal(errl)
	}

	return file, filel, err
}
