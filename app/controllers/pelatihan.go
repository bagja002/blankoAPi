package controllers

import (
	"log"

	"strings"
	"template/app/entity"
	"template/app/models"
	"template/pkg/config"
	"template/pkg/database"
	"template/pkg/tools"

	"github.com/gofiber/fiber/v2"
)

func TestPreloadPencapaian(c *fiber.Ctx) error {

	var Pelatihan entity.Pelatihan

	id := c.Query("id")

	database.DB.Preload("SarprasPelatihan").Where("id_pelatihan = ?", id).Find(&Pelatihan)

	return c.JSON(fiber.Map{
		"Pesan": "Sukses",
		"data":  Pelatihan,
	})
}

func CreatePelatihan(c *fiber.Ctx) error {

	//Pake Role Super admin/ admin pusat
	id_admin, _ := c.Locals("id_admin").(int)
	role, _ := c.Locals("role").(string)
	names, _ := c.Locals("name").(string)

	tools.ValidationJwtLemdik(c, role, id_admin, names)

	//Foto terlebih dahulu

	file, err := c.FormFile("photo_pelatihan")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"Message": "Failed to retrieve file", "Error": err.Error()})
	}

	//Inputan Biasa
	var request models.Pelatihan
	if err := c.BodyParser(&request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"Message": "Failed to parse request body", "Error": err.Error()})
	}

	newPelatihan := entity.Pelatihan{
		IdLemdik:                 uint(id_admin),
		KodePelatihan:            request.KodePelatihan,
		NamaPelatihan:            request.NamaPelatihan,
		PenyelenggaraPelatihan:   request.PenyelenggaraPelatihan,
		DetailPelatihan:          request.DetailPelatihan,
		FotoPelatihan:            strings.ReplaceAll(file.Filename, " ", ""),
		JenisPelatihan:           request.JenisPelatihan,
		BidangPelatihan:          request.BidangPelatihan,
		DukunganProgramTerobosan: request.DukunganProgramTerobosan,
		TanggalMulaiPelatihan:    request.TanggalMulaiPelatihan,
		TanggalBerakhirPelatihan: request.TanggalBerakhirPelatihan,
		HargaPelatihan:           tools.StringToInt(request.HargaPelatihan),
		Instruktur:               request.Instruktur,
		Status:                   "Belum Publish",
		MemoPusat:                request.MemoPusat,
		SilabusPelatihan:         request.SilabusPelatihan,
		LokasiPelatihan:          request.LokasiPelatihan,
		PelaksanaanPelatihan:     request.PelaksanaanPelatihan,

		//Ketika Uji Ada Uji kom
		UjiKompotensi:   request.UjiKompotensi,
		KoutaPelatihan:  request.KoutaPelatihan,
		AsalPelatihan:   request.AsalPelatihan,
		AsalSertifikat:  request.AsalSertifikat,
		JenisSertifikat: request.JenisSertifikat,
		TtdSertifikat:   request.TtdSertifikat,
		NoSertifikat:    request.NoSertifikat,

		IdKonsumsi: request.IdKonsumsi,
		CreateAt:   tools.TimeNowJakarta(),
		UpdateAt:   tools.TimeNowJakarta(),
	}

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

	//CEK LALO HASILNYA ITU TRUE
	//tambahkan pelatihan

	//tampabih sarpras

	id_sarpras := request.IdSaranaPrasarana

	list_id_sarpras := strings.Split(id_sarpras, ",")

	if id_sarpras != "" {
		for _, lis_id := range list_id_sarpras {
			newSarprasPelatihan := entity.SarprasPelatihan{
				IdPelatihan: newPelatihan.IdPelatihan,
				IdLemdik:    uint(id_admin),
				IdSarpras:   uint(tools.StringToInt(lis_id)),
			}

			if err := database.DB.Create(&newSarprasPelatihan).Error; err != nil {
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"Message": "Gagal Nambahin prass", "Error": err.Error()})
			}
		}
	}

	if newPelatihan.UjiKompotensi == "true" {
		//Kirim data pelatihan ke table sertiifikasi

	}

	//Menambahkan Masukan materi

	// Simpan file ke dalam direktori static/merchant
	if err := c.SaveFile(file, "public/static/pelatihan/"+strings.ReplaceAll(file.Filename, " ", "")); err != nil {
		log.Println("Berubah")
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"Message": "Failed to save file", "Error": err.Error()})
	}

	return c.JSON(fiber.Map{"Message": "Successfully Add Pelatihan"})

}

func GetPelatihan(c *fiber.Ctx) error {

	viper := config.NewViper()
	baseUrl := viper.GetString("web.baseUrl")

	//ambil By Id

	id := c.Query("id")
	bidangPelatihan := c.Query("bidang_pelatihan")
	penyelenggaraPelatihan := c.Query("penyelenggara_pelatihan")
	idLemdik := c.Query("id_lemdik")

	var pelatihan []entity.Pelatihan

	queryBase := database.DB

	if id != "" {
		queryBase = queryBase.Where("id_pelatihan = ?", id)
	}
	if bidangPelatihan != "" {
		queryBase = queryBase.Where("bidang_pelatihan = ?", bidangPelatihan)
	}
	if penyelenggaraPelatihan != "" {
		queryBase = queryBase.Where("penyelenggara_pelatihan = ?", penyelenggaraPelatihan)
	}
	if idLemdik != "" {
		queryBase = queryBase.Where("id_lemdik = ? ", idLemdik)
	}

	queryBase.Find(&pelatihan)

	for i, _ := range pelatihan {
		pelatihan[i].FotoPelatihan = baseUrl + "/public/static/pelatihan/" + pelatihan[i].FotoPelatihan
	}

	return c.JSON(fiber.Map{
		"Pesan": "Sukses Mengambil Data",
		"data":  pelatihan,
	})
}

