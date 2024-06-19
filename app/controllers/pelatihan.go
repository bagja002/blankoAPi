package controllers

import (
	"log"
	"os"

	"strings"
	"template/app/entity"
	"template/app/models"
	"template/pkg/config"
	"template/pkg/database"
	"template/pkg/generator"
	"template/pkg/tools"

	"github.com/gofiber/fiber/v2"
)

// Materi Pelatihan

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
		//SilabusPelatihan:         request.SilabusPelatihan,  nanti jadikan ini sebagai file jadikan dia wajib
		LokasiPelatihan:      request.LokasiPelatihan,
		PelaksanaanPelatihan: request.PelaksanaanPelatihan,

		//Ketika Uji Ada Uji kom
		UjiKompotensi:   request.UjiKompotensi, //isinya true or false
		KoutaPelatihan:  request.KoutaPelatihan,
		AsalPelatihan:   request.AsalPelatihan, //Asal Peserta pelatihan
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

	//tambahin sarpras

	//Inputan sarpras ini adalah stirng yang nanti di pecah menjadi list

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

	//Menambahkan Data Pelatihan Materi

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

	queryBase.Preload("UserPelatihan").
		Preload("SarprasPelatihan").
		Preload("MateriPelatihan").
		Find(&pelatihan)

	for i, _ := range pelatihan {
		pelatihan[i].FotoPelatihan = baseUrl + "/public/static/pelatihan/" + pelatihan[i].FotoPelatihan
	}

	//File File yang harus di Tambahkan

	//Surat Pemberitahuan
	for z, _ := range pelatihan {
		pelatihan[z].SuratPemberitahuan = baseUrl + "/public/static/suratPemberitahuan/" + pelatihan[z].SuratPemberitahuan
	}

	//BeritaAcara
	for x, _ := range pelatihan {
		pelatihan[x].BeritaAcara = baseUrl + "/public/static/BeritaAcara/" + pelatihan[x].BeritaAcara
	}

	//Silabus Pelatihan
	for y, _ := range pelatihan {
		pelatihan[y].SilabusPelatihan = baseUrl + "/public/silabus/pelatihan/" + pelatihan[y].SilabusPelatihan
	}
	//Module Materi Pelatihan
	for p, _ := range pelatihan {
		pelatihan[p].ModuleMateri = baseUrl + "/public/module/pelatihan/" + pelatihan[p].SilabusPelatihan
	}
	//Memo
	for p, _ := range pelatihan {
		pelatihan[p].MemoPusat = baseUrl + "/public/static/memo/" + pelatihan[p].MemoPusat
	}

	return c.JSON(fiber.Map{
		"Pesan": "Sukses Mengambil Data",
		"data":  pelatihan,
	})
}

