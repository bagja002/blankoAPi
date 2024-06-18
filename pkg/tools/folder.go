package tools

import (
	"fmt"
	"os"
)

func CreateFolder() {
	folders := []string{
		"public/silabus/pelatihan",
		"public/silabus/sertifikasi",
		"public/module/pelatihan",

		"public/static/sertifikasi",
		"public/static/prasarana",
		"public/static/pelatihan",
		"public/static/profile/fotoProfile",
		"public/static/profile/ijazah",
		"public/static/profile/kk",
		"public/static/profile/ktp",
		"public/static/profile/suratSehat",

		"public/static/suratPemberitahuan",
		"public/static/BeritaAcara",
		"public/static/fileSertifikat",

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
