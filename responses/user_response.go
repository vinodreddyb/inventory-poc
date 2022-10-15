package responses

import "github.com/gofiber/fiber/v2"

type UserResponse struct {
	Status  int        `json:"status"`
	Message string     `json:"message"`
	Body    *fiber.Map `json:"body"`
}
