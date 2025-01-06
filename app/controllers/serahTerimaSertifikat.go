package controllers

import (
	"os"
	"template/app/entity"
	"template/pkg/database"
	"template/pkg/tools"

	"github.com/gofiber/fiber/v2"
)

/*
"public/static/ttd-penerima",

	"public/static/ttd-pemberi",
	"public/static/bukti-serah-terima",

	"public/static/ttd-pengiriman",
	"public/static/bukti-resi",
	"public/static/bukti-pengiriman-sertifikat",
	"public/static/bukti-penerimaan-sertifikat",
*/
func CreateSerahterimaSertifikat(c *fiber.Ctx) error {

	var request entity.SerahTerimaSertifikat

	if err := c.BodyParser(&request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"Message": "Failed to parse request body", "Error": err.Error()})
	}

	//File File Duluu
	ttd_penerima, _ := c.FormFile("ttd_penerima")

	ttd_pemberi, _ := c.FormFile("ttd_pemberi")

	bukti_serah_terima, _ := c.FormFile("bukti_serah_terima")

	newData := entity.SerahTerimaSertifikat{
		NamaPenerima:          request.NamaPenerima,
		Jabatan:               request.Jabatan,
		Instansi:              request.Instansi,
		NamaKegiatan:          request.NamaKegiatan,
		TanggalPengambilan:    request.TanggalPengambilan,
		NoSeriBlanko:          request.NoSeriBlanko,
		JenisSertifikat:       request.JenisSertifikat,
		TandaTanganPenerima:   ttd_penerima.Filename, // File (string atau path file)
		TandaTanganPemberi:    ttd_pemberi.Filename,  // File (string atau path file)
		NamaPemberiSertifikat: request.NamaPemberiSertifikat,
		BuktiSerahTerima:      bukti_serah_terima.Filename, // File (string atau path file)
		CreateAt:              tools.TimeNowJakarta(),      // Atur waktu sekarang untuk CreateAt
		UpdateAt:              tools.TimeNowJakarta(),      // Atur waktu sekarang untuk UpdateAt
		Status:                "Active",                    //
	}

	//tambahkan Data data
	database.DB.Create(&newData)

	// Simpan data Ke dalam Database
	if bukti_serah_terima != nil {
		if err := c.SaveFile(ttd_pemberi, "public/static/bukti-serah-terima/"+tools.RemoverSpaci(ttd_pemberi.Filename)); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Gagal Menyimpan File "})
		}
	}

	if ttd_pemberi != nil {
		if err := c.SaveFile(ttd_pemberi, "public/static/ttd-pemberi/"+tools.RemoverSpaci(ttd_pemberi.Filename)); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Gagal Menyimpan File "})
		}
	}
	if ttd_penerima != nil {
		if err := c.SaveFile(ttd_penerima, "public/static/ttd-penerima/"+tools.RemoverSpaci(ttd_penerima.Filename)); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Gagal Menyimpan File Dulu"})
		}
	}

	return c.JSON(fiber.Map{
		"Pesan": "Berhasil Membuat Serahterima Sertifikat",
	})
}
func UpdateSerahterimaSertifikat(c *fiber.Ctx) error {
	// Validasi otentikasi dan data pengguna
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

	// Parse request body
	var request entity.SerahTerimaSertifikat
	if err := c.BodyParser(&request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"Message": "Failed to parse request body", "Error": err.Error()})
	}

	// Ambil ID dari parameter URL
	id := c.Params("id")
	var serahTerima entity.SerahTerimaSertifikat

	// Cari data berdasarkan ID
	if err := database.DB.First(&serahTerima, id).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"Message": "Data not found"})
	}

	// Handle file upload
	ttd_penerima, _ := c.FormFile("ttd_penerima")

	ttd_pemberi, _ := c.FormFile("ttd_pemberi")

	bukti_serah_terima, _ := c.FormFile("bukti_serah_terima")

	tx := database.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if ttd_penerima != nil {
		newPath := "public/static/ttd-penerima/" + tools.RemoverSpaci(ttd_penerima.Filename)
		if err := c.SaveFile(ttd_penerima, newPath); err != nil {
			tx.Rollback()
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Gagal Menyimpan File Dulu",
				"path":  newPath,
			})
		}

		oldPath := "public/static/ttd-penerima/" + serahTerima.TandaTanganPenerima
		if serahTerima.TandaTanganPenerima != "" {
			os.Remove(oldPath)
		}
	}
	if ttd_pemberi != nil {
		newPath := "public/static/ttd-pemberi/" + tools.RemoverSpaci(ttd_pemberi.Filename)
		if err := c.SaveFile(ttd_pemberi, newPath); err != nil {
			tx.Rollback()
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Gagal Menyimpan File Dulu",
				"path":  newPath,
			})
		}

		oldPath := "public/static/ttd-pemberi/" + serahTerima.TandaTanganPemberi
		if serahTerima.TandaTanganPemberi != "" {
			os.Remove(oldPath)
		}
	}

	if bukti_serah_terima != nil {
		newPath := "public/static/bukti-serah-terima/" + tools.RemoverSpaci(bukti_serah_terima.Filename)
		if err := c.SaveFile(bukti_serah_terima, newPath); err != nil {
			tx.Rollback()
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Gagal Menyimpan File Dulu",
				"path":  newPath,
			})
		}

		oldPath := "public/static/bukti-serah-terima/" + serahTerima.BuktiSerahTerima
		if serahTerima.BuktiSerahTerima != "" {
			os.Remove(oldPath)
		}
	}

	// Update data lainnya
	serahTerima.NamaPenerima = request.NamaPenerima
	serahTerima.Jabatan = request.Jabatan
	serahTerima.Instansi = request.Instansi
	serahTerima.NamaKegiatan = request.NamaKegiatan
	serahTerima.TanggalPengambilan = request.TanggalPengambilan
	serahTerima.NoSeriBlanko = request.NoSeriBlanko
	serahTerima.JenisSertifikat = request.JenisSertifikat
	serahTerima.NamaPemberiSertifikat = request.NamaPemberiSertifikat
	serahTerima.Status = request.Status

	// Perbarui timestamp
	serahTerima.UpdateAt = tools.TimeNowJakarta()

	// Simpan perubahan ke database
	if err := database.DB.Save(&serahTerima).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"Message": "Failed to update data"})
	}

	return c.JSON(fiber.Map{"Message": "Serah Terima Sertifikat updated successfully"})
}

