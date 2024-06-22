package controllers

import (
	"fmt"
	"path/filepath"
	"template/app/entity"
	"template/pkg/database"
	"template/pkg/tools"

	"github.com/xuri/excelize/v2"

	"github.com/gofiber/fiber/v2"
)

type ResponDatas struct {
	NamaSoal     string
	IdLemdik     uint
	IdPelatihan  uint
	JawabanBenar string
	Jawaban1     string
	Jawaban2     string
	Jawaban3     string
	Jawaban4     string
}

func ImportSoal(c *fiber.Ctx) error {

	id_admin, _ := c.Locals("id_admin").(int)
	role, _ := c.Locals("role").(string)
	names, _ := c.Locals("name").(string)

	tools.ValidationJwtLemdik(c, role, id_admin, names)

	// Membaca file Excel dari request
	file, err := c.FormFile("file")
	if err != nil {
		return err
	}

	var reques ResponDatas

	errs := c.BodyParser(&reques)
	if errs != nil {
		return c.Status(400).JSON(fiber.Map{
			"Pesan": "Gagal menconverisi kan form",
		})
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

	var ListSoal []ResponDatas

	//coba untuk list nya dai berupa

	//Pengubahan Pengubuhan yang dengan AI untuk mengubah terkait dengan Pendidikan DLL, Penangkapan dll
	for _, rowData := range rows[1:] {
		soal := ResponDatas{}
		// Loop melalui setiap kolom dan menambahkan nilai ke objek iku sesuai dengan nama kolom
		for i, columnName := range rows[0] {
			if i >= len(rowData) {
				// Jika indeks i melebihi panjang rowData, lanjutkan ke baris berikutnya
				continue
			}
			switch columnName {
			//Penambahan Validasi Nik dan bersihkand data yang sudah bersih lurr
			case "soal":
				soal.NamaSoal = rowData[i]
			case "jawaban_benar":
				soal.JawabanBenar = rowData[i]
			case "jawaban1":
				soal.Jawaban1 = rowData[i]
			case "jawaban2":
				soal.Jawaban2 = rowData[i]
			case "jawaban3":
				soal.Jawaban3 = rowData[i]
			case "jawaban4":
				soal.Jawaban4 = rowData[i]
				//Bisa ga untuk ke kota aja gitu untuk perapihannya
			}

		}

		soal.IdPelatihan = reques.IdPelatihan

		soal.IdLemdik = uint(id_admin)

		ListSoal = append(ListSoal, soal)
		// Set nilai-nilai lainnya seperti IdPenyuluh, Create_at, Status, dll
	}

	fmt.Println(ListSoal)

	//lakukanperulangan menggunakan for range

	for _, data := range ListSoal {
		fmt.Println("Soal ini adalah : ", data.NamaSoal)
		Pertanyaan := entity.SoalUjianLemdik{
			Soal:        data.NamaSoal,
			IdPelatihan: data.IdPelatihan,
			IdLemdik:    data.IdLemdik,
			Status:      "Active",
		}

		database.DB.Create(&Pertanyaan)

		data_jawaban := []entity.Jawaban{
			{IdSoalUjian: Pertanyaan.IdSoalUjian, NameJawaban: data.JawabanBenar},
			{IdSoalUjian: Pertanyaan.IdSoalUjian, NameJawaban: data.Jawaban1},
			{IdSoalUjian: Pertanyaan.IdSoalUjian, NameJawaban: data.Jawaban2},
			{IdSoalUjian: Pertanyaan.IdSoalUjian, NameJawaban: data.Jawaban3},
			{IdSoalUjian: Pertanyaan.IdSoalUjian, NameJawaban: data.Jawaban4},
		}

		database.DB.Create(&data_jawaban)
		jawabanBenar := data_jawaban[0].NameJawaban

		update := entity.SoalUjianLemdik{
			JawabanBenar: jawabanBenar,
		}
		database.DB.Model(&Pertanyaan).Updates(&update)
	}

	// Logic dari Yang Create Biasa

	return c.JSON(fiber.Map{
		"Pesan": "Soal Berhasil di Import",
	})
}

func GetPertanyaanRandom(c *fiber.Ctx) error {

	id_admin, _ := c.Locals("id_admin").(int)
	role, _ := c.Locals("role").(string)
	names, _ := c.Locals("name").(string)

	tools.ValidationJwtLemdik(c, role, id_admin, names)

	var soal []entity.SoalUjianLemdik

	database.DB.Order("RAND()").Limit(1).Preload("Jawaban").Find(&soal)
	return c.JSON(fiber.Map{
		"Pesan": "Sukses MEngambil Data Random",
		"data":  soal,
	})
}

func AddSoalToUsers(c *fiber.Ctx) error {
	// Ambil id_admin, role, dan names dari context dengan pengecekan error
	idAdmin, ok := c.Locals("id_admin").(int)
	if !ok {
		return c.Status(400).JSON(fiber.Map{
			"Pesan": "ID Admin tidak valid",
		})
	}
	role, ok := c.Locals("role").(string)
	if !ok {
		return c.Status(400).JSON(fiber.Map{
			"Pesan": "Role tidak valid",
		})
	}
	names, ok := c.Locals("name").(string)
	if !ok {
		return c.Status(400).JSON(fiber.Map{
			"Pesan": "Name tidak valid",
		})
	}

	// Validasi JWT dan akses
	if err := tools.ValidationJwtLemdik(c, role, idAdmin, names); err != nil {
		return c.Status(403).JSON(fiber.Map{
			"Pesan": "Unauthorized",
		})
	}

	// Parse body request
	var data map[string]string
	if err := c.BodyParser(&data); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"Pesan": "Gagal menconversi kan form",
		})
	}

	idPelatihan, ok := data["id_pelatihan"]
	if !ok {
		return c.Status(400).JSON(fiber.Map{
			"Pesan": "ID Pelatihan tidak ditemukan",
		})
	}

	var pelatihan entity.Pelatihan
	var soal []entity.SoalUjianLemdik

	// Ambil informasi pelatihan dan juga users pelatihannya
	if err := database.DB.Where("id_pelatihan = ? AND id_lemdik = ?", idPelatihan, idAdmin).Preload("UserPelatihan").Find(&pelatihan).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{
			"Pesan": "Gagal mengambil data pelatihan",
		})
	}

	if err := database.DB.Where("id_pelatihan = ? AND id_lemdik = ?", idPelatihan, idAdmin).Order("RAND()").Limit(10).Preload("Jawaban").Find(&soal).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{
			"Pesan": "Gagal mengambil data soal",
		})
	}

	// Pelatihan Users
	usersPelatihan := pelatihan.UserPelatihan

	// Setiap 1 user dia akan menerima soal random sebanyak 10
	for _, user := range usersPelatihan {
		// Cek setiap perulangan user
		fmt.Println("id_users", user.IdUserPelatihan)
		var users []entity.UsersSoal
		if err := database.DB.Where("id_user_pelatihan = ?", user.IdUserPelatihan).Find(&users).Error; err != nil {
			return c.Status(500).JSON(fiber.Map{
				"Pesan": "Gagal memeriksa data users soal",
			})
		}

		if len(users) > 10 {
			continue
		}

		user.CodeAksess = tools.RandomString(6)
		if err := database.DB.Model(&user).Updates(&user).Error; err != nil {
			return c.Status(500).JSON(fiber.Map{
				"Pesan": "Gagal memperbarui kode akses user",
			})
		}

		for _, s := range soal {
			// Users soal
			usersSoal := entity.UsersSoal{
				IdUserPelatihan: user.IdUserPelatihan,
				IdSoalUjian:     s.IdSoalUjian,
				CreateAt:        tools.TimeNowJakarta(),
			}

			if err := database.DB.Create(&usersSoal).Error; err != nil {
				return c.Status(500).JSON(fiber.Map{
					"Pesan": "Gagal menambah soal ke user",
				})
			}
		}
	}

	return c.JSON(fiber.Map{
		"Pesan": "Berhasil menambah soal ke users",
	})
}

