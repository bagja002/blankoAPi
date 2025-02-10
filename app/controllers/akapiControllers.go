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

func GetDataByNameUserSertifika(c *fiber.Ctx) error {

	startDate := c.Query("start_date", "2024-06-01")
	endDate := c.Query("end_date", "2024-12-31")
	isPrint := c.Query("is_print", "1")
	typeBlanko := c.Query("type_blanko", "")

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
	var Selec string
	var Condition1 string
	var Condition2 string
	var Where string

	type Result struct {
		NomorSertifikat   string    `json:"nomor_sertifikat"`
		NomorBlanko       string    `json:"nomor_blanko"`
		JenisDiklat       string    `json:"jenis_diklat"`
		TempatDiklat      string    `json:"tempat_diklat"`
		TanggalSertifikat time.Time `json:"tanggal_sertifikat"`
		NamaLengkap       string    `json:"nama_lengkap"`
		TempatLahir       string    `json:"tempat_lahir"`
		NIK               string    `json:"nik"`
		TanggalLahir      time.Time `json:"tanggal_lahir"`
		Alamat            string    `json:"alamat"`
		JenisSertifikat   string    `json:"jenis_sertifikat"`
		Lokasi            string    `json:"lokasi"`
	}

	var results []Result

	if typeBlanko == "COC" {
		Selec = "s.s_nomor_sertifikat as nomor_sertifikat, s.s_serial_no as nomor_blanko, ru.ru_jenis_setifikasi as jenis_diklat, ru.ru_tempat_ujian as tempat_diklat, s.s_tanggal as tangal_sertifikat, a.nama_lengkap as nama_lengkap, a.tempat_lahir as tempat_lahir, a.nik, a.tanggal_lahir, a.alamat, s.s_jenis_sertifikat"
		Condition1 = "JOIN anggota a ON s.anggota_id = a.id"
		Condition2 = "JOIN rencana_ujian ru ON s.d_id = ru.ru_id"
		Where = `s.isprint = ? AND s.created_on BETWEEN ? AND ?`
	} else if typeBlanko == "COP" {
		Selec = "s.s_nomor_sertifikat as nomor_sertifikat, s.s_serial_no as nomor_blanko, d.d_nama as jenis_diklat, d.d_tempat as tempat_diklat, s.s_tanggal as tanggal_sertifikat, a.nama_lengkap as nama_lengkap, a.tempat_lahir as tempat_lahir, a.nik, a.tanggal_lahir, a.alamat, s.s_jenis_sertifikat"
		Condition1 = "JOIN anggota a ON s.anggota_id = a.id"
		Condition2 = " JOIN master_diklat d ON s.d_id = d.d_id"
		Where = "s.isprint = ? AND s.created_on BETWEEN ? AND ?"
	}

	baseQuery := database.DB1.Table("sertifikat AS s").
		Select(Selec, isPrint, startDateTime, endDateTime).
		Joins(Condition1).
		Joins(Condition2).
		Where(Where).
		Order("s.s_serial_no")

	baseQuery.Scan(&results)

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

type CombinedResponse struct {
	DataLembaga   []LembagaResponse   `json:"data_lembaga"`
	DataUnitKerja []UnitKerjaResponse `json:"data_unit_kerja"`
}

type RespostSertifikat struct {
	DataCOP []SertifikatPerbalai `json:"data"`
	DataCOC []SertifikatPerbalai
}

type SertifikatPerbalai struct {
	Setifikat string `json:"setifikat"`
	Lembaga   []Lembaga
}

type Lembaga struct {
	NamaLembaga string `json:"nama_embaga"`
	Total       int    `json:"total"`
}

type LembagaResponse struct {
	Lembaga    string       `json:"Lembaga"`
	Sertifikat []DiklatStat `json:"sertifikat"`
}

type UnitKerjaResponse struct {
	UnitKerja  string            `json:"UnitKerja"`
	Sertifikat []SertifikasiStat `json:"sertifikat"`
}

type DiklatStat struct {
	NamaDiklat string `json:"Nama Diklat"`
	Total      int    `json:"total"`
}

type SertifikasiStat struct {
	JenisSertifikasi string `json:"Jenis Sertifikasi"`
	Total            int    `json:"total"`
}
type APIResponse struct {
	Lembaga    string `json:"Lembaga"`
	Sertifikat []struct {
		NamaDiklat string `json:"Nama Diklat"`
		Total      int    `json:"total"`
	} `json:"sertifikat"`
}

func GetDataBalaiSertifikats(c *fiber.Ctx) error {

	var resultsLembaga []struct {
		Lembaga     string
		JenisDiklat string
		Jumlah      int
	}

	// Ambil parameter waktu dari query string
	waktuMulai := c.Query("waktu_awal")
	waktuBerakhir := c.Query("waktu_berakhir")

	// Pastikan parameter waktu tidak kosong, jika kosong gunakan default
	if waktuMulai == "" {
		waktuMulai = "2024-06-01"
	}
	if waktuBerakhir == "" {
		waktuBerakhir = "2024-12-31"
	}

	err := database.DB1.Table("sertifikat s").
		Select("pl.pl_nama_lembaga as lembaga, d.d_sub_jenis_pendidikan as jenis_diklat, COUNT(*) as jumlah").
		Joins("JOIN master_diklat d ON s.d_id = d.d_id").
		Joins("JOIN master_lembaga ml ON d.l_id = ml.l_id").
		Joins("JOIN master_profil_lembaga pl ON ml.pl_id = pl.pl_id").
		Where("s.isprint = ? AND s.created_on BETWEEN ? AND ?", 1, waktuMulai, waktuBerakhir).
		Group("pl.pl_nama_lembaga, d.d_sub_jenis_pendidikan").
		Order("pl.pl_nama_lembaga ASC").
		Scan(&resultsLembaga).Error

	total := 0
	for _, jumlah := range resultsLembaga {
		total += jumlah.Jumlah // Menggunakan += tanpa ++
	}

	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	// Query Kedua (Unit Kerja)
	var resultsUnitKerja []struct {
		UnitKerja        string
		JenisSertifikasi string
		Jumlah           int
	}

	err = database.DB1.Table("sertifikat s").
		Select("ml.uk_nama as unit_kerja, d.ru_jenis_setifikasi as jenis_sertifikasi, COUNT(*) as jumlah").
		Joins("JOIN rencana_ujian d ON s.d_id = d.ru_id").
		Joins("JOIN master_unit_kerja ml ON d.ru_unit_kerja = ml.uk_id").
		Where("s.isprint = ? AND s.created_on BETWEEN ? AND ?", 1, waktuMulai, waktuBerakhir).
		Group("ml.uk_nama, d.ru_jenis_setifikasi").
		Order("ml.uk_nama ASC").
		Scan(&resultsUnitKerja).Error

	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	// Process Lembaga Data
	lembagaMap := make(map[string]LembagaResponse)
	for _, item := range resultsLembaga {
		if _, ok := lembagaMap[item.Lembaga]; !ok {
			lembagaMap[item.Lembaga] = LembagaResponse{
				Lembaga: item.Lembaga,
			}
		}
		entry := lembagaMap[item.Lembaga]
		entry.Sertifikat = append(entry.Sertifikat, DiklatStat{
			NamaDiklat: item.JenisDiklat,
			Total:      item.Jumlah,
		})
		lembagaMap[item.Lembaga] = entry
	}

	// Process Unit Kerja Data
	unitKerjaMap := make(map[string]UnitKerjaResponse)
	for _, item := range resultsUnitKerja {
		if _, ok := unitKerjaMap[item.UnitKerja]; !ok {
			unitKerjaMap[item.UnitKerja] = UnitKerjaResponse{
				UnitKerja: item.UnitKerja,
			}
		}
		entry := unitKerjaMap[item.UnitKerja]
		entry.Sertifikat = append(entry.Sertifikat, SertifikasiStat{
			JenisSertifikasi: item.JenisSertifikasi,
			Total:            item.Jumlah,
		})
		unitKerjaMap[item.UnitKerja] = entry
	}

	// Convert maps to slices
	finalLembaga := make([]LembagaResponse, 0, len(lembagaMap))
	for _, v := range lembagaMap {
		finalLembaga = append(finalLembaga, v)
	}

	finalUnitKerja := make([]UnitKerjaResponse, 0, len(unitKerjaMap))
	for _, v := range unitKerjaMap {
		finalUnitKerja = append(finalUnitKerja, v)
	}

	respon := CombinedResponse{
		DataLembaga:   finalLembaga,
		DataUnitKerja: finalUnitKerja,
	}

	return c.JSON(fiber.Map{
		"Pesan": "Berhasil Get Data",
		"data":  respon,
	})
}

func GetDataBalaiSertifikatsReverse(c *fiber.Ctx) error {

	var resultsLembaga []struct {
		JenisDiklat string
		Lembaga     string
		Jumlah      int
	}

	// Ambil parameter waktu dari query string
	waktuMulai := c.Query("waktu_awal")
	waktuBerakhir := c.Query("waktu_berakhir")

	// Pastikan parameter waktu tidak kosong, jika kosong gunakan default
	if waktuMulai == "" {
		waktuMulai = "2024-06-01"
	}
	if waktuBerakhir == "" {
		waktuBerakhir = "2024-12-31"
	}

	err := database.DB1.Table("sertifikat s").
		Select("d.d_sub_jenis_pendidikan AS jenis_diklat, pl.pl_nama_lembaga AS lembaga, COUNT(*) AS jumlah").
		Joins("JOIN master_diklat d ON s.d_id = d.d_id").
		Joins("JOIN master_lembaga ml ON d.l_id = ml.l_id").
		Joins("JOIN master_profil_lembaga pl ON ml.pl_id = pl.pl_id").
		Where("s.isprint = ? AND s.created_on BETWEEN ? AND ?", 1, waktuMulai, waktuBerakhir).
		Group("d.d_sub_jenis_pendidikan, pl.pl_nama_lembaga").
		Order("d.d_sub_jenis_pendidikan ASC").
		Scan(&resultsLembaga).Error

	total := 0
	for _, jumlah := range resultsLembaga {
		total += jumlah.Jumlah // Menggunakan += tanpa ++
	}

	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	// Query Kedua (Unit Kerja)
	var resultsUnitKerja []struct {
		UnitKerja        string
		JenisSertifikasi string
		Jumlah           int
	}

	err = database.DB1.Table("sertifikat s").
		Select("d.ru_jenis_setifikasi AS jenis_sertifikasi, ml.uk_nama AS unit_kerja, COUNT(*) AS jumlah").
		Joins("JOIN rencana_ujian d ON s.d_id = d.ru_id").
		Joins("JOIN master_unit_kerja ml ON d.ru_unit_kerja = ml.uk_id").
		Where("s.isprint = ? AND s.created_on BETWEEN ? AND ?", 1, waktuMulai, waktuBerakhir).
		Group("d.ru_jenis_setifikasi, ml.uk_nama").
		Order("d.ru_jenis_setifikasi ASC").
		Scan(&resultsUnitKerja).Error

	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	// Inisialisasi struct untuk response
	var response RespostSertifikat

	// Memproses resultsLembaga ke dalam DataCOC
	lembagaMap := make(map[string][]Lembaga)
	for _, r := range resultsLembaga {
		lembagaMap[r.JenisDiklat] = append(lembagaMap[r.JenisDiklat], Lembaga{
			NamaLembaga: r.Lembaga,
			Total:       r.Jumlah,
		})
	}

	for jenisDiklat, lembagaList := range lembagaMap {
		response.DataCOC = append(response.DataCOC, SertifikatPerbalai{
			Setifikat: jenisDiklat,
			Lembaga:   lembagaList,
		})
	}

	// Memproses resultsUnitKerja ke dalam DataCOP
	unitKerjaMap := make(map[string][]Lembaga)
	for _, r := range resultsUnitKerja {
		unitKerjaMap[r.JenisSertifikasi] = append(unitKerjaMap[r.JenisSertifikasi], Lembaga{
			NamaLembaga: r.UnitKerja,
			Total:       r.Jumlah,
		})
	}

	for jenisSertifikasi, unitKerjaList := range unitKerjaMap {
		response.DataCOP = append(response.DataCOP, SertifikatPerbalai{
			Setifikat: jenisSertifikasi,
			Lembaga:   unitKerjaList,
		})
	}

	// response sekarang berisi hasil yang telah dikelompokkan ke dalam struct yang sesuai

	return c.JSON(fiber.Map{
		"Pesan": "Berhasil Get Data",
		"data":  response,
	})
}
