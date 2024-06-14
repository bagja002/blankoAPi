package generator

import (
	"fmt"
	"strings"
	"template/app/entity"
	"template/app/models"
	"template/pkg/database"
	"template/pkg/tools"
	"time"

	"github.com/gofiber/fiber/v2"
)

func GeneratorNoRegister(name string, jenis string, idPel uint, idUsers uint, idLemdik int) string {
	// Cek di database apakah nomor yang akan di-generate sudah ada atau belum
	var noSertif int64
	database.DB.Model(&entity.UsersPelatihan{}).Where("id_pelatihan = ?", idPel).Count(&noSertif)

	// Ambil dari data sebelumnya
	newNoSertif := noSertif + 1

	/* A.201.5.24.0001
	Lemdiklat (A : Tegal),
	Jenis sertifikat (CBIB : 201,
	Waktu ; bulan dan tahun 5.24)
	dan Urutan 0001
	*/

	lowBidangs := strings.ToLower(jenis)

	namaBalai := strings.ToLower(name)
	newBalai := ""
	switch namaBalai {
	case "bppp medan":
		newBalai = "B"
	case "puslat":
		newBalai = "A"
	case "bppp tegal":
		newBalai = "C"
	case "bppp bitung":
		newBalai = "E"
	case "bppp ambon":
		newBalai = "F"
	case "bppp banyuwangi":
		newBalai = "D"
	default:
		newBalai = "BPPSDM"
	}

	//Sertifikasi di mulai dari 200

	//Mungkin pelatihan di mulai dari 100

	bidangPelatihan := ""
	switch lowBidangs {
	case "cpib":
		bidangPelatihan = "201"
	case "cppib":
		bidangPelatihan = "202"
	case "haccp":
		bidangPelatihan = "203"
	case "api":
		bidangPelatihan = "204"
	case "mpm cpib":
		bidangPelatihan = "205"
	case "wisata bahari":
		bidangPelatihan = "WB"
	default:
		bidangPelatihan = "200"
	}

	/* A.201.5.24.0001
	Lemdiklat (A : Tegal),
	Jenis sertifikat (CBIB : 201,
	Waktu ; bulan dan tahun 5.24)
	dan Urutan 0001
	*/

	// Membuat nomor registrasi baru
	tahun := time.Now().Format("2006")
	bulan := time.Now().Format("02")
	/*bulanRomawi := map[string]string{
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

	*/

	bulanString := bulan
	//bulanRomawiString := bulanRomawi[bulanString]
	noSertifFormatted := fmt.Sprintf("%04d", newNoSertif)
	//noPelatihanFormatted := fmt.Sprintf("%04d", idPel)
	NoRegis := fmt.Sprintf("%s.%s.%s.%s.%s", newBalai, bidangPelatihan, bulanString, tahun, noSertifFormatted)

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

func convertUsersPelatihan(entityUsers []entity.Users) []models.Users {
	modelUsersPelatihan := make([]models.Users, len(entityUsers))
	for i, eup := range entityUsers {
		modelUsersPelatihan[i] = models.Users{
			IdUsers: eup.IdUsers,
			Nama:    eup.Nama,
			// Salin field dari eup ke ModelUsersPelatihan sesuai kebutuhan
		}
	}
	return modelUsersPelatihan
}
