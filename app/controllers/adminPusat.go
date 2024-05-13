package controllers

import (
	"template/app/entity"
	//"template/app/models"
	"template/pkg/database"
	"template/pkg/tools"
	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
)

func SuperAdminLogin(c *fiber.Ctx) error {
	var data map[string]string

	if err := c.BodyParser(&data); err != nil {
		c.JSON(fiber.Map{
			"pesan": err.Error(),
		})
	}

	var users entity.SuperAdmin

	database.DB.Where("username = ? ", data["username"]).First(&users)
	if users.IdSuperAdmin == 0 {
		c.Status(fiber.StatusUnauthorized)
		return c.JSON(fiber.Map{
			"pesan": "Username tidak di temukan",
		})
	}

	if err := bcrypt.CompareHashAndPassword([]byte(users.Password), []byte(data["password"])); err != nil {
		c.Status(fiber.StatusUnauthorized)
		return c.JSON(fiber.Map{
			"pesan": "Incorrect password!",
		})
	} else {
		t := tools.GenerateToken(users)
		return c.JSON(fiber.Map{
			"t": t,
		})
	}

}

func CreateAdminPusat(c *fiber.Ctx) error {

	id_admin, _ := c.Locals("id_admin").(string)
	role, _ := c.Locals("role").(string)
	names, _ := c.Locals("name").(string)

	if role != "1" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"pesan": "Role Bukan Super Admin ",
		})
	}
	if id_admin == "" {
		return c.JSON(fiber.Map{
			"pesan": "Admin tidak terdaftar ",
		})
	}
	if names == "" {
		return c.JSON(fiber.Map{
			"pesan": "Tidak ada Nama di dalam Jwt",
		})
	}

	var data map[string]string

	if err := c.BodyParser(&data); err != nil {
		c.JSON(fiber.Map{
			"pesan": err.Error(),
		})
	}

	//cek email apakah sudah ada

	var existingEmail entity.Users
	email := data["email"]

	err := database.DB.Where("email = ?", email).Find(&existingEmail).Error
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"Message": err,
		})
	}

	//cek email
	if existingEmail.Email == email {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"Message": "This Email is Register",
		})
	}

	return nil
}
