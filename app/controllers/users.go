package controllers

import (
	//"backend-elaut/app/entity"

	"template/app/entity"
	"template/app/models"
	"template/pkg/database"
	"template/pkg/tools"

	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
)

func CreateUser(c *fiber.Ctx)error{

	return c.SendString("Testing")
}



func Login(c *fiber.Ctx) error {

	var data map[string]string

	if err := c.BodyParser(&data); err != nil {
		c.JSON(fiber.Map{
			"pesan": err.Error(),
		})
	}

	var users entity.Users

	database.DB.Where("username = ? ", data["username"]).First(&users)
	if users.IdUsers == 0 {
		c.Status(fiber.StatusUnauthorized)
		return c.JSON(fiber.Map{
			"pesan": "Username tidak di temukan",
		})
	}

	if err := bcrypt.CompareHashAndPassword([]byte("Password database"), []byte(data["password"])); err != nil {
		c.Status(fiber.StatusUnauthorized)
		return c.JSON(fiber.Map{
			"pesan": "Incorrect password!",
		})
	} else {
		t:= tools.GenerateToken(models.User{})
		return c.JSON(fiber.Map{
			"t":t,
		})
	}

	
}

