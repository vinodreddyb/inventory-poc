package controllers

import (
	"github.com/gofiber/fiber/v2"
	"mongo-rest/models"
	"mongo-rest/responses"
	"mongo-rest/services"
	"net/http"
)

func CreateUser(c *fiber.Ctx) error {

	var user models.User

	if err := c.BodyParser(&user); err != nil {
		return c.Status(http.StatusBadRequest).JSON(responses.UserResponse{Status: http.StatusBadRequest, Message: "error", Body: &fiber.Map{"data": err.Error()}})
	}

	usr, err := services.AddNewUser(&user)

	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(responses.UserResponse{Status: http.StatusInternalServerError, Message: "error", Body: &fiber.Map{"data": err.Error()}})
	}

	return c.Status(http.StatusCreated).JSON(responses.UserResponse{Status: http.StatusCreated, Message: "success", Body: &fiber.Map{"user": usr}})
}

func GetAllUsers(c *fiber.Ctx) error {
	users := services.GetAllUsers()
	return c.Status(http.StatusOK).JSON(
		responses.UserResponse{Status: http.StatusOK, Message: "success", Body: &fiber.Map{"users": users}},
	)
}

func GetAllCivil(c *fiber.Ctx) error {

	civils := services.GetCivils()

	return c.Status(http.StatusOK).JSON(
		responses.UserResponse{Status: http.StatusOK, Message: "success", Body: &fiber.Map{"civils": civils}},
	)
}
