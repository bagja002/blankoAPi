package static

import "github.com/gofiber/fiber/v2"

func StaticBelankoRusak(c *fiber.Ctx) error {
	params := c.Params("string")
	return c.SendFile("public/static/foto-blanko-rusak/" + params)
}
