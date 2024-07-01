package controllers

import (
	"log"
	"template/app/entity"
	"template/pkg/database"
	"template/pkg/tools"

	"github.com/gofiber/fiber/v2"
)

// CreteDataBlanko handles the creation of a new Blanko record
func CreteDataBlanko(c *fiber.Ctx) error {
	// Get role and admin details from JWT token
	idAdmin, ok := c.Locals("id_admin").(int)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"Message": "Unauthorized"})
	}
	role, ok := c.Locals("role").(string)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"Message": "Unauthorized"})
	}
	name, ok := c.Locals("name").(string)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"Message": "Unauthorized"})
	}

	// Validate the JWT token
	tools.ValidationJwtLemdik(c, role, idAdmin, name)

	// Parse the request body
	var data map[string]string
	if err := c.BodyParser(&data); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"Message": "Failed to parse request body", "Error": err.Error()})
	}

	// Create a new Blanko record
	dataBlanko := entity.Blanko{
		Jumlah:           tools.StringToInt(data["jumlah"]),
		NoSeri:           data["no_seri"],
		TipeBlanko:       data["tipe_blanko"],
		TanggalPengadaan: data["tanggal_pengadaan"],
		JumlahPengadaan:  tools.StringToInt(data["jumlah_pengadaan"]),
		CreateAt:         tools.TimeNowJakarta(),
	}
	if result := database.DB.Create(&dataBlanko); result.Error != nil {
		log.Printf("Failed to create Blanko record: %v", result.Error)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"Message": "Failed to create Blanko record", "Error": result.Error.Error()})
	}

	return c.JSON(fiber.Map{"Pesan": "Telah Berhasil Membuat Data Blanko"})
}

// GetBlanko handles fetching of Blanko records
func GetBlanko(c *fiber.Ctx) error {
	id := c.Query("id_blanko")

	var blanko []entity.Blanko
	query := database.DB

	if id != "" {
		query = query.Where("id_blanko = ?", id)
	}

	if result := query.Find(&blanko); result.Error != nil {
		log.Printf("Failed to fetch Blanko records: %v", result.Error)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"Message": "Failed to fetch Blanko records", "Error": result.Error.Error()})
	}

	return c.JSON(fiber.Map{"Pesan": "Data Blanko Berhasil Didapatkan", "data": blanko})
}

// UpdateBlanko handles updating an existing Blanko record
func UpdateBlanko(c *fiber.Ctx) error {
	// Get role and admin details from JWT token
	idAdmin, ok := c.Locals("id_admin").(int)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"Message": "Unauthorized"})
	}
	role, ok := c.Locals("role").(string)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"Message": "Unauthorized"})
	}
	name, ok := c.Locals("name").(string)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"Message": "Unauthorized"})
	}

	// Validate the JWT token
	tools.ValidationJwtLemdik(c, role, idAdmin, name)

	// Parse the request body
	var request entity.Blanko
	if err := c.BodyParser(&request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"Message": "Failed to parse request body", "Error": err.Error()})
	}

	// Get the ID of the Blanko to update
	id := c.Query("id")
	var blanko entity.Blanko

	if result := database.DB.Where("id_blanko = ?", id).Find(&blanko); result.Error != nil {
		log.Printf("Failed to find Blanko record: %v", result.Error)
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"Message": "Blanko record not found", "Error": result.Error.Error()})
	}

	// Update the Blanko record
	updates := entity.Blanko{
		Jumlah:           request.Jumlah,
		NoSeri:           request.NoSeri,
		TipeBlanko:       request.TipeBlanko,
		TanggalPengadaan: request.TanggalPengadaan,
		JumlahPengadaan:  request.JumlahPengadaan,
		CreateAt:         tools.TimeNowJakarta(),
	}

	if result := database.DB.Model(&blanko).Updates(&updates); result.Error != nil {
		log.Printf("Failed to update Blanko record: %v", result.Error)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"Message": "Failed to update Blanko record", "Error": result.Error.Error()})
	}

	return c.JSON(fiber.Map{"Pesan": "Data Blanko Berhasil di ubah", "data": blanko})
}

// DeleteBlanko handles deletion of a Blanko record
func DeleteBlanko(c *fiber.Ctx) error {
	// Get role and admin details from JWT token
	idAdmin, ok := c.Locals("id_admin").(int)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"Message": "Unauthorized"})
	}
	role, ok := c.Locals("role").(string)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"Message": "Unauthorized"})
	}
	name, ok := c.Locals("name").(string)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"Message": "Unauthorized"})
	}

	// Validate the JWT token
	tools.ValidationJwtLemdik(c, role, idAdmin, name)

	// Get the ID of the Blanko to delete
	id := c.Query("id")
	var blanko entity.Blanko

	if result := database.DB.Where("id_blanko = ?", id).Find(&blanko); result.Error != nil {
		log.Printf("Failed to find Blanko record: %v", result.Error)
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"Message": "Blanko record not found", "Error": result.Error.Error()})
	}

	// Delete the Blanko record
	if result := database.DB.Delete(&blanko); result.Error != nil {
		log.Printf("Failed to delete Blanko record: %v", result.Error)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"Message": "Failed to delete Blanko record", "Error": result.Error.Error()})
	}

	return c.JSON(fiber.Map{"Pesan": "Data Blanko Berhasil di hapus"})
}
