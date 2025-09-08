package handlers

import (
	"miniMutualFund/models"

	"github.com/gofiber/fiber/v2"
)

type AuthHandler struct{}

func (h *AuthHandler) Login(c *fiber.Ctx) error {
	var req models.LoginRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid request"})
	}

	if req.UserID != "" || req.Password != "" {
		resp := models.LoginResponse{
			Status:  "success",
			Message: "login successful",
			Token:   "mock-jwt-token",
		}
		return c.JSON(resp)
	}

	return c.Status(401).JSON(models.LoginResponse{
		Status:  "error",
		Message: "invalid credentials",
	})
}
