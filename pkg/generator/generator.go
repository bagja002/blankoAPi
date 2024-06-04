package generator

import (
	"fmt"
	"strings"
	"template/app/entity"
	"template/pkg/database"
	//"template/pkg/tools"
	"time"
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
		bidangPelatihan = "BP"
	case "pengolahan":
		bidangPelatihan = "PP"
	case "pemasaran":
		bidangPelatihan = "PM"
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
