package controllers

import (
	"fmt"
	"log"
	"path/filepath"
	"template/app/entity"
	"template/pkg/database"
	"template/pkg/tools"

	"github.com/gofiber/fiber/v2"
	"github.com/xuri/excelize/v2"
)

func CreateMateriPelatihan(c *fiber.Ctx) error {

	id_admin, _ := c.Locals("id_admin").(int)
	role, _ := c.Locals("role").(string)
	names, _ := c.Locals("name").(string)

	tools.ValidationJwtLemdik(c, role, id_admin, names)

	var data map[string]string

	if err := c.BodyParser(&data); err != nil {
		c.JSON(fiber.Map{
			"pesan": err.Error(),
		})
	}

	idPelatihan := c.Query("id_pelatihan")

	var pelatihan entity.Pelatihan

	database.DB.Where("id_pelatihan = ? ", idPelatihan).Find(&pelatihan)
	if pelatihan.IdPelatihan == 0 {
		return c.Status(400).JSON(fiber.Map{
			"Pesan": "tidak ada pelatihan",
		})
	}

	newMateriPelatihan := entity.MateriPelatihan{
		IdPelatihan: uint(tools.StringToInt(idPelatihan)),
		NamaMateri:  data["nama_materi"],
		Deskripsi:   data["deskripsi"],
		JamTeory:    data["jam_teory"],
		JamPraktek:  data["jam_praktek"],
	}

	database.DB.Create(&newMateriPelatihan)

	return c.JSON(fiber.Map{
		"Pesan": "Sukses Membuat Materi Pelatihan",
		"data":  newMateriPelatihan,
	})
}

func ImportMateriPelatihan(c *fiber.Ctx) error {
	id_admin, _ := c.Locals("id_admin").(int)
	role, _ := c.Locals("role").(string)
	names, _ := c.Locals("name").(string)

	tools.ValidationJwtLemdik(c, role, id_admin, names)

	// Membaca file Excel dari request
	file, err := c.FormFile("file")
	if err != nil {
		return err
	}

	// Memeriksa ekstensi file
	ext := filepath.Ext(file.Filename)
	if ext != ".xlsx" && ext != ".xls" {
		return fmt.Errorf("file harus berupa file Excel (.xlsx atau .xls)")
	}

	// Membuka file Excel
	excelFile, err := file.Open()
	if err != nil {
		return err
	}
	defer excelFile.Close()

	// Membaca file Excel menggunakan excelize
	f, err := excelize.OpenReader(excelFile)
	if err != nil {
		return err
	}

	// Mendapatkan nama semua sheet dalam file Excel
	sheets := f.GetSheetList()

	// Membaca data dari sheet pertama
	rows, err := f.GetRows(sheets[0])
	if err != nil {
		return err
	}

	var materiPelaihan []entity.MateriPelatihan

	var reques entity.MateriPelatihan

	if err := c.BodyParser(&reques); err != nil {
		// Log the error for debugging purposes
		log.Printf("Error parsing request body: %v", err)

		// Respond with an appropriate error message and status code
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "Failed to parse request body",
			"message": err.Error(),
		})
	}
	idPelatihan := reques.IdPelatihan
	fmt.Println(idPelatihan)
	for _, rowsMateri := range rows[1:] {

		materiPelatihans := entity.MateriPelatihan{}

		for i, columnName := range rows[0] {
			if i >= len(rowsMateri) {
				continue
			}

			switch columnName {
			case "nama_materi":
				materiPelatihans.NamaMateri = rowsMateri[i]
			case "deskripsi":
				materiPelatihans.Deskripsi = rowsMateri[i]
			case "jam_teory":
				materiPelatihans.JamTeory = rowsMateri[i]
			case "jam_praktek":
				materiPelatihans.JamPraktek = rowsMateri[i]
			}
			materiPelatihans.IdPelatihan = uint(idPelatihan)
			materiPelatihans.CreateAt = tools.TimeNowJakarta()
		}
		materiPelaihan = append(materiPelaihan, materiPelatihans)

	}

	for _, AllMateri := range materiPelaihan {
		if err := database.DB.Create(&AllMateri).Error; err != nil {
			// Log the error for debugging purposes
			fmt.Println("Gagal")
			return c.JSON(err)

		}
	}

	return c.JSON(fiber.Map{
		"Pesan": "Sukses Upload Data Materi Pelatihan",
	})
}