func AuthExam(c *fiber.Ctx) error {

	var data map[string]string

	if err := c.BodyParser(&data); err != nil {
		c.JSON(fiber.Map{
			"pesan": err.Error(),
		})
	}

	codeAkses := data["code_akses"]
	if codeAkses == "" {
		return c.Status(400).JSON(fiber.Map{
			"Pesan": "Mohon Maaf Masukan Kode Akses Anda",
		})
	}

	var users entity.UsersPelatihan
	database.DB.Where("code_aksess = ? ", codeAkses).Find(&users)
	if users.CodeAksess == "" {
		return c.Status(400).JSON(fiber.Map{
			"Pesan": "Silahkan Daftar Pelatihan",
		})
	}

	//Generate Token
	t := tools.GenerateTokenExam(users)

	return c.JSON(fiber.Map{
		"t": t,
	})
}

//Fungsi Unuk mengambil soal dari uersSoal Pelatihan

func GetSoalUsers(c *fiber.Ctx) error {

	id_users, _ := c.Locals("id_users").(int)
	id_users_pelatihan, _ := c.Locals("id_users_pelatihan").(int)
	types, _ := c.Locals("types").(string)

	if id_users == 0 {

	}

	if types == "" {

	}

	//ambil soal dari table users soal berdasarkan id users pelatihan
	var soalUsers []entity.UsersSoal

	database.DB.Where("id_user_pelatihan = ? ", id_users_pelatihan).Find(&soalUsers)

	idSoalUjian := []int64{}

	for _, id_soal := range soalUsers {
		idSoalUjian = append(idSoalUjian, int64(id_soal.IdSoalUjian))
	}

	//Mengambil soal dan jawaban dari users

	var soal []entity.SoalUjianLemdik

	database.DB.Where("id_soal_ujian IN ?", idSoalUjian).Preload("Jawaban").Find(&soal)

	//ambil soal berdasarkan

	return c.JSON(fiber.Map{
		"Soal":   soal,
		"jumlah": len(idSoalUjian),
	})

}

