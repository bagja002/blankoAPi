package controllers

import (
	"os"
	"template/app/entity"
	"template/pkg/database"
	"template/pkg/tools"

	"github.com/gofiber/fiber/v2"
)

func CreatePengirimanSertifikat(c *fiber.Ctx) error {
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

	// Validate JWT token
	tools.ValidationJwtLemdik(c, role, idAdmin, name)

	var request entity.PengirimanSertifikat
	if err := c.BodyParser(&request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"Message": "Failed to parse request body", "Error": err.Error()})
	}

	// Handle file uploads
	buktiResi, err := c.FormFile("bukti_resi")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Bukti Resi file not found"})
	}

	ttdTerima, _ := c.FormFile("ttd_terima")

	buktiPengiriman, _ := c.FormFile("bukti_pengiriman")

	// Create Pengiriman Sertifikat record
	newData := entity.PengirimanSertifikat{
		NamaPenerima:              request.NamaPenerima,
		NomorTelpon:               request.NomorTelpon,
		Alamat:                    request.Alamat,
		NoResi:                    request.NoResi,
		BuktiResi:                 buktiResi.Filename,
		NominalPengiriman:         request.NominalPengiriman,
		TtdTerimaPengiriman:       ttdTerima.Filename,
		BuktiPengirimanSertifikat: buktiPengiriman.Filename,
		BuktiPenerimaanSertikat:   request.BuktiPenerimaanSertikat,
		ListSertifikatDikirimkan:  request.ListSertifikatDikirimkan,
		CreateAt:                  tools.TimeNowJakarta(),
		UpdateAt:                  tools.TimeNowJakarta(),
		Status:                    request.Status,
	}

	// Save data to the database
	database.DB.Create(&newData)

	// Save files
	if err := c.SaveFile(buktiResi, "public/static/bukti-resi/"+tools.RemoverSpaci(buktiResi.Filename)); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to save file"})
	}
	if err := c.SaveFile(ttdTerima, "public/static/ttd-penerima/"+tools.RemoverSpaci(ttdTerima.Filename)); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to save file"})
	}
	if err := c.SaveFile(buktiPengiriman, "public/static/bukti-pengiriman-sertifikat/"+tools.RemoverSpaci(buktiPengiriman.Filename)); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to save file"})
	}

	return c.JSON(fiber.Map{
		"Pesan": "Berhasil Membuat Pengiriman Sertifikat",
	})
}

func GetPengirimanSertifikat(c *fiber.Ctx) error {
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

	// Validate JWT token
	tools.ValidationJwtLemdik(c, role, idAdmin, name)

	// Get the ID from the URL parameters
	id := c.Query("id")
	var pengiriman []entity.PengirimanSertifikat

	baseQuery := database.DB
	if id != "" {
		baseQuery = baseQuery.Where("id_pengiriman_sertifikat =?", id)
	}

	// Find the data based on ID
	if err := baseQuery.Find(&pengiriman).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"Message": "Data not found"})
	}

	return c.JSON(fiber.Map{
		"Pesan": "Berhasl get",
		"data":  pengiriman,
	})
}

func UpdatePengirimanSertifikat(c *fiber.Ctx) error {
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

	// Validate JWT token
	tools.ValidationJwtLemdik(c, role, idAdmin, name)

	var request entity.PengirimanSertifikat
	if err := c.BodyParser(&request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"Message": "Failed to parse request body", "Error": err.Error()})
	}

	// Get the ID from the URL parameters
	id := c.Params("id")
	var pengiriman entity.PengirimanSertifikat

	// Find the data based on ID
	if err := database.DB.First(&pengiriman, id).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"Message": "Data not found"})
	}

	buktiResi, _ := c.FormFile("bukti_resi")

	ttdTerima, _ := c.FormFile("ttd_terima")

	buktiPengiriman, _ := c.FormFile("bukti_pengiriman")

	tx := database.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if ttdTerima != nil {
		newPath := "public/static/ttd-penerima/" + tools.RemoverSpaci(ttdTerima.Filename)
		if err := c.SaveFile(ttdTerima, newPath); err != nil {
			tx.Rollback()
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Gagal Menyimpan File Dulu",
				"path":  newPath,
			})
		}

		oldPath := "public/static/ttd-penerima/" + pengiriman.TtdTerimaPengiriman
		if pengiriman.TtdTerimaPengiriman != "" {
			os.Remove(oldPath)
		}
	}

	if buktiResi != nil {
		newPath := "public/static/bukti-resi/" + tools.RemoverSpaci(buktiResi.Filename)
		if err := c.SaveFile(buktiResi, newPath); err != nil {
			tx.Rollback()
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Gagal Menyimpan File Dulu",
				"path":  newPath,
			})
		}

		oldPath := "public/static/bukti-resi/" + pengiriman.BuktiResi
		if pengiriman.BuktiResi != "" {
			os.Remove(oldPath)
		}
	}

	if buktiPengiriman != nil {
		newPath := "public/static/bukti-pengiriman-sertifikat/" + tools.RemoverSpaci(buktiPengiriman.Filename)
		if err := c.SaveFile(buktiPengiriman, newPath); err != nil {
			tx.Rollback()
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Gagal Menyimpan File Dulu",
				"path":  newPath,
			})
		}

		oldPath := "public/static/bukti-pegiriman-sertifikat/" + pengiriman.BuktiPengirimanSertifikat
		if pengiriman.BuktiPengirimanSertifikat != "" {
			os.Remove(oldPath)
		}
	}
	// Update the data
	pengiriman.NamaPenerima = request.NamaPenerima
	pengiriman.NomorTelpon = request.NomorTelpon
	pengiriman.Alamat = request.Alamat
	pengiriman.NoResi = request.NoResi
	pengiriman.NominalPengiriman = request.NominalPengiriman
	pengiriman.ListSertifikatDikirimkan = request.ListSertifikatDikirimkan
	pengiriman.Status = request.Status
	pengiriman.UpdateAt = tools.TimeNowJakarta()

	// Save the updated data to the database
	if err := database.DB.Save(&pengiriman).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"Message": "Failed to update data"})
	}

	return c.JSON(fiber.Map{"Message": "Pengiriman Sertifikat updated successfully"})
}

func DeletePengirimanSertifikat(c *fiber.Ctx) error {
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

	// Validate JWT token
	tools.ValidationJwtLemdik(c, role, idAdmin, name)

	// Get the ID from the URL parameters
	id := c.Params("id")
	var pengiriman entity.PengirimanSertifikat

	// Find the data based on ID
	if err := database.DB.First(&pengiriman, id).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"Message": "Data not found"})
	}

	// Delete the data
	if err := database.DB.Delete(&pengiriman).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"Message": "Failed to delete data"})
	}

	return c.JSON(fiber.Map{"Message": "Pengiriman Sertifikat deleted successfully"})
}
