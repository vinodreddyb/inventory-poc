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
		return c.Status(http.StatusBadRequest).JSON(responses.APIResponse{Status: http.StatusBadRequest, Message: "error", Body: &fiber.Map{"data": err.Error()}})
	}

	usr, err := services.AddNewUser(&user)

	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(responses.APIResponse{Status: http.StatusInternalServerError, Message: "error", Body: &fiber.Map{"data": err.Error()}})
	}

	return c.Status(http.StatusCreated).JSON(responses.APIResponse{Status: http.StatusCreated, Message: "success", Body: &fiber.Map{"user": usr}})
}

func GetAllUsers(c *fiber.Ctx) error {
	users := services.GetAllUsers()
	return c.Status(http.StatusOK).JSON(
		responses.APIResponse{Status: http.StatusOK, Message: "success", Body: users},
	)
}

func GetAllCivil(c *fiber.Ctx) error {

	path := c.Query("path", "")
	civils, err := services.GetCivils(path)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(responses.APIResponse{Status: http.StatusInternalServerError, Message: "error", Body: &fiber.Map{"data": err.Error()}})
	}
	return c.Status(http.StatusOK).JSON(
		responses.APIResponse{Status: http.StatusOK, Message: "success", Body: civils},
	)
}

func GetAllCivilFields(c *fiber.Ctx) error {

	logr.Info("This is civil fields")
	civils, err := services.GetCivilFields()
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(responses.APIResponse{Status: http.StatusInternalServerError, Message: "error", Body: &fiber.Map{"data": err.Error()}})
	}
	return c.Status(http.StatusOK).JSON(
		responses.APIResponse{Status: http.StatusOK, Message: "success", Body: civils},
	)
}

func AddNewNode(c *fiber.Ctx) error {

	logr.Info("Add new node")
	var civilNode models.CivilDTO
	nodeId := c.Params("nodeId", "")
	if err := c.BodyParser(&civilNode); err != nil {
		return c.Status(http.StatusBadRequest).JSON(responses.APIResponse{Status: http.StatusBadRequest, Message: "error", Body: &fiber.Map{"data": err.Error()}})
	}

	nodes, err := services.AddCivilNode(nodeId, civilNode)

	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(responses.APIResponse{Status: http.StatusInternalServerError, Message: "error", Body: &fiber.Map{"data": err.Error()}})
	}

	return c.Status(http.StatusCreated).JSON(responses.APIResponse{Status: http.StatusCreated, Message: "success", Body: nodes})
}

func UpdateNodeValues(c *fiber.Ctx) error {

	logr.Info("Update node values")
	var civilNode models.CivilDTO

	if err := c.BodyParser(&civilNode); err != nil {
		return c.Status(http.StatusBadRequest).JSON(responses.APIResponse{Status: http.StatusBadRequest, Message: "error", Body: &fiber.Map{"data": err.Error()}})
	}

	result, err := services.UpdateCivilNode(civilNode)

	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(responses.APIResponse{Status: http.StatusInternalServerError, Message: "error", Body: &fiber.Map{"data": err.Error()}})
	}

	return c.Status(http.StatusCreated).JSON(responses.APIResponse{Status: http.StatusCreated, Message: "success", Body: result})
}

func AddStatus(c *fiber.Ctx) error {

	logr.Info("Add status")
	var civilNode models.CivilProgressDTO

	if err := c.BodyParser(&civilNode); err != nil {
		return c.Status(http.StatusBadRequest).JSON(responses.APIResponse{Status: http.StatusBadRequest, Message: "error", Body: &fiber.Map{"data": err.Error()}})
	}

	result, err := services.AddWorkStatus(civilNode)

	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(responses.APIResponse{Status: http.StatusInternalServerError, Message: "error", Body: &fiber.Map{"data": err.Error()}})
	}

	return c.Status(http.StatusCreated).JSON(responses.APIResponse{Status: http.StatusCreated, Message: "success", Body: result})
}

func GetCivilProgress(c *fiber.Ctx) error {

	nodeId := c.Params("nodeId", "")
	graph := c.Query("graph", "false")
	civilProgress, graphData, err := services.GetStatusGraph(nodeId)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(responses.APIResponse{Status: http.StatusInternalServerError, Message: "error", Body: &fiber.Map{"data": err.Error()}})
	}
	if graph == "true" {
		return c.Status(http.StatusOK).JSON(
			responses.APIResponse{Status: http.StatusOK, Message: "success", Body: graphData},
		)
	}
	return c.Status(http.StatusOK).JSON(
		responses.APIResponse{Status: http.StatusOK, Message: "success", Body: civilProgress},
	)

}
