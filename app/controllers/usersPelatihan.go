package controllers

import (
	"template/app/entity"
	"template/pkg/database"
	"template/pkg/tools"

	"github.com/gofiber/fiber/v2"
)

//From Register Akun User
func CreateUserPelatihan(c *fiber.Ctx )error{

	id_admin, _ := c.Locals("id_admin").(int)
	role, _ := c.Locals("role").(string)
	names, _ := c.Locals("name").(string)

	tools.ValidationJwtUsers(c, role, id_admin, names)


	var data map[string]string

	if err := c.BodyParser(&data); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"Message": "Failed to parse request body", "Error": err.Error()})
	}

	newUserPelatihan := entity.UsersPelatihan{
		IdUsers: uint(id_admin),
		IdPelatihan: uint(tools.StringToInt(data["id_pelatihan"])),
		StatusPembayaran: "pending",
		CreteAt: tools.TimeNowJakarta(),
	}


	database.DB.Create(&newUserPelatihan)


	var testing entity.UsersPelatihan

	database.DB.Preload("Users").Find(&testing)

	return c.JSON(fiber.Map{
		"pesan":"berhasil membuat data",
		"data":testing,
	})
}