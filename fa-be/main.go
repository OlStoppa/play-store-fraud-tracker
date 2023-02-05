package main

import (
	"encoding/json"
	"fa-be/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func main() {
	app := fiber.New()
	app.Use(cors.New())
	app.Post("/api/search", func(c *fiber.Ctx) error {
		payload := struct {
			Locales []string `json:"locales" form:"locales"`
		}{}

		err := c.BodyParser(&payload)

		if c.Query("searchTerm") == "" {
			return c.Status(400).SendString("No Search Term Supplied")
		}

		if err != nil {
			return c.Status(400).SendString("No Locales Supplied")
		}

		data, _ := utils.ScrapeData(payload.Locales, c.Query("searchTerm"), c.Query("keyword"))

		byteArr, _ := json.Marshal(data)
		return c.SendString(string(byteArr))
	})

	app.Listen(":9000")
}
