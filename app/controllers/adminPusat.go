package controllers

import (
	"template/app/entity"
	//"template/app/models"
	"template/pkg/database"
	"template/pkg/tools"

	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
)

func SuperAdminLogin(c *fiber.Ctx)error{
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
		t:= tools.GenerateToken(users)
		return c.JSON(fiber.Map{
			"t":t,
		})
	}

	


}

func CreateAdminPusat(c *fiber.Ctx)error{

	


		return nil
	}