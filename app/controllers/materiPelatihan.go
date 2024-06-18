package controllers

import (
	"template/app/entity"
	"template/pkg/database"
	"template/pkg/tools"

	"github.com/gofiber/fiber/v2"
)

func CreateMateriPelatihan(c *fiber.Ctx) error {

	id_admin, _ := c.Locals("id_admin").(int)
	role, _ := c.Locals("role").(string)
	names, _ := c.Locals("name").(string)

	tools.ValidationJwtLemdik(c, role, id_admin, names)

	var data map[string]string

	if err := c.BodyParser(&data); err != nil {
		c.JSON(fiber.Map{
			"pesan": err.Error(),
		})
	}

	idPelatihan := c.Query("id_pelatihan")

	var pelatihan entity.Pelatihan

	database.DB.Where("id_pelatihan = ? ", idPelatihan).Find(&pelatihan)
	if pelatihan.IdPelatihan == 0 {
		return c.Status(400).JSON(fiber.Map{
			"Pesan": "tidak ada pelatihan",
		})
	}

	newMateriPelatihan := entity.MateriPelatihan{
		IdPelatihan: uint(tools.StringToInt(idPelatihan)),
		NamaMateri:  data["nama_materi"],
		Deskripsi:   data["deskripsi"],
		JamTeory:    data["jam_teory"],
		JamPraktek:  data["jam_praktek"],
	}

	database.DB.Create(&newMateriPelatihan)

	return c.JSON(fiber.Map{
		"Pesan": "Sukses Membuat Materi Pelatihan",
		"data":  newMateriPelatihan,
	})
}
