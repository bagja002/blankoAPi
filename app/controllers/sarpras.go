package controllers

import (
	"os"
	"strings"
	"template/app/entity"
	"template/pkg/config"
	"template/pkg/database"
	"template/pkg/tools"

	"github.com/gofiber/fiber/v2"
)

func CreateSarpras(c *fiber.Ctx) error {

	id_admin, _ := c.Locals("id_admin").(int)
	role, _ := c.Locals("role").(string)
	names, _ := c.Locals("name").(string)

	tools.ValidationJwtLemdik(c, role, id_admin, names)

	var request entity.Sarpras

	fotoSarpras, err := c.FormFile("FotoSarpras")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"Message": "Failed to retrieve file", "Error": err.Error()})
	}

	if err := c.BodyParser(&request); err != nil {
		c.JSON(fiber.Map{
			"pesan": err.Error(),
		})
	}

	//masukan ke sarpras untuk membuat sarpras baru

	newSarpras := entity.Sarpras{
		IdLemdik:    uint(id_admin),
		NamaSarpras: request.NamaSarpras,
		Harga:       request.Harga,
		Deskripsi:   request.Deskripsi,
		Jenis:       request.Jenis,
		FotoSarpras: strings.ReplaceAll(fotoSarpras.Filename, " ", ""),
		CreateAt:    tools.TimeNowJakarta(),
	}

	//save to database
	database.DB.Create(&newSarpras)

	//Save Foto Sarpras

	if err := c.SaveFile(fotoSarpras, "public/static/prasarana/"+strings.ReplaceAll(fotoSarpras.Filename, " ", "")); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"Message": "Failed to save file", "Error": err.Error()})
	}
	return c.JSON(fiber.Map{
		"pesan": "Sukses Menambahkan Database Sarpras",
	})
}

func GetSarpras(c *fiber.Ctx) error {
	id_admin, _ := c.Locals("id_admin").(int)
	role, _ := c.Locals("role").(string)
	names, _ := c.Locals("name").(string)

	tools.ValidationJwtLemdik(c, role, id_admin, names)

	id := c.Query("id")

	baseQuery := database.DB.Model(&entity.Sarpras{})

	// Tambahkan kondisi filter berdasarkan ID jika ada
	if id != "" {
		baseQuery = baseQuery.Where("id_sarpras = ?", id)
	}

	jenisSarpras := c.Query("jenis_sarpras")
	// Tambahkan kondisi filter berdasarkan jenis sarpras jika ada
	if jenisSarpras != "" {
		baseQuery = baseQuery.Where("jenis= ?", jenisSarpras)
	}

	// Eksekusi query dan simpan hasilnya ke dalam slice Sarpras
	var Sarpras []entity.Sarpras
	baseQuery.Where("id_lemdik = ?", id_admin).Find(&Sarpras)

	viper := config.NewViper()
	baseUrl := viper.GetString("web.baseUrl")

	for i, _ := range Sarpras {
		Sarpras[i].FotoSarpras = baseUrl + "/public/static/prasarana/" + Sarpras[i].FotoSarpras
	}

	return c.JSON(fiber.Map{
		"Pesan": "Sukses mendapatkan data",
		"data":  Sarpras,
	})
}

func UpdateSarpras(c *fiber.Ctx) error {
	id_admin, _ := c.Locals("id_admin").(int)
	role, _ := c.Locals("role").(string)
	names, _ := c.Locals("name").(string)

	tools.ValidationJwtLemdik(c, role, id_admin, names)

	//var request entity.Sarpras

	id := c.Query("id")
	fotoSarpras, _ := c.FormFile("FotoSarpras")

	var sarpras entity.Sarpras
	if err := database.DB.Where("id_sarpras = ?", id).First(&sarpras).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "sarpras not found",
		})
	}

	tx := database.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if fotoSarpras != nil {
		newPath := "public/static/prasarana/" + tools.RemoverSpaci(fotoSarpras.Filename)
		if err := c.SaveFile(fotoSarpras, newPath); err != nil {
			tx.Rollback()
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"message": "Failed to save Silabus Pelatihan",
				"error":   err.Error(),
			})
		}
		oldPath := "public/static/prasarana/" + sarpras.FotoSarpras
		if sarpras.FotoSarpras != "" {
			os.Remove(oldPath)
		}
		sarpras.FotoSarpras = tools.RemoverSpaci(fotoSarpras.Filename)
	}

	var request entity.Sarpras

	if err := c.BodyParser(&request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Failed to parse request body",
			"error":   err.Error(),
		})
	}

	sarpras.NamaSarpras = request.NamaSarpras
	sarpras.Harga = request.Harga
	sarpras.Deskripsi = request.Deskripsi
	sarpras.Jenis = request.Jenis
	sarpras.UpdateAt = tools.TimeNowJakarta()

	if err := tx.Model(&sarpras).Where("id_sarpras = ?", id).Updates(&sarpras).Error; err != nil {
		tx.Rollback()
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to update pelatihan",
			"error":   err.Error(),
		})
	}

	if err := tx.Commit().Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to commit transaction",
			"error":   err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "Sarpras updated successfully",
		"data":    sarpras,
	})
}
func DeleteSarpras(c *fiber.Ctx) error {
	// Mendapatkan nilai dari konten lokal
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

	// Mendapatkan ID dari query parameter
	id := c.Query("id")
	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"Pesan": "Bad Request: Missing ID parameter",
		})
	}

	var sarpras entity.Sarpras

	// Mencari data Sarpras berdasarkan ID
	result := database.DB.Where("id_sarpras = ?", id).Find(&sarpras)
	if result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"Pesan": "Error: " + result.Error.Error(),
		})
	}

	// Memeriksa apakah data Sarpras ditemukan
	if result.RowsAffected == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"Pesan": "Data Sarpras tidak ditemukan",
		})
	}

	// Menghapus data Sarpras
	deleteResult := database.DB.Delete(&sarpras)
	if deleteResult.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"Pesan": "Error: " + deleteResult.Error.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"Pesan": "Data Sarpras berhasil dihapus",
	})
}