func UpdatePelatihan(c *fiber.Ctx) error {

	id_admin, _ := c.Locals("id_admin").(int)
	role, _ := c.Locals("role").(string)
	names, _ := c.Locals("name").(string)

	tools.ValidationJwtLemdik(c, role, id_admin, names)

	id := c.Query("id")

	SilabusPelatihan, _ := c.FormFile("SilabusPelatihan")
	ModuleMateri, _ := c.FormFile("ModuleMateri")
	SuratPemberitahuan, _ := c.FormFile("SuratPemberitahuan")
	BeritaAcara, _ := c.FormFile("BeritaAcara")

	var pelatihan entity.Pelatihan

	database.DB.Where("id_pelatihan = ?", id).Find(&pelatihan)

	// Menginisialisasi koneksi database
	tx := database.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if SilabusPelatihan != nil || ModuleMateri != nil || SuratPemberitahuan != nil || BeritaAcara != nil {
		if SilabusPelatihan != nil {
			pelatihan.SilabusPelatihan = SilabusPelatihan.Filename
			if err := c.SaveFile(SilabusPelatihan, "public/silabus/pelatihan/"+SilabusPelatihan.Filename); err != nil {
				tx.Rollback()
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
					"pesan": "Gagal menyimpan file EvaluasiRenaksi",
					"error": err.Error(),
				})
			}
		}

		if ModuleMateri != nil {
			pelatihan.ModuleMateri = ModuleMateri.Filename
			if err := c.SaveFile(ModuleMateri, "public/silabus/pelatihan/"+ModuleMateri.Filename); err != nil {
				tx.Rollback()
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
					"pesan": "Gagal menyimpan file Module Materi",
					"error": err.Error(),
				})
			}
		}

		if SuratPemberitahuan != nil {
			pelatihan.SuratPemberitahuan = SuratPemberitahuan.Filename
			if err := c.SaveFile(SuratPemberitahuan, "public/static/suratPemberitahuan/"+SuratPemberitahuan.Filename); err != nil {
				tx.Rollback()
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
					"pesan": "Gagal menyimpan file Surat Pemberitahuan",
					"error": err.Error(),
				})
			}
		}

		if BeritaAcara != nil {
			pelatihan.BeritaAcara = BeritaAcara.Filename
			if err := c.SaveFile(BeritaAcara, "public/static/BeritaAcara/"+BeritaAcara.Filename); err != nil {
				tx.Rollback()
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
					"pesan": "Gagal menyimpan file Surat Pemberitahuan",
					"error": err.Error(),
				})
			}
		}

		if err := tx.Model(&pelatihan).Where("id_pelatihan = ?", id).Updates(&pelatihan).Error; err != nil {
			tx.Rollback()
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"pesan": "Gagal memperbarui MonitoringEvaluasi",
				"error": err.Error(),
			})
		}

	}

	//Yang biasanya
	var request entity.Pelatihan
	if err := c.BodyParser(&request); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"pesan": "gagal reques",
		})
	}

	updates := entity.Pelatihan{
		NamaPelatihan:            request.NamaPelatihan,
		PenyelenggaraPelatihan:   request.PenyelenggaraPelatihan,
		DetailPelatihan:          request.DetailPelatihan,
		JenisPelatihan:           request.JenisPelatihan,
		BidangPelatihan:          request.BidangPelatihan,
		DukunganProgramTerobosan: request.DukunganProgramTerobosan,
		TanggalMulaiPelatihan:    request.TanggalMulaiPelatihan,
		TanggalBerakhirPelatihan: request.TanggalBerakhirPelatihan,
		HargaPelatihan:           request.HargaPelatihan,
		Instruktur:               request.Instruktur,
		Status:                   request.Status,
		MemoPusat:                request.MemoPusat,
		SilabusPelatihan:         request.SilabusPelatihan,
		LokasiPelatihan:          request.LokasiPelatihan,
		PelaksanaanPelatihan:     request.PelaksanaanPelatihan,
		UjiKompotensi:            request.UjiKompotensi,
		KoutaPelatihan:           request.KoutaPelatihan,
		AsalPelatihan:            request.AsalPelatihan,
		AsalSertifikat:           pelatihan.AsalSertifikat,
		JenisSertifikat:          request.JenisPelatihan,
		TtdSertifikat:            request.TtdSertifikat,
		NoSertifikat:             request.NoSertifikat,

		StatusApproval: request.Status,

		UpdateAt: tools.TimeNowJakarta(),

		PemberitahuanDiterima: request.PemberitahuanDiterima,

		CatatanPemberitahuanByPusat:  request.CatatanPemberitahuanByPusat,
		PenerbitanSertifikatDiterima: request.PenerbitanSertifikatDiterima,
		CatatanPenerbitanByPusat:     request.CatatanPemberitahuanByPusat,
	}

	if err := tx.Model(&pelatihan).Where("id_pelatihan = ?", id).Updates(&updates).Error; err != nil {
		tx.Rollback()
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"pesan": "Gagal memperbarui MonitoringEvaluasi",
			"error": err.Error(),
		})
	}

	if err := tx.Commit().Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"pesan": "Gagal melakukan commit transaksi",
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"Pesan": "Sukses Update Pelatihan",
		"data":  pelatihan,
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
