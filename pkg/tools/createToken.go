package tools

import (
	"template/app/models"
	//"errors"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

//Model Role model yang harus di tambahkan




func GenerateToken(entitys interface{}) (string) {
	var name, idAdmin, role, satminkal string

	switch e := entitys.(type) {
	case models.User :
		name = e.Name
		idAdmin = e.ID
		role = e.Role
		satminkal = e.Satminkal
	default:
		return ""
	}

	claims := jwt.MapClaims{
		"name":      name,
		"id_admin":  idAdmin,
		"role":      role,
		"satminkal": satminkal,
		"exp":       time.Now().Add(time.Hour * 72).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	t,_:= token.SignedString([]byte("secret"))
	return t
}
