package controllers

import (
	"template/app/entity"
	"template/pkg/database"
	"template/pkg/tools"

	"github.com/gofiber/fiber/v2"
)

// From Register Akun User
func CreateUserPelatihan(c *fiber.Ctx) error {

	id_admin, _ := c.Locals("id_admin").(int)
	role, _ := c.Locals("role").(string)
	names, _ := c.Locals("name").(string)

	tools.ValidationJwtUsers(c, role, id_admin, names)

	var data map[string]string

	if err := c.BodyParser(&data); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"Message": "Failed to parse request body", "Error": err.Error()})
	}

	newUserPelatihan := entity.UsersPelatihan{
		IdUsers:          uint(id_admin),
		Nama:             names,
		IdPelatihan:      uint(tools.StringToInt(data["id_pelatihan"])),
		StatusPembayaran: "pending",
		CreteAt:          tools.TimeNowJakarta(),
	}

	database.DB.Create(&newUserPelatihan)

	var testing entity.UsersPelatihan

	database.DB.Preload("Users").Find(&testing)

	return c.JSON(fiber.Map{
		"pesan": "berhasil membuat data",
		"data":  testing,
	})
}

func GetUserPelatihan(c *fiber.Ctx) error {

	/*
		id_admin, _ := c.Locals("id_admin").(int)
		role, _ := c.Locals("role").(string)
		names, _ := c.Locals("name").(string)

		tools.ValidationJwtUsers(c, role, id_admin, names)

	*/

	id_users := c.Query("idUsers")
	id_pelatihan := c.Query("idPelatihan")

	var usersPelatihan []entity.UsersPelatihan
	baseQuey := database.DB

	if id_users != "" {
		baseQuey = baseQuey.Where("id_users = ?", id_users)
	}
	if id_pelatihan != "" {
		baseQuey = baseQuey.Where("id_pelatihan = ?", id_pelatihan)
	}

	baseQuey.Find(&usersPelatihan)

	return c.JSON(fiber.Map{
		"pesan": "Sukses Mengambil data",
		"data":  usersPelatihan,
	})
}

// Test dengan Gorm relasi cuy
func GetPelatihanByUser(c *fiber.Ctx) error {
	userID := c.Query("userID")
	var user entity.Users

	if err := database.DB.Preload("Pelatihan").First(&user, userID).Error; err != nil {
		return c.Status(404).SendString(err.Error())
	}

	return c.JSON(user)
}

func GetUsersByPelatihan(c *fiber.Ctx) error {
	/*
		//JWT nya siapa ntar ?
		id_admin, _ := c.Locals("id_admin").(int)
		role, _ := c.Locals("role").(string)
		names, _ := c.Locals("name").(string)

		tools.ValidationJwt(c, role, id_admin, names)

	*/

	idPelatihan := c.Query("idPelatihan")

	var pelatihan entity.Pelatihan

	if err := database.DB.Preload("UserPelatihan").Find(&pelatihan, idPelatihan).Error; err != nil {
		return c.Status(404).SendString(err.Error())
	}

	return c.JSON(pelatihan)
}
