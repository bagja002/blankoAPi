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

	//Pemberian Akses 	Admin Pusat dan juga Admin lemdik
	var lemdik entity.Lemdiklat
	database.DB.Where("id_lemdik = ?", id_admin).Preload("Pelatihan").Find(&lemdik)

	return c.JSON(fiber.Map{
		"Pesan": "Sukses Mengambil Data",
		"data":  lemdik,
	})
}

func UpdateLemdik(c *fiber.Ctx) error {
	idAdmin, ok := c.Locals("id_admin").(int)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"Pesan": "Unauthorized: Invalid admin ID",
		})
	}
	role, ok := c.Locals("role").(string)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"Pesan": "Unauthorized: Invalid role",
		})
	}
	names, ok := c.Locals("name").(string)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"Pesan": "Unauthorized: Invalid name",
		})
	}

	// Validasi JWT
	tools.ValidationJwtLemdik(c, role, idAdmin, names)

	// Mendapatkan data yang akan diupdate dari request body
	var data map[string]string
	if err := c.BodyParser(&data); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"Pesan": "Bad Request: " + err.Error(),
		})
	}

	// Memulai transaksi
	tx := database.DB.Begin()
	if tx.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"Pesan": "Failed to begin transaction: " + tx.Error.Error(),
		})
	}

	var lemdik entity.Lemdiklat
	if err := tx.Where("id_lemdik = ?", idAdmin).First(&lemdik).Error; err != nil {
		tx.Rollback()
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"Pesan": "Data not found: " + err.Error(),
		})
	}

	// Membuat struct update dengan data yang baru
	update := entity.Lemdiklat{
		NamaLemdik:   data["nama_lemdik"],
		NoTelpon:     tools.StringToInt(data["no_telpon"]),
		Email:        data["email"],
		Password:     tools.GeneratePassword(data["password"]),
		Alamat:       data["alamat"],
		Deskripsi:    data["deskripsi"],
		LastNosertif: data["no_last_sertifikat"],
		CreateAt:     tools.TimeNowJakarta(),
	}

	// Menambahkan record sertifikat jika ada update pada no_last_sertifikat
	if update.LastNosertif != "" {
		newLastRecord := entity.Sertifikat{
			IdLemdik:    uint(idAdmin),
			NoSertfikat: data["no_last_sertifikat"],
		}
		if err := tx.Create(&newLastRecord).Error; err != nil {
			tx.Rollback()
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"Pesan": "Failed to create new certificate record: " + err.Error(),
			})
		}
	}

	// Melakukan update pada data Lemdiklat
	if err := tx.Model(&lemdik).Updates(&update).Error; err != nil {
		tx.Rollback()
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"Pesan": "Failed to update Lemdiklat: " + err.Error(),
		})
	}

	// Commit transaksi
	if err := tx.Commit().Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"Pesan": "Failed to commit transaction: " + err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"Pesan": "Sukses Update Data",
		"Data":  lemdik,
	})
}

func DeleteLemdik(c *fiber.Ctx) error {

	return nil
}

func LastNomorSertifBalai(c *fiber.Ctx) error {

	id_admin, _ := c.Locals("id_admin").(int)
	role, _ := c.Locals("role").(string)
	names, _ := c.Locals("name").(string)

	tools.ValidationJwtLemdik(c, role, id_admin, names)

	var data map[string]string

	err := c.BodyParser(&data)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"Pesan": "Gagal Mendapatkan Body data",
		})
	}

	newLastRecord := entity.Sertifikat{
		IdLemdik:    uint(id_admin),
		NoSertfikat: data["no_last_sertifikat"],
	}

	database.DB.Create(&newLastRecord)

	return c.JSON(fiber.Map{
		"Pesan": "Sukses Menambahkan data balai",
	})
}