func GetSerahterimaSertifikat(c *fiber.Ctx) error {
	// Validasi otentikasi dan data pengguna
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

	// Ambil ID dari parameter URL
	id := c.Query("id")
	var serahTerima []entity.SerahTerimaSertifikat

	query := database.DB

	if id != "" {
		query = query.Where("id_serah_terima_sertifikat =?", id)
	}

	// Cari data berdasarkan ID
	if err := query.Find(&serahTerima).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"Message": "Data not found"})
	}

	return c.JSON(fiber.Map{
		"Message": "Data Serah Terima Sertifikat",
		"data":    serahTerima,
	})
}

func DeleteSerahterimaSertifikat(c *fiber.Ctx) error {
	// Validasi otentikasi dan data pengguna
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

	// Ambil ID dari parameter URL
	id := c.Params("id")
	var serahTerima entity.SerahTerimaSertifikat

	// Cari data berdasarkan ID
	if err := database.DB.First(&serahTerima, id).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"Message": "Data not found"})
	}

	// Lokasi file yang terkait
	path1 := "public/static/ttd-penerima/" + serahTerima.TandaTanganPenerima
	path2 := "public/static/ttd-pemberi/" + serahTerima.TandaTanganPemberi
	path3 := "public/static/bukti-serah-terima/" + serahTerima.BuktiSerahTerima

	// Hapus file jika ada
	if err := os.Remove(path1); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"Message": "Failed to delete Tanda Tangan Penerima file"})
	}

	if err := os.Remove(path2); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"Message": "Failed to delete Tanda Tangan Pemberi file"})
	}

	if err := os.Remove(path3); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"Message": "Failed to delete Bukti Serah Terima file"})
	}

	// Hapus data dari database
	if err := database.DB.Delete(&serahTerima).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"Message": "Failed to delete data"})
	}

	// Return response
	return c.JSON(fiber.Map{"Message": "Serah Terima Sertifikat deleted successfully"})
}
