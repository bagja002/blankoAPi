package static

import "github.com/gofiber/fiber/v2"

func StaticPelatihan(c *fiber.Ctx) error {
	params := c.Params("string")
	return c.SendFile("public/static/pelatihan/"+ params)
}

func StaticPrasarana(c *fiber.Ctx)error{
	params := c.Params("string")
	return c.SendFile("public/static/prasarana/"+ params)
}

func StaticProfile(c *fiber.Ctx)error{
	params := c.Params("string")
	return c.SendFile("public/static/profile/"+ params)
}

func StaticSertifikasi(c *fiber.Ctx)error{
	params := c.Params("string")
	return c.SendFile("public/static/sertifikasi/"+ params)
}