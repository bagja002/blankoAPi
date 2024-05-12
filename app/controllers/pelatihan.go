package controllers

import (
	"fmt"
	"strings"
	"template/app/entity"
	"template/pkg/config"
	"template/pkg/database"
	"template/pkg/tools"

	"github.com/gofiber/fiber/v2"
)

type Pelatihan struct {
	IdPelatihan              uint   `gorm:"primary_key;auto_increment" json:"id_pelatihan"`
	IdLemdik                 string   `json:"IdLemdik"`
	KodePelatihan            string `json:"KodePelatihan"`
	NamaPelatihan            string `json:"NamaPelatihan"`
	PenyelenggaraPelatihan   string `json:"PenyelenggaraPelatihan"`
	DetailPelatihan          string `json:"DetailPelatihan"`
	JenisPelatihan           string `json:"JenisPelatihan"`
	BidangPelatihan          string `json:"BidangPelatihan"`
	DukunganProgramTerobosan string `json:"DukunganProgramTerobosan"`
	TanggalMulaiPelatihan    string `json:"TanggalMulaiPelatihan"`
	TanggalBerakhirPelatihan string `json:"TanggalBerakhirPelatihan"`
	HargaPelatihan           string    `json:"HargaPelatihan"`
	Instruktur               string `json:"instruktur"`
	FotoPelatihan string	
	Status                   string `json:"status"`
	MemoPusat                string `json:"memo_pusat"`
	SilabusPelatihan         string `json:"silabus_pelatihan"`
	LokasiPelatihan          string `json:"lokasi_pelatihan"`
	PelaksanaanPelatihan     string `json:"pelaksanaan_pelatihan"`
	UjiKompotensi            string `json:"uji_kompetensi"`
	KoutaPelatihan           string `json:"kouta_pelatihan"`
	AsalPelatihan            string `json:"asal_pelatihan"`
	AsalSertifikat           string `json:"asal_sertifikat"`
	JenisSertifikat          string `json:"jenis_sertifikat"`
	TtdSertifikat            string `json:"ttd_sertifikat"`
	NoSertifikat             string `json:"no_sertifikat"`
	IdSaranaPrasarana        string `json:"id_sarana_prasarana"`
	IdKonsumsi               string `json:"id_konsumsi"`
	CreateAt                 string `json:"created_at"`
	UpdateAt                 string `json:"updated_at"`
}

func CreatePelatihan(c *fiber.Ctx) error {

	//Pakai JWT

	//Foto terlebih dahulu


	file, err := c.FormFile("photo_pelatihan")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"Message": "Failed to retrieve file", "Error": err.Error()})
	}
	// Simpan file ke dalam direktori static/merchant
	if err := c.SaveFile(file, "public/static/pelatihan/"+strings.ReplaceAll(file.Filename, " ", "")); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"Message": "Failed to save file", "Error": err.Error()})
	}

	//Inputan Biasa
	var request Pelatihan
	if err := c.BodyParser(&request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"Message": "Failed to parse request body", "Error": err.Error()})
	}
	
	fmt.Println(request)

	newPelatihan := entity.Pelatihan{
		IdLemdik:                 uint(tools.StringToInt(request.IdLemdik)),
		KodePelatihan:            request.KodePelatihan,
		NamaPelatihan:            request.NamaPelatihan,
		PenyelenggaraPelatihan:   request.PenyelenggaraPelatihan,
		DetailPelatihan:          request.DetailPelatihan,
		FotoPelatihan: strings.ReplaceAll(file.Filename, " ", ""),
		JenisPelatihan:           request.JenisPelatihan,
		BidangPelatihan:          request.BidangPelatihan,
		DukunganProgramTerobosan: request.DukunganProgramTerobosan,
		TanggalMulaiPelatihan:    request.TanggalMulaiPelatihan,
		TanggalBerakhirPelatihan: request.TanggalBerakhirPelatihan,
		HargaPelatihan:           tools.StringToInt(request.HargaPelatihan),
		Instruktur:               request.Instruktur,
		Status:                   request.Status,
		MemoPusat:                request.MemoPusat,
		SilabusPelatihan:         request.SilabusPelatihan,
		LokasiPelatihan:          request.LokasiPelatihan,
		PelaksanaanPelatihan:     request.PelaksanaanPelatihan,
		UjiKompotensi:            request.UjiKompotensi,
		KoutaPelatihan:           request.KoutaPelatihan,
		AsalPelatihan:            request.AsalPelatihan,
		AsalSertifikat:           request.AsalSertifikat,
		JenisSertifikat:          request.JenisSertifikat,
		TtdSertifikat:            request.TtdSertifikat,
		NoSertifikat:             request.NoSertifikat,
		IdSaranaPrasarana:        request.IdSaranaPrasarana,
		IdKonsumsi:               request.IdKonsumsi,
		CreateAt:                 tools.TimeNowJakarta(),
		UpdateAt:                 tools.TimeNowJakarta(),
	}

	fmt.Println(request.IdLemdik)

	tx := database.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Create(&newPelatihan).Error; err != nil {
		tx.Rollback()
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"Message": "Failed to create merchant", "Error": err.Error()})
	}

	// Commit transaksi jika semuanya berhasil
	if err := tx.Commit().Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"Message": "Failed to commit transaction", "Error": err.Error()})
	}

	return c.JSON(fiber.Map{"Message": "Successfully Add Pelatihan"})

}

func GetPelatihan(c *fiber.Ctx) error {

	viper := config.NewViper()
	baseUrl := viper.GetString("web.baseUrl")

	//ambil By Id
	id:= c.Query("id")
	if id != "" {	

		var pelatihan entity.Pelatihan


		database.DB.Where("id_Pelatihan = ?", id).Find(&pelatihan)

		pelatihan.FotoPelatihan = baseUrl + "/public/static/pelatihan/" + pelatihan.FotoPelatihan
		//jangan lupa filter

		return c.JSON(fiber.Map{
			"Pesan":"",
		})
	}
	//Bikin filter berdasarkan
	//Bidang Pelatihan
	//Penyelenggara Pelatihan
	//Pelaksanaan Pelatihan

	var pelatihan []entity.Pelatihan


	database.DB.Find(&pelatihan)

	for i, _:= range pelatihan {
		pelatihan[i].FotoPelatihan = baseUrl+"http://127.0.0.1:3000/public/static/pelatihan/" + pelatihan[i].FotoPelatihan
	}


	return c.JSON(fiber.Map{
		"Pesan": "Sukses Mengambil Data",
		"data":  pelatihan,
	})
}

func UpdatePelatihan(c *fiber.Ctx) error {

	//id:= c.Query("id")

	return c.JSON(fiber.Map{
		"Pesan": "Sukses Update Pelatihan",
		"data":  "",
	})
}

func DeletePelatihan(c *fiber.Ctx) error {

	return c.JSON(fiber.Map{
		"Pesan": "Sukses Hapus pelatihan",
	})
}

func SearchPelatihan(c *fiber.Ctx) error {

	//Memakai Query Like

	return c.JSON(fiber.Map{
		"Pesan": "Berhasil Mencari Pelatihan ",
		"data":  "",
	})
}
