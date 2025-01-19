package controllers

import (
	"fmt"
	"template/pkg/database"
	"time"

	"github.com/gofiber/fiber/v2"
)

type Sertifikat struct {
	JenisSertifikat  string `gorm:"column:s_jenis_sertifikat"`
	JumlahSertifikat int    `gorm:"column:jumlah_sertifikat"`
}

// CustomResult struct to map the aggregated query result
type CustomResult struct {
	JenisSertifikat  string `json:"jenis_sertifikat"`
	JumlahSertifikat int    `json:"jumlah_sertifikat"`
}

type ResponseDatas struct {
	NamaLembaga        string `json:"nama_lembaga"`
	SubJenisPendidikan string `json:"sub_jenis_pendidikan"`
	Jumlah             int64  `json:"jumlah"`
}

type ResponseDataCOc struct {
	NamaUnit         string `json:"nama_unit"`
	JenisSertifikasi string `json:"jenis_sertifikasi"`
	Jumlah           int64  `json:"jumlah"`
}

type ResponseData struct {
	DID                uint   `json:"d_id"`
	LID                uint   `json:"l_id"`
	PLID               uint   `json:"pl_id"`
	NamaLembaga        string `json:"nama_lembaga"`
	SubJenisPendidikan string `json:"sub_jenis_pendidikan"`
	Lokasi             string `json:"lokasi"`
	Jumlah             int64  `json:"jumlah"`
}

// get jumlah COP dan COC totalan
func GetDataSertifikat(c *fiber.Ctx) error {
	// Default date range and query parameter values
	startDate := c.Query("start_date", "2024-06-01")
	endDate := c.Query("end_date", "2024-12-31")
	isPrint := c.Query("is_print", "1")      // Default value is 1
	typeBlanko := c.Query("type_blanko", "") // Expected values: "COP" or "COC"

	// Validate is_print input (must be 0 or 1)
	if isPrint != "0" && isPrint != "1" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid value for is_print. Must be '0' or '1'.",
		})
	}

	// Validate type_blanko input
	if typeBlanko != "COP" && typeBlanko != "COC" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid value for type_blanko. Must be 'COP' or 'COC'.",
		})
	}

	// Parse the date range
	startDateTime, err := time.Parse("2006-01-02", startDate)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid start_date format. Expected format: YYYY-MM-DD.",
		})
	}
	endDateTime, err := time.Parse("2006-01-02", endDate)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid end_date format. Expected format: YYYY-MM-DD.",
		})
	}

	// Choose the correct join based on type_blanko
	var joinCondition string
	if typeBlanko == "COP" {
		joinCondition = "JOIN master_diklat d ON s.d_id = d.d_id"
	} else if typeBlanko == "COC" {
		joinCondition = "JOIN rencana_ujian d ON s.d_id = d.ru_id"
	}

	// Execute the query
	var results []CustomResult
	err = database.DB1.Table("sertifikat AS s").
		Select("s.s_jenis_sertifikat AS jenis_sertifikat, COUNT(*) AS jumlah_sertifikat").
		Joins(joinCondition).
		Where("s.isprint = ? AND s.created_on BETWEEN ? AND ?", isPrint, startDateTime, endDateTime).
		Group("s.s_jenis_sertifikat").
		Order("jumlah_sertifikat DESC").
		Scan(&results).Error
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch data from the database: " + err.Error(),
		})
	}

	// Return the query results as JSON
	return c.JSON(fiber.Map{
		"message": "Successfully retrieved data from the database.",
		"data":    results,
	})
}