// UpdatePelatihan updates the pelatihan by ID
func UpdatePelatihan(c *fiber.Ctx) error {
	idAdmin, ok := c.Locals("id_admin").(int)
	role, okRole := c.Locals("role").(string)
	names, okName := c.Locals("name").(string)
	if !ok || !okRole || !okName {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid or missing id_admin, role, or name",
		})
	}

	tools.ValidationJwtLemdik(c, role, idAdmin, names)

	id := c.Query("id")

	SilabusPelatihan, _ := c.FormFile("SilabusPelatihan")
	ModuleMateri, _ := c.FormFile("ModuleMateri")
	SuratPemberitahuan, _ := c.FormFile("SuratPemberitahuan")
	BeritaAcara, _ := c.FormFile("BeritaAcara")
	photoPelatihan, _ := c.FormFile("photo_pelatihan")

	var pelatihan entity.Pelatihan
	if err := database.DB.Where("id_pelatihan = ?", id).First(&pelatihan).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "Pelatihan not found",
		})
	}

	tx := database.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Handle file uploads and updates
	if SilabusPelatihan != nil {
		newPath := "public/silabus/pelatihan/" + SilabusPelatihan.Filename
		if err := c.SaveFile(SilabusPelatihan, newPath); err != nil {
			tx.Rollback()
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"message": "Failed to save Silabus Pelatihan",
				"error":   err.Error(),
			})
		}
		oldPath := "public/silabus/pelatihan/" + pelatihan.SilabusPelatihan
		if pelatihan.SilabusPelatihan != "" {
			os.Remove(oldPath)
		}
		pelatihan.SilabusPelatihan = SilabusPelatihan.Filename
	}

	if photoPelatihan != nil {
		newPath := "public/static/pelatihan/" + photoPelatihan.Filename
		if err := c.SaveFile(photoPelatihan, newPath); err != nil {
			tx.Rollback()
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"message": "Failed to save photo Pelatihan",
				"error":   err.Error(),
			})
		}
		oldPath := "public/static/pelatihan/" + pelatihan.FotoPelatihan
		if pelatihan.FotoPelatihan != "" {
			os.Remove(oldPath)
		}
		pelatihan.FotoPelatihan = photoPelatihan.Filename
	}

	if ModuleMateri != nil {
		newPath := "public/silabus/pelatihan/" + ModuleMateri.Filename
		if err := c.SaveFile(ModuleMateri, newPath); err != nil {
			tx.Rollback()
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"message": "Failed to save Module Materi",
				"error":   err.Error(),
			})
		}
		oldPath := "public/silabus/pelatihan/" + pelatihan.ModuleMateri
		if pelatihan.ModuleMateri != "" {
			os.Remove(oldPath)
		}
		pelatihan.ModuleMateri = ModuleMateri.Filename
	}

	if SuratPemberitahuan != nil {
		newPath := "public/static/suratPemberitahuan/" + SuratPemberitahuan.Filename
		if err := c.SaveFile(SuratPemberitahuan, newPath); err != nil {
			tx.Rollback()
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"message": "Failed to save Surat Pemberitahuan",
				"error":   err.Error(),
			})
		}
		oldPath := "public/static/suratPemberitahuan/" + pelatihan.SuratPemberitahuan
		if pelatihan.SuratPemberitahuan != "" {
			os.Remove(oldPath)
		}
		pelatihan.SuratPemberitahuan = SuratPemberitahuan.Filename
	}

	if BeritaAcara != nil {
		newPath := "public/static/BeritaAcara/" + BeritaAcara.Filename
		if err := c.SaveFile(BeritaAcara, newPath); err != nil {
			tx.Rollback()
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"message": "Failed to save Berita Acara",
				"error":   err.Error(),
			})
		}
		oldPath := "public/static/BeritaAcara/" + pelatihan.BeritaAcara
		if pelatihan.BeritaAcara != "" {
			os.Remove(oldPath)
		}
		pelatihan.BeritaAcara = BeritaAcara.Filename
	}

	// Update pelatihan fields
	var request entity.Pelatihan
	if err := c.BodyParser(&request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Failed to parse request body",
			"error":   err.Error(),
		})
	}

	pelatihan.NamaPelatihan = request.NamaPelatihan
	pelatihan.PenyelenggaraPelatihan = request.PenyelenggaraPelatihan
	pelatihan.DetailPelatihan = request.DetailPelatihan
	pelatihan.JenisPelatihan = request.JenisPelatihan
	pelatihan.BidangPelatihan = request.BidangPelatihan
	pelatihan.DukunganProgramTerobosan = request.DukunganProgramTerobosan
	pelatihan.TanggalMulaiPelatihan = request.TanggalMulaiPelatihan
	pelatihan.TanggalBerakhirPelatihan = request.TanggalBerakhirPelatihan
	pelatihan.HargaPelatihan = request.HargaPelatihan
	pelatihan.Instruktur = request.Instruktur
	pelatihan.Status = request.Status
	pelatihan.MemoPusat = request.MemoPusat
	pelatihan.LokasiPelatihan = request.LokasiPelatihan
	pelatihan.PelaksanaanPelatihan = request.PelaksanaanPelatihan
	pelatihan.UjiKompotensi = request.UjiKompotensi
	pelatihan.KoutaPelatihan = request.KoutaPelatihan
	pelatihan.AsalPelatihan = request.AsalPelatihan
	pelatihan.AsalSertifikat = request.AsalSertifikat
	pelatihan.JenisSertifikat = request.JenisPelatihan
	pelatihan.TtdSertifikat = request.TtdSertifikat
	pelatihan.NoSertifikat = request.NoSertifikat
	pelatihan.StatusApproval = request.StatusApproval
	pelatihan.UpdateAt = tools.TimeNowJakarta()
	pelatihan.PemberitahuanDiterima = request.PemberitahuanDiterima
	pelatihan.CatatanPemberitahuanByPusat = request.CatatanPemberitahuanByPusat
	pelatihan.PenerbitanSertifikatDiterima = request.PenerbitanSertifikatDiterima
	pelatihan.CatatanPenerbitanByPusat = request.CatatanPemberitahuanByPusat

	if err := tx.Model(&pelatihan).Where("id_pelatihan = ?", id).Updates(&pelatihan).Error; err != nil {
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
		"message": "Pelatihan updated successfully",
		"data":    pelatihan,
	})
}

// Publish Sertifikat Pake ID Pelatihan
func PublishSertifikat(c *fiber.Ctx) error {

	id_admin, _ := c.Locals("id_admin").(int)
	role, _ := c.Locals("role").(string)
	names, _ := c.Locals("name").(string)

	tools.ValidationJwtLemdik(c, role, id_admin, names)

	id := c.Query("id")

	if id == "" {
		return c.Status(400).JSON(fiber.Map{
			"Pesan": "Masukan Parameter ID Pelatihan",
		})

	}

	var pelatihan entity.Pelatihan

	database.DB.Where("id_pelatihan = ?", id).Find(&pelatihan)

	if pelatihan.NoSertifikat != "" {
		return c.Status(400).JSON(fiber.Map{
			"Pesan": "Sertifikat Sudah ada",
			"data":  pelatihan.NoSertifikat,
		})
	}
	//BeritaAcara, _ := c.FormFile("BeritaAcara")
	//c.SaveFile(BeritaAcara, "public/static/BeritaAcara/"+BeritaAcara.Filename)

	//Generate Sertifikat per balai
	var request entity.Pelatihan
	if err := c.BodyParser(&request); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"pesan": "gagal reques",
		})
	}

	NewNoSertif := generator.GenerateSertifikat(tools.IntToString(id_admin), id, c)

	updates := entity.Pelatihan{
		//BeritaAcara:   BeritaAcara.Filename,
		TtdSertifikat: request.TtdSertifikat,
		NoSertifikat:  NewNoSertif,
	}

	database.DB.Model(&pelatihan).Updates(&updates)

	return c.JSON(fiber.Map{
		"Pesan": "Sukses Generate Sertifikat",
		"Data":  pelatihan,
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
