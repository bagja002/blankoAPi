package tools

import (
	"template/app/entity"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

func GenerateTokenExam(entitys interface{}) string {
	var types string
	var id_users float64 // Mengubah tipe data menjadi string
	var id_users_pelatihan float64

	switch e := entitys.(type) {
	case entity.UsersPelatihan:
		id_users_pelatihan = float64((int(e.IdUserPelatihan)))
		id_users = float64((int(e.IdUsers)))
		// Contoh penggunaan ID, sesuaikan dengan kebutuhan And
		types = "PreTest"
	default:
		return ""
	}

	claims := jwt.MapClaims{
		"id_users_pelatihan": id_users_pelatihan,
		"id_users":           id_users,
		"type":               types,
		"exp":                time.Now().Add(time.Hour * 72).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	t, _ := token.SignedString([]byte("secret"))
	return t
}
func GenerateTokenExamPostTest(entitys interface{}) string {
	var types string
	var id_users float64 // Mengubah tipe data menjadi string
	var id_users_pelatihan float64

	switch e := entitys.(type) {
	case entity.UsersPelatihan:
		id_users_pelatihan = float64((int(e.IdUserPelatihan)))
		id_users = float64((int(e.IdUsers)))
		types = "PostTest"
	default:
		return ""
	}

	claims := jwt.MapClaims{
		"id_users_pelatihan": id_users_pelatihan,
		"id_users":           id_users,
		"type":               types,
		"exp":                time.Now().Add(time.Hour * 72).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	t, _ := token.SignedString([]byte("secret"))
	return t
}
