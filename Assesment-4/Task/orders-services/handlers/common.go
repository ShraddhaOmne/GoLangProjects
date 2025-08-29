package handlers

import (
	"github.com/gofiber/fiber/v2"
	_ "github.com/gofiber/fiber/v2"
)

func Root(c *fiber.Ctx) error {
	return c.SendString("Hello OmneNEXT")
}
func Ping(c *fiber.Ctx) error {
	return c.SendString("pong")
}

func Health(c *fiber.Ctx) error {
	return c.SendString("ok")
}

/*
 func(c *fiber.Ctx) error {
		return c.SendString("Hello, Fiber!")
	})

*/
