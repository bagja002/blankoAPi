package controllers

import (
	"log"
	"template/app/entity"
	"template/pkg/database"
	"template/pkg/tools"

	"github.com/gofiber/fiber/v2"
)

// CreateBlankoRusak handles the creation of a new BlankoRusak record
func CreateBlankoRusak(c *fiber.Ctx) error {
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
	var request entity.BlankoRusak
	if err := c.BodyParser(&request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"Message": "Failed to parse request body", "Error": err.Error()})
	}

	// Create a new BlankoRusak record
	dataBlankoRusak := entity.BlankoRusak{
		IdBlankoKeluar: request.IdBlankoKeluar,
		NoSeri:         request.NoSeri,
		Tipe:           request.Tipe,
		Keterangan:     request.Keterangan,
		TanggalRusak:   tools.TimeNowJakarta(),
	}

	if result := database.DB.Create(&dataBlankoRusak); result.Error != nil {
		log.Printf("Failed to create BlankoRusak record: %v", result.Error)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"Message": "Failed to create BlankoRusak record", "Error": result.Error.Error()})
	}

	return c.JSON(fiber.Map{"Pesan": "Telah Berhasil Membuat Data Blanko Rusak"})
}

// GetBlankoRusak handles fetching of BlankoRusak records
func GetBlankoRusak(c *fiber.Ctx) error {
	id := c.Query("id_blanko_keluar")
	CoC := c.Query("tipe_blanko")
	var blankoRusak []entity.BlankoRusak
	query := database.DB

	if id != "" {
		query = query.Where("id_blanko_rusak = ?", id)
	}

	if CoC != "" {
		query = query.Where("tipe_blanko = ? ", CoC)
	}

	if result := query.Find(&blankoRusak); result.Error != nil {
		log.Printf("Failed to fetch BlankoRusak records: %v", result.Error)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"Message": "Failed to fetch BlankoRusak records", "Error": result.Error.Error()})
	}

	return c.JSON(fiber.Map{"Pesan": "Data Blanko Rusak Berhasil Didapatkan", "data": blankoRusak})
}

// UpdateBlankoRusak handles updating an existing BlankoRusak record
func UpdateBlankoRusak(c *fiber.Ctx) error {
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
	var request entity.BlankoRusak
	if err := c.BodyParser(&request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"Message": "Failed to parse request body", "Error": err.Error()})
	}

	// Get the ID of the BlankoRusak to update
	id := c.Query("id")
	var blankoRusak entity.BlankoRusak

	if result := database.DB.Where("id_blanko_rusak = ?", id).Find(&blankoRusak); result.Error != nil {
		log.Printf("Failed to find BlankoRusak record: %v", result.Error)
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"Message": "BlankoRusak record not found", "Error": result.Error.Error()})
	}

	// Update the BlankoRusak record
	updates := entity.BlankoRusak{
		IdBlankoKeluar: request.IdBlankoRusak,
		NoSeri:         request.NoSeri,
		Tipe:           request.Tipe,
		Keterangan:     request.Keterangan,
		TanggalRusak:   request.TanggalRusak,
	}

	if result := database.DB.Model(&blankoRusak).Updates(&updates); result.Error != nil {
		log.Printf("Failed to update BlankoRusak record: %v", result.Error)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"Message": "Failed to update BlankoRusak record", "Error": result.Error.Error()})
	}

	return c.JSON(fiber.Map{"Pesan": "Data Blanko Rusak Berhasil di ubah", "data": blankoRusak})
}

// DeleteBlankoRusak handles deletion of a BlankoRusak record
func DeleteBlankoRusak(c *fiber.Ctx) error {
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

	// Get the ID of the BlankoRusak to delete
	id := c.Query("id")
	var blankoRusak entity.BlankoRusak

	if result := database.DB.Where("id_blanko_rusak = ?", id).Find(&blankoRusak); result.Error != nil {
		log.Printf("Failed to find BlankoRusak record: %v", result.Error)
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"Message": "BlankoRusak record not found", "Error": result.Error.Error()})
	}

	// Delete the BlankoRusak record
	if result := database.DB.Delete(&blankoRusak); result.Error != nil {
		log.Printf("Failed to delete BlankoRusak record: %v", result.Error)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"Message": "Failed to delete BlankoRusak record", "Error": result.Error.Error()})
	}

	return c.JSON(fiber.Map{"Pesan": "Data Blanko Rusak Berhasil di hapus"})
}