func Jawab(c *fiber.Ctx) error {

	id_users, _ := c.Locals("id_users").(int)
	id_users_pelatihan, _ := c.Locals("id_users_pelatihan").(int)
	types, _ := c.Locals("types").(string)

	if id_users == 0 {

	}

	if types == "" {

	}

	var jawaban []struct {
		IdSoalLemdik    string `json:"id_soal_lemdik"`
		JawabanPengguna string `json:"jawaban_pengguna"`
	}

	// Parse request body
	if err := c.BodyParser(&jawaban); err != nil {
		fmt.Println("Error parsing body:", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Gagal memparsing body",
		})
	}

	ScoreBenar := 0
	//var jawabanBenarCount int
	for _, jwb := range jawaban {
		var soal entity.SoalUjianLemdik
		result := database.DB.Where("id_soal_ujian = ?", jwb.IdSoalLemdik).Find(&soal)

		if result.Error != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"message": "Gagal mendapatkan soal",
			})
		}

		if jwb.JawabanPengguna == soal.JawabanBenar {
			ScoreBenar = ScoreBenar + 1
		}
	}

	//Perhitungan Score Berdasarkan jumlan soal dan soal benar

	jumlahSoal := len(jawaban)
	fmt.Println(jumlahSoal)

	finalScore := (float64(ScoreBenar) / float64(jumlahSoal)) * 100.0

	//Kirim data ke API E-Laut untuk Peserta Ujjian
	var usersPelatihan entity.UsersPelatihan
	database.DB.Where("id_user_pelatihan = ?", id_users_pelatihan).Find(&usersPelatihan)

	//ubah nilainya
	usersPelatihan.PreTest = int(finalScore)

	database.DB.Model(&usersPelatihan).Updates(&usersPelatihan)

	//Jika telah mengerjakan Post test Hapus Soal yang ada di userSoal Agar data tidak numpuk

	return c.JSON(fiber.Map{
		"Pesan": "Terima Kasih Telah Megerjakan Ujian",
	})
}
