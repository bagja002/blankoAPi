package generator

import (
	"fmt"
	"strings"
	"template/app/entity"
	"template/pkg/database"
	"template/pkg/tools"
	"time"

	"github.com/gofiber/fiber/v2"
)

func GeneratorNoRegister(name string, bidang string, idPel uint, idUsers uint, idLemdik int) string {
	// Cek di database apakah nomor yang akan di-generate sudah ada atau belum
	var noSertif int64
	database.DB.Model(&entity.UsersPelatihan{}).Where("id_pelatihan = ?", idPel).Count(&noSertif)

	// Ambil dari data sebelumnya
	newNoSertif := noSertif + 1

	lowBidangs := strings.ToLower(bidang)

	namaBalai := strings.ToLower(name)
	newBalai := ""
	switch namaBalai {
	case "bppp medan":
		newBalai = "MDN"
	case "bppp tegal":
		newBalai = "TGL"
	case "bppp bitung":
		newBalai = "BTNG"
	case "bppp ambon":
		newBalai = "AMBN"
	case "bppp banyuwangi":
		newBalai = "BWI"
	default:
		newBalai = "BPPSDM"
	}

	bidangPelatihan := ""
	switch lowBidangs {
	case "budidaya":
		bidangPelatihan = "BD"
	case "pengolahan dan pemasaran":
		bidangPelatihan = "PM"
	case "penangkapan":
		bidangPelatihan = "PK"
	case "mesin perikanan":
		bidangPelatihan = "MP"
	case "konservasi":
		bidangPelatihan = "KV"
	case "wisata bahari":
		bidangPelatihan = "WB"
	default:
		bidangPelatihan = "XX"
	}

	// Membuat nomor registrasi baru
	tanggal := time.Now().Format("2006")
	bulan := time.Now().Format("02")
	bulanRomawi := map[string]string{
		"01": "I",
		"02": "II",
		"03": "III",
		"04": "IV",
		"05": "V",
		"06": "VI",
		"07": "VII",
		"08": "VIII",
		"09": "IX",
		"10": "X",
		"11": "XI",
		"12": "XII",
	}

	bulanString := bulan
	bulanRomawiString := bulanRomawi[bulanString]
	noSertifFormatted := fmt.Sprintf("%04d", newNoSertif)
	noPelatihanFormatted := fmt.Sprintf("%04d", idPel)
	NoRegis := fmt.Sprintf("%s.%s.%s.%s.%s.%s", newBalai, bidangPelatihan, noPelatihanFormatted, bulanRomawiString, tanggal, noSertifFormatted)

	// Simpan nomor registrasi ke database
	/*
		newRegistrasi := entity.NoRegistrasi{
			Nomor:       int(newNoSertif),
			IdPelatihan: idPel,
			IdUsers:     idUsers,
			NamaLemdik:  name,
			Bidang:      bidang,
			CreateAt:    tools.TimeNowJakarta(),
		}
		database.DB.Create(&newRegistrasi)

	*/

	return NoRegis
}

func GenerateSertifikat(idLemdik string, idPelatihan string, c *fiber.Ctx) string {

	var lastSertif entity.Sertifikat

	database.DB.Where("id_lemdik = ?", idLemdik).Last(&lastSertif)

	lastSertifikat := lastSertif.NoSertfikat

	splitted := strings.Split(lastSertifikat, "/")

	newSertifkatBaru := ""

	if len(splitted) >= 5 {
		nomorSeriAwal := splitted[0]
		balai := splitted[1]
		rsdm := splitted[2]
		bulan := time.Now().Month().String()

		bulanRomawi := map[string]string{
			"January":   "I",
			"February":  "II",
			"March":     "III",
			"April":     "IV",
			"May":       "V",
			"June":      "VI",
			"July":      "VII",
			"August":    "VIII",
			"September": "IX",
			"October":   "X",
			"November":  "XI",
			"December":  "XII",
		}
		bulanString := bulan
		bulanRomawiString := bulanRomawi[bulanString]
		tahun := splitted[4]

		// Mengubah nomor seri awal menjadi integer
		NoSertifAwal := tools.StringToInt(nomorSeriAwal)

		// Mendapatkan nomor seri berikutnya
		nomorSeriBaru := NoSertifAwal + 1

		// Memformat nomor seri baru
		nomorSeriBaruString := fmt.Sprintf("%03d", nomorSeriBaru)

		// Membentuk string sertifikat baru
		sertifikatBaru := nomorSeriBaruString + "/" + balai + "/" + rsdm + "/" + bulanRomawiString + "/" + tahun

		// Mengembalikan string sertifikat baru
		newSertifkatBaru = sertifikatBaru
	}
	//simpan NewSertif baru di simpan di database
	newDatabaseSertif := entity.Sertifikat{
		IdLemdik:    uint(tools.StringToInt(idLemdik)),
		IdPelatihan: uint(tools.StringToInt(idPelatihan)),
		NoSertfikat: newSertifkatBaru,
		CreateAt:    tools.TimeNowJakarta(),
	}

	database.DB.Create(&newDatabaseSertif)

	return newSertifkatBaru
}
