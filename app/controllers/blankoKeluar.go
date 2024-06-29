package controllers

import (
	"log"
	"template/app/entity"
	"template/pkg/database"
	"template/pkg/tools"

	"github.com/gofiber/fiber/v2"
)

// CreateBlankoKeluar handles the creation of a new BlankoKeluar record
func CreateBlankoKeluar(c *fiber.Ctx) error {
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
	var request entity.BlankoKeluar
	if err := c.BodyParser(&request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"Message": "Failed to parse request body", "Error": err.Error()})
	}

	// Create a new BlankoKeluar record
	dataBlankoKeluar := entity.BlankoKeluar{
		IdBlanko:              request.IdBlanko,
		TipeBlanko:            request.TipeBlanko,
		TanggalKeluar:         request.TanggalKeluar,
		NamaLemdiklat:         request.NamaLemdiklat,
		NamaPelaksana:         request.NamaPelaksana,
		TanggalPermohonan:     request.TanggalPermohonan,
		LinkPermohonan:        request.LinkPermohonan,
		NamaProgram:           request.NamaProgram,
		TanggalPelaksanaan:    request.TanggalPelaksanaan,
		JumlahPesertaLulus:    request.JumlahPesertaLulus,
		JumlahBlankoDiajukan:  request.JumlahBlankoDiajukan,
		JumlahBlankoDisetujui: request.JumlahBlankoDisetujui,
		NoSeriBlanko:          request.NoSeriBlanko,
		Status:                request.Status,
		IsSd:                  request.IsSd,
		IsCetak:               request.IsCetak,
		TipePengambilan:       request.TipePengambilan,
		PetugasYangMenerima:   request.PetugasYangMenerima,
		PetugasYangMemberi:    request.PetugasYangMemberi,
		LinkDataDukung:        request.LinkDataDukung,
		CreatedAt:             tools.TimeNowJakarta(),
		UpdatedAt:             tools.TimeNowJakarta(),
		Keterangan:            request.Keterangan,
	}

	if result := database.DB.Create(&dataBlankoKeluar); result.Error != nil {
		log.Printf("Failed to create BlankoKeluar record: %v", result.Error)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"Message": "Failed to create BlankoKeluar record", "Error": result.Error.Error()})
	}

	return c.JSON(fiber.Map{"Pesan": "Telah Berhasil Membuat Data Blanko Keluar"})
}

// GetBlankoKeluar handles fetching of BlankoKeluar records
func GetBlankoKeluar(c *fiber.Ctx) error {
	id := c.Query("id_blanko_keluar")

	var blankoKeluar []entity.BlankoKeluar
	query := database.DB

	if id != "" {
		query = query.Where("id_blanko_keluar = ?", id)
	}

	if result := query.Find(&blankoKeluar); result.Error != nil {
		log.Printf("Failed to fetch BlankoKeluar records: %v", result.Error)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"Message": "Failed to fetch BlankoKeluar records", "Error": result.Error.Error()})
	}

	return c.JSON(fiber.Map{"Pesan": "Data Blanko Keluar Berhasil Didapatkan", "data": blankoKeluar})
}

// UpdateBlankoKeluar handles updating an existing BlankoKeluar record
func UpdateBlankoKeluar(c *fiber.Ctx) error {
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
	var request entity.BlankoKeluar
	if err := c.BodyParser(&request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"Message": "Failed to parse request body", "Error": err.Error()})
	}

	// Get the ID of the BlankoKeluar to update
	id := c.Query("id")
	var blankoKeluar entity.BlankoKeluar

	if result := database.DB.Where("id_blanko_keluar = ?", id).Find(&blankoKeluar); result.Error != nil {
		log.Printf("Failed to find BlankoKeluar record: %v", result.Error)
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"Message": "BlankoKeluar record not found", "Error": result.Error.Error()})
	}

	// Update the BlankoKeluar record
	updates := entity.BlankoKeluar{
		IdBlanko:              request.IdBlanko,
		TipeBlanko:            request.TipeBlanko,
		TanggalKeluar:         request.TanggalKeluar,
		NamaLemdiklat:         request.NamaLemdiklat,
		NamaPelaksana:         request.NamaPelaksana,
		TanggalPermohonan:     request.TanggalPermohonan,
		LinkPermohonan:        request.LinkPermohonan,
		NamaProgram:           request.NamaProgram,
		TanggalPelaksanaan:    request.TanggalPelaksanaan,
		JumlahPesertaLulus:    request.JumlahPesertaLulus,
		JumlahBlankoDiajukan:  request.JumlahBlankoDiajukan,
		JumlahBlankoDisetujui: request.JumlahBlankoDisetujui,
		NoSeriBlanko:          request.NoSeriBlanko,
		Status:                request.Status,
		IsSd:                  request.IsSd,
		IsCetak:               request.IsCetak,
		TipePengambilan:       request.TipePengambilan,
		PetugasYangMenerima:   request.PetugasYangMenerima,
		PetugasYangMemberi:    request.PetugasYangMemberi,
		LinkDataDukung:        request.LinkDataDukung,
		UpdatedAt:             tools.TimeNowJakarta(),
		Keterangan:            request.Keterangan,
	}

	if result := database.DB.Model(&blankoKeluar).Updates(&updates); result.Error != nil {
		log.Printf("Failed to update BlankoKeluar record: %v", result.Error)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"Message": "Failed to update BlankoKeluar record", "Error": result.Error.Error()})
	}

	return c.JSON(fiber.Map{"Pesan": "Data Blanko Keluar Berhasil di ubah", "data": blankoKeluar})
}

// DeleteBlankoKeluar handles deletion of a BlankoKeluar record
func DeleteBlankoKeluar(c *fiber.Ctx) error {
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

	// Get the ID of the BlankoKeluar to delete
	id := c.Query("id")
	var blankoKeluar entity.BlankoKeluar

	if result := database.DB.Where("id_blanko_keluar = ?", id).Find(&blankoKeluar); result.Error != nil {
		log.Printf("Failed to find BlankoKeluar record: %v", result.Error)
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"Message": "BlankoKeluar record not found", "Error": result.Error.Error()})
	}

	// Delete the BlankoKeluar record
	if result := database.DB.Delete(&blankoKeluar); result.Error != nil {
		log.Printf("Failed to delete BlankoKeluar record: %v", result.Error)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"Message": "Failed to delete BlankoKeluar record", "Error": result.Error.Error()})
	}

	return c.JSON(fiber.Map{"Pesan": "Data Blanko Keluar Berhasil di hapus"})
}
