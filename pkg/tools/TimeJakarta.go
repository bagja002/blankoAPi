package tools

import "time"

func TimeNowJakarta() string {
	loc, err := time.LoadLocation("Asia/Jakarta")
	if err != nil {
		panic(err)
	}

	// Mendapatkan waktu saat ini dalam zona waktu Jakarta
	now := time.Now().In(loc)

	// Menggunakan waktu Jakarta untuk mengatur iku.Create_at
	timestring := now.Format("02 January 2006, 03:04 PM")

	return string(timestring)
}