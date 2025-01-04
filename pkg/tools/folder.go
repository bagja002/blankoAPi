package tools

import (
	"fmt"
	"os"
)

func CreateFolder() {
	folders := []string{

		"public/static/foto-blanko-rusak",

		"public/static/ttd-penerima",
		"public/static/ttd-pemberi",
		"public/static/bukti-serah-terima",

		"public/static/ttd-pengiriman",
		"public/static/bukti-resi",
		"public/static/bukti-pengiriman-sertifikat",
		"public/static/bukti-penerimaan-sertifikat",

		//Folder Untuk Foto Sarpras

	}

	for _, folder := range folders {
		err := os.MkdirAll(folder, os.ModePerm)
		if err != nil {
			fmt.Printf("Error creating directory %s: %v\n", folder, err)
			return
		}

	}

}
