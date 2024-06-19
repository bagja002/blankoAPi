package controllers

import (
	"fmt"
	"path/filepath"
	"template/app/entity"
	"template/pkg/database"
	"template/pkg/generator"
	"template/pkg/tools"

	"github.com/gofiber/fiber/v2"
	"github.com/xuri/excelize/v2"
)

func ExportPesertaPelatihan(c *fiber.Ctx) error {

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

	var models entity.Pelatihan

	if err := c.BodyParser(&models); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"Message": "Failed to parse request body", "Error": err.Error()})
	}

	idPelatihan := models.IdPelatihan

	//var userPelatihanList []entity.UsersPelatihan
	var users []entity.Users

	for _, rowUsers := range rows[1:] {
		user := entity.Users{}

		for i, columnName := range rows[0] {
			if i >= len(rowUsers) {
				// Jika indeks i melebihi panjang rowData, lanjutkan ke baris berikutnya
				continue
			}
			switch columnName {
			case "nama":
				user.Nama = rowUsers[i]
			case "no_telpon":
				user.NoTelpon = tools.StringToInt(rowUsers[i])
			case "email":
				user.Email = rowUsers[i]
			case "password":
				user.Password = tools.GeneratePassword(rowUsers[i])
			case "kota":
				user.Kota = rowUsers[i]
			case "provinsi":
				user.Provinsi = rowUsers[i]
			case "alamat":
				user.Alamat = rowUsers[i]
			case "nik":
				user.Nik = tools.StringToInt(rowUsers[i])
			case "tempat_lahir":
				user.TempatLahir = rowUsers[i]
			case "tanggal_lahir":
				user.TanggalLahir = rowUsers[i]
			case "jenis_kelamin":
				user.JenisKelamin = rowUsers[i]
			case "pekerjaan":
				user.Pekerjaan = rowUsers[i]
			case "golongan_darah":
				user.GolonganDarah = rowUsers[i]
			case "status_menikah":
				user.StatusMenikah = rowUsers[i]
			case "kewarganegaraan":
				user.Kewarganegaraan = rowUsers[i]
			case "ibu_kandung":
				user.IbuKandung = rowUsers[i]
			case "negara_tujuan_bekerja":
				user.NegaraTujuanBekerja = rowUsers[i]
			case "pendidikan_terakhir":
				user.PendidikanTerakhir = rowUsers[i]
			case "agama":
				user.Agama = rowUsers[i]
			case "foto":
				user.Foto = rowUsers[i]
			case "ktp":
				user.Ktp = rowUsers[i]
			case "kk":
				user.KK = rowUsers[i]
			case "surat_kesehatan":
				user.SuratKesehatan = rowUsers[i]
			case "ijazah_users":
				user.Ijazah = rowUsers[i]
			}

			//Masukan List Ke dalam Users List  List

		}
		users = append(users, user)

		//Masukan Data Semuanya

	}

	for _, AllUsers := range users {

		//Buat Id Users terlebih Dahulu
		database.DB.Create(&AllUsers)
		var Pelatihan entity.Pelatihan
		database.DB.Where("id_pelatihan =?", idPelatihan).Find(&Pelatihan)

		//ambil id lemdiknya trius ambil Id namanya
		var lemdik entity.Lemdiklat
		database.DB.Where("id_lemdik = ? ", Pelatihan.IdLemdik).Find(&lemdik)

		NoRegistrasi := generator.GeneratorNoRegister(lemdik.NamaLemdik, Pelatihan.BidangPelatihan, Pelatihan.IdPelatihan, uint(id_admin), int(lemdik.IdLemdik))

		dataUsers := entity.UsersPelatihan{
			IdUsers:            AllUsers.IdUsers,
			IdPelatihan:        idPelatihan,
			Nama:               AllUsers.Nama,
			NoRegistrasi:       NoRegistrasi,
			TotalBayar:         tools.IntToString(Pelatihan.HargaPelatihan),
			TempatTanggalLahir: AllUsers.TempatLahir + AllUsers.TanggalLahir,
			NamaPelatihan:      Pelatihan.NamaPelatihan,
			BidangPelatihan:    Pelatihan.BidangPelatihan,
			DetailPelatihan:    Pelatihan.DetailPelatihan,
			StatusAproval:      Pelatihan.StatusApproval,
			TanggalMulai:       Pelatihan.TanggalMulaiPelatihan,
			TanggalBerakhir:    Pelatihan.TanggalBerakhirPelatihan,
			StatusPembayaran:   "done",
		}

		database.DB.Create(&dataUsers)

	}

	return c.JSON(fiber.Map{
		"Pesan":      "Sukses Upload Data Peserta Pelatihan ",
		"data":       users,
		"Total Data": len(users),
	})
}

func ExportMateriPelatihan(c *fiber.Ctx) error {

	return nil
}
