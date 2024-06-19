package static

import "github.com/gofiber/fiber/v2"

// Area Static Pelatihan
func StaticPelatihan(c *fiber.Ctx) error {
	params := c.Params("string")
	return c.SendFile("public/static/pelatihan/" + params)
}

func StaticSilabusPelatihan(c *fiber.Ctx) error {
	params := c.Params("string")
	return c.SendFile("public/silabus/pelatihan/" + params)
}

func StaticSilabusSertifikasi(c *fiber.Ctx) error {
	params := c.Params("string")
	return c.SendFile("public/silabus/sertifikasi/" + params)
}

// Done
func StaticModulePelatihan(c *fiber.Ctx) error {
	params := c.Params("string")
	return c.SendFile("public/module/pelatihan/" + params)
}

// Done
func StaticBeritaAcara(c *fiber.Ctx) error {
	params := c.Params("string")
	return c.SendFile("public/static/BeritaAcara/" + params)
}
func StaticMemo(c *fiber.Ctx) error {
	params := c.Params("string")
	return c.SendFile("public/static/memo/" + params)
}

func StaticSuratPemberitahuan(c *fiber.Ctx) error {
	params := c.Params("string")
	return c.SendFile("public/static/suratPemberitahuan/" + params)
}

//----------------_End Static Pelatihan

//Usesr Static Area

//Foto, Ktp, KK, Ijzah, surat kesehatan

// Ini Static Foto Users
func StaticFotoUsers(c *fiber.Ctx) error {
	params := c.Params("string")
	return c.SendFile("public/static/profile/fotoProfile/" + params)
}

func StaticIjazah(c *fiber.Ctx) error {
	params := c.Params("string")
	return c.SendFile("public/static/profile/ijazah/" + params)
}
func StaticKK(c *fiber.Ctx) error {
	params := c.Params("string")
	return c.SendFile("public/static/profile/kk/" + params)
}

func StaticKtp(c *fiber.Ctx) error {
	params := c.Params("string")
	return c.SendFile("public/static/profile/ktp/" + params)
}
func StaticSuratSehat(c *fiber.Ctx) error {
	params := c.Params("string")
	return c.SendFile("public/static/profile/suratSehat/" + params)
}

// -------------end users area
func StaticPrasarana(c *fiber.Ctx) error {
	params := c.Params("string")
	return c.SendFile("public/static/prasarana/" + params)
}

func StaticProfile(c *fiber.Ctx) error {
	params := c.Params("string")
	return c.SendFile("public/static/profile/" + params)
}

func StaticSertifikasi(c *fiber.Ctx) error {
	params := c.Params("string")
	return c.SendFile("public/static/sertifikasi/" + params)
}

func StaticSarpras(c *fiber.Ctx) error {
	params := c.Params("string")
	return c.SendFile("public/static/sarpras/" + params)
}
