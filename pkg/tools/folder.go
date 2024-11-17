package tools

import (
	"fmt"
	"os"
)

func CreateFolder() {
	folders := []string{

		"public/static/foto-blanko-rusak",

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
