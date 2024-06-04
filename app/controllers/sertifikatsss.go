package controllers

import (
	"fmt"
	"strings"
	"template/app/entity"
	"template/pkg/database"
	"template/pkg/tools"
	"time"

	"github.com/gofiber/fiber/v2"
)

func AddLastSertifLowBalai(c *fiber.Ctx) error {

	id_admin, _ := c.Locals("id_admin").(int)
	role, _ := c.Locals("role").(string)
	names, _ := c.Locals("name").(string)

	tools.ValidationJwtLemdik(c, role, id_admin, names)

	//ambil pake query

	var data map[string]string

	c.BodyParser(&data)

	var lastRow entity.Lemdiklat

	database.DB.Where("id_lemdik = ?", id_admin).Find(&lastRow)

	//ambil untuk update yang bagian lashrow sertifikat
	noSertif := data["LastNoSertif"]
	updates := entity.Lemdiklat{
		LastNosertif: noSertif,
	}

	database.DB.Model(&lastRow).Updates(&updates)

	//Data Baru

	//Models data last Sertif
	// Kode yang diberikan
	code := noSertif

	// Memisahkan berdasarkan delimiter "/"
	parts := strings.Split(code, "/")

	//numbers := ""

	if len(parts) >= 3 {
		// Memisahkan Nomor
		//numPart := strings.Split(parts[0], ".")[1] // Memisahkan dengan delimiter "."

		// Memisahkan TGL
		//tglPart := strings.Split(parts[1], ".")[1] // Memisahkan dengan delimiter "."

		// Memisahkan Bulan (dalam kasus ini III diwakili sebagai Maret)
		bulanPart := parts[3] // III diwakili sebagai Maret

		// Memisahkan Tahun
		//tahunPart := parts[4]

		// Konversi bulan dari III ke Maret
		switch bulanPart {
		case "I":
			bulanPart = "Januari"
		case "II":
			bulanPart = "Februari"
		case "III":
			bulanPart = "Maret"
		case "IV":
			bulanPart = "April"
		case "V":
			bulanPart = "Mei"
		case "VI":
			bulanPart = "Juni"
		case "VII":
			bulanPart = "Juli"
		case "VIII":
			bulanPart = "Agustus"
		case "IX":
			bulanPart = "September"
		case "X":
			bulanPart = "Oktober"
		case "XI":
			bulanPart = "November"
		case "XII":
			bulanPart = "Desember"
		default:
			bulanPart = "Bulan tidak valid"
		}

		// Output hasil pemisahan
		//numbers = numPart

	} else {
		fmt.Println("Format kode tidak valid")
	}

	NewTableSertifikat := entity.NoSertfikat{
		//IdLemdik:            uint(id_admin),
		//Nomor:               numbers,
		NoLengkapSertifikat: code,
		CreateAt:            tools.TimeNowJakarta(),
	}

	database.DB.Create(&NewTableSertifikat)

	return c.JSON(fiber.Map{
		"Pesan": "Sukses Menambahkan Data No Sertif Keluar",
		"Data":  lastRow.LastNosertif,
	})
}

func GenerateNoSertifikat(c *fiber.Ctx, admin uint) error {

	id_admin, _ := c.Locals("id_admin").(int)
	role, _ := c.Locals("role").(string)
	names, _ := c.Locals("name").(string)

	tools.ValidationJwtLemdik(c, role, id_admin, names)

	//ambil data sertifikat sebelumnya yang oaling terakhir pada last row
	var NoSertif entity.NoSertfikat

	database.DB.Where("id_lemdik = ?", id_admin).Last(&NoSertif)

	//Penambahan tambah dataa nya yaitu 1

	Sertifikat := NoSertif.NoLengkapSertifikat

	parts := strings.Split(Sertifikat, "/")

	singkatan := strings.ToLower(names)
	namess := strings.Split(singkatan, " ")

	balai := ""

	switch namess[1] {
	case "tegal":
		balai = "TGL"
	case "banyuwangi":
		balai = "BWI"
	case "ambon":
		balai = "AMBN"
	case "medan":
		balai = "MDN"
	case "bitung":
		balai = "BTG"
	}

	numbers := 0
	Bulan := time.Now().Month().String()
	Tahun := tools.IntToString(time.Now().Year())

	if len(parts) >= 3 {
		// Memisahkan Nomor
		numPart := strings.Split(parts[0], ".")[1] // Memisahkan dengan delimiter "."

		// Memisahkan TGL
		//tglPart := strings.Split(parts[1], ".")[1] // Memisahkan dengan delimiter "."

		// Memisahkan Bulan (dalam kasus ini III diwakili sebagai Maret)
		// III diwakili sebagai Maret

		// Memisahkan Tahun
		//tahunPart := parts[4]

		// Konversi bulan dari III ke Maret
		switch Bulan {
		case "January":
			Bulan = "I"
		case "Febuari":
			Bulan = "II"
		default:
			Bulan = "Bulan tidak valid"
		}

		// Output hasil pemisahan
		numbers = tools.StringToInt(numPart) + 1

	} else {
		fmt.Println("Format kode tidak valid")
	}
	numbersCuy := tools.IntToString(numbers)
	//Gabungin No Sertifikat nya bree
	NewNoSertifikat := "B." + numbersCuy + "/BPPP." + balai + "/RSDM.510/" + Bulan + "/" + Tahun
	newSertifikat := entity.NoSertfikat{
		//Nomor:               tools.IntToString(numbers),
		NoLengkapSertifikat: NewNoSertifikat,
	}

	database.DB.Create(&newSertifikat)

	return nil

	// B. nomor/
}
