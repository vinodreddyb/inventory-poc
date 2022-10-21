package controllers

import (
	"github.com/gofiber/fiber/v2"
	logr "github.com/sirupsen/logrus"
	"mongo-rest/models"
	"mongo-rest/responses"
	"mongo-rest/services"
	"net/http"
)

func CreateUser(c *fiber.Ctx) error {

	logr.Info("This is create user")
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
		responses.UserResponse{Status: http.StatusOK, Message: "success", Body: users},
	)
}

func GetAllCivil(c *fiber.Ctx) error {

	civils, err := services.GetCivils()
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(responses.UserResponse{Status: http.StatusInternalServerError, Message: "error", Body: &fiber.Map{"data": err.Error()}})
	}
	return c.Status(http.StatusOK).JSON(
		responses.UserResponse{Status: http.StatusOK, Message: "success", Body: civils},
	)
}

func GetAllCivilFields(c *fiber.Ctx) error {

	logr.Info("This is civil fields")
	civils, err := services.GetCivilFields()
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(responses.UserResponse{Status: http.StatusInternalServerError, Message: "error", Body: &fiber.Map{"data": err.Error()}})
	}
	return c.Status(http.StatusOK).JSON(
		responses.UserResponse{Status: http.StatusOK, Message: "success", Body: civils},
	)
}
