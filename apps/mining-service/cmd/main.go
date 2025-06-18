package main

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
)

func main() {
	fiberApp := fiber.New()

	fiberApp.Get("/", func(c *fiber.Ctx) error {
		fmt.Println("Mining service is running")
		return c.SendString("Mining service is running")
	})

	fiberApp.Listen(":4000")
}
