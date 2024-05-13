package tools

import (
	"fmt"

	"template/app/entity"
	//"template/app/models"

	//"errors"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

//Model Role model yang harus di tambahkan



func GenerateToken(entitys interface{}) string {
    var name, role, types string
    var idAdmin float64 // Mengubah tipe data menjadi string

	fmt.Println(entitys)
    switch e := entitys.(type) {
    case entity.Users:
        name = e.Nama
        idAdmin = float64((int(e.IdUsers))) // Contoh penggunaan ID, sesuaikan dengan kebutuhan Anda
        role = "5"
        types = "Peserta"
    case entity.SuperAdmin:
        name = e.Nama
        idAdmin = float64((int(e.IdSuperAdmin))) // Konversi ID ke string
        role = "1"
        types = "SuperAdmin"

    case entity.Lemdiklat:
        name = e.NamaLemdik
        idAdmin =float64(int(e.IdLemdik))
        role = IntToString(roleLemdiklat)
        types = "Lemdiklat"
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