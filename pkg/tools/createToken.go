package tools

import (
	"github.com/golang-jwt/jwt/v4"
	"template/app/entity"
	"time"
)

//Model Role model yang harus di tambahkan

func GenerateToken(entitys interface{}) string {
	var name, role, types string
	var idAdmin float64 // Mengubah tipe data menjadi string

	switch e := entitys.(type) {
	case entity.Admin:
		name = e.Nama
		idAdmin = float64((int(e.IdAdmin))) // Contoh penggunaan ID, sesuaikan dengan kebutuhan Anda
		role = "1"
		types = "Admin Pusat"
	case entity.SuperAdmin:
		name = e.Nama
		idAdmin = float64((int(e.IdSuperAdmin))) // Konversi ID ke string
		role = "99"
		types = "SuperAdmin"
	default:
		return ""
	}

	claims := jwt.MapClaims{
		"name":     name,
		"id_admin": idAdmin,
		"role":     role,
		"type":     types,
		"exp":      time.Now().Add(time.Hour * 72).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	t, _ := token.SignedString([]byte("secret"))
	return t
}
