package main

import (
	"fa-be/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func main() {
	app := fiber.New()
	app.Use(cors.New())
	app.Get("/api", func(c *fiber.Ctx) error {
		loacales := [7]string{
			"TW",
			"NG",
			"US",
			"GH",
			"ES",
			"NL",
			"DE",
		}
		data, _ := utils.ScrapeData(loacales[:])
		return c.SendString(string(data))
	})

	app.Listen(":9000")
}
