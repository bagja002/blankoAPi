package controllers

import (
	"template/app/entity"
	"template/pkg/database"
	"template/pkg/tools"

	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
)

// Register lemdik hanya bisa di daftrkan olah admin pusat
func RegisterLemdik(c *fiber.Ctx) error {

	//Pake Role Super admin/ admin pusat
	id_admin, _ := c.Locals("id_admin").(int)
	role, _ := c.Locals("role").(string)
	names, _ := c.Locals("name").(string)

	tools.ValidationJwt(c, role, id_admin, names)
	var data map[string]string

	if err := c.BodyParser(&data); err != nil {
		c.JSON(fiber.Map{
			"pesan": err.Error(),
		})
	}

	//cek email apakah sudah ada

	var existingEmail entity.Lemdiklat
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

	newLemdik := entity.Lemdiklat{
		NamaLemdik: data["nama_lemdik"],
		NoTelpon:   tools.StringToInt(data["no_telpon"]),
		Email:      email,
		Password:   tools.GeneratePassword(data["password"]),
		Alamat:     data["alamat"],
		Deskripsi:  data["deskripsi"],
		CreateAt:   tools.TimeNowJakarta(),
	}

	database.DB.Create(&newLemdik)

	return c.JSON(fiber.Map{
		"Pesan": "Sukses Membuat Lemdik",
	})
}

func LoginLemdik(c *fiber.Ctx) error {

	var data map[string]string

	if err := c.BodyParser(&data); err != nil {
		c.JSON(fiber.Map{
			"pesan": err.Error(),
		})
	}

	var users entity.Lemdiklat

	database.DB.Where("email = ? ", data["email"]).First(&users)
	if users.IdLemdik == 0 {
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

func GetLemdik(c *fiber.Ctx) error {

	//Pake Role Super admin/ admin pusat
	id_admin, _ := c.Locals("id_admin").(int)
	role, _ := c.Locals("role").(string)
	names, _ := c.Locals("name").(string)

	tools.ValidationJwtLemdik(c, role, id_admin, names)

	//Pemberian Akses Admin Pusat dan juga Admin lemdik
	var lemdik entity.Lemdiklat
	database.DB.Where("id_lemdik = ?", id_admin).Preload("Pelatihan").Find(&lemdik)

	return c.JSON(fiber.Map{
		"Pesan": "Sukses Mengambil Data",
		"data":  lemdik,
	})
}

func UpdateLemdik(c *fiber.Ctx) error {

	id_admin, _ := c.Locals("id_admin").(int)
	role, _ := c.Locals("role").(string)
	names, _ := c.Locals("name").(string)

	tools.ValidationJwtLemdik(c, role, id_admin, names)

	var lemdik entity.Lemdiklat

	database.DB.Where("id_lemdik = ?", id_admin).Find(&lemdik)

	var data map[string]string
	var existingEmail entity.Lemdiklat
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

	update := entity.Lemdiklat{
		NamaLemdik: data["nama_lemdik"],
		NoTelpon:   tools.StringToInt(data["no_telpon"]),
		Email:      email,
		Password:   tools.GeneratePassword(data["password"]),
		Alamat:     data["alamat"],
		Deskripsi:  data["deskripsi"],
		CreateAt:   tools.TimeNowJakarta(),
	}

	database.DB.Model(&lemdik).Updates(&update)

	return c.JSON(fiber.Map{
		"Pesan": "Sukses Update Data",
		"Data":  lemdik,
	})
}

func DeleteLemdik(c *fiber.Ctx) error {

	return nil
}