// Get Diklat dan Solikasi dan Juga Dimana Lokasi Diklatnya
// jadi Ini Yang harus di jabarkan di diklat
func GetDataBalaiSertifikatLokasi(c *fiber.Ctx) error {
	startDate := c.Query("start_date", "2024-06-01")
	endDate := c.Query("end_date", "2024-12-31")
	isPrint := c.Query("is_print", "1")
	typeBlanko := c.Query("type_blanko", "")

	// Validate is_print input (must be 0 or 1)
	if isPrint != "0" && isPrint != "1" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid value for is_print. Must be '0' or '1'.",
		})
	}

	// Validate type_blanko input
	if typeBlanko != "COP" && typeBlanko != "COC" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid value for type_blanko. Must be 'COP' or 'COC'.",
		})
	}

	// Parse the date range
	startDateTime, err := time.Parse("2006-01-02", startDate)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid start_date format. Expected format: YYYY-MM-DD.",
		})
	}
	endDateTime, err := time.Parse("2006-01-02", endDate)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid end_date format. Expected format: YYYY-MM-DD.",
		})
	}

	var joinCondition string
	var joinCondition2 string
	var joinCondition3 string
	var selectAwal string
	var GrouBy string
	if typeBlanko == "COP" {
		joinCondition = "JOIN master_diklat d ON s.d_id = d.d_id"
		joinCondition2 = "JOIN master_lembaga ml ON d.l_id = ml.l_id"
		joinCondition3 = "JOIN master_profil_lembaga pl ON ml.pl_id = pl.pl_id"
		GrouBy = " pl.pl_nama_lembaga, d.d_sub_jenis_pendidikan, d.d_lokasi"
		selectAwal = `pl.pl_nama_lembaga as nama_lembaga, 
			d.d_sub_jenis_pendidikan as sub_jenis_pendidikan, 
			d.d_lokasi as lokasi, 
			COUNT(*) AS jumlah`
	} else if typeBlanko == "COC" {
		joinCondition = "JOIN rencana_ujian d ON s.d_id = d.ru_id"
		joinCondition2 = "JOIN master_unit_kerja ml ON d.ru_unit_kerja = ml.uk_id"
		joinCondition3 = ""
		GrouBy = "ml.uk_nama, d.ru_jenis_setifikasi, d.ru_tempat_ujian"
		selectAwal = `ml.uk_nama as nama_lembaga, 
			d.ru_jenis_setifikasi as sub_jenis_pendidikan, 
			d.ru_tempat_ujian as lokasi, 
			COUNT(*) AS jumlah`
	}

	var results []ResponseData
	query := database.DB1.Table("sertifikat s").
		Select(selectAwal).
		Joins(joinCondition).
		Joins(joinCondition2)

	if typeBlanko != "COC" {
		query = query.Joins(joinCondition3)
	}

	err = query.
		Where("s.created_on BETWEEN ? AND ? AND s.isprint = ?", startDateTime, endDateTime, isPrint).
		Group(GrouBy).
		Order("jumlah DESC").
		Scan(&results).Error

	fmt.Println("querynya", query)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "Failed to fetch data",
			"err":   err,
		})
	}

	return c.JSON(results)
}

// Get Jumlah Per Balai per COC dan COP nya
func GetDataBalaiSertifikat(c *fiber.Ctx) error {
	startDate := c.Query("start_date", "2024-06-01")
	endDate := c.Query("end_date", "2024-12-31")
	isPrint := c.Query("is_print", "1")

	var results1 []ResponseData
	var results2 []ResponseData

	err1 := database.DB1.Table("sertifikat s").
		Select(`pl.pl_nama_lembaga as nama_lembaga, 
			d.d_sub_jenis_pendidikan as sub_jenis_pendidikan, d.d_lokasi as lokasi, COUNT(*) AS jumlah`).
		Joins("JOIN master_diklat d ON s.d_id = d.d_id").
		Joins("JOIN master_lembaga ml ON d.l_id = ml.l_id").
		Joins("JOIN master_profil_lembaga pl ON ml.pl_id = pl.pl_id").
		Where("s.isprint = ? AND s.created_on BETWEEN ? AND ?", isPrint, startDate, endDate).
		Group("pl.pl_nama_lembaga, d.d_sub_jenis_pendidikan, d.d_lokasi").
		Order("pl.pl_nama_lembaga, jumlah DESC").
		Scan(&results1).Error

	if err1 != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "Failed to fetch data for Balai Sertifikat Cop",
			"err":   err1,
		})
	}

	err2 := database.DB1.Table("sertifikat s").
		Select(`ml.uk_nama as nama_lembaga, 
			d.ru_jenis_setifikasi as sub_jenis_pendidikan, d.ru_tempat_ujian as lokasi, COUNT(*) AS jumlah`).
		Joins("JOIN rencana_ujian d ON s.d_id = d.ru_id").
		Joins("JOIN master_unit_kerja ml ON d.ru_unit_kerja = ml.uk_id").
		Where("s.isprint = 1 AND s.created_on BETWEEN ? AND ?", startDate, endDate).
		Group("ml.uk_nama, d.ru_jenis_setifikasi, d.ru_tempat_ujian").
		Order("ml.uk_nama ASC, jumlah DESC").
		Scan(&results2).Error

	if err2 != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "Failed to fetch data for Balai Sertifikat CoC",
			"err":   err2,
		})
	}

	return c.JSON(fiber.Map{
		"balai_sertifikat_cop": results1,
		"balai_sertifikat_coc": results2,
	})
}
