package controllers

import (
	//"backend-elaut/app/entity"

	"os"
	"strings"
	"template/app/entity"
	"template/app/models"

	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"

	//"template/app/models"
	"template/pkg/database"
	"template/pkg/tools"
)

func CreateUser(c *fiber.Ctx) error {

	var data map[string]string

	if err := c.BodyParser(&data); err != nil {
		c.JSON(fiber.Map{
			"pesan": err.Error(),
		})
	}

	//cek email apakah sudah ada

	var existingEmail entity.Users
	nik := tools.StringToInt(data["nik"])
	err := database.DB.Where("nik = ?", nik).Find(&existingEmail).Error
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"Message": err,
		})
	}

	//Integrasi dengan pengecekan sistem dari DUkcapil

	//cek email
	if existingEmail.Nik == nik {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"Message": "NIK Telah Terdaftar is Register",
		})
	}

	NewUser := entity.Users{
		Nama:        data["nama"],
		Nik:         nik,
		Password:    tools.GeneratePassword(data["password"]),
		NoTelpon:    tools.StringToInt(data["no_number"]),
		KusukaUsers: data["kusuka_users"],
		CreateAt:    tools.TimeNowJakarta(),
	}
	//Memulai Transaksi

	tx := database.DB.Begin()

	if err := tx.Create(&NewUser).Error; err != nil {
		tx.Rollback()
		return err
	}

	// Commit transaksi jika semuanya berhasil
	if err := tx.Commit().Error; err != nil {
		return err
	}

	//database.DB.Create(&NewUser)

	return c.JSON(fiber.Map{
		"Pesan": "Selamat Anda Sudah mendaftar",
	})
}

func LoginUsers(c *fiber.Ctx) error {

	var data map[string]string

	if err := c.BodyParser(&data); err != nil {
		c.JSON(fiber.Map{
			"pesan": err.Error(),
		})
	}

	var users entity.Users

	database.DB.Where("nik = ? ", data["nik"]).Find(&users)
	if users.IdUsers == 0 {
		c.Status(fiber.StatusUnauthorized)
		return c.JSON(fiber.Map{
			"pesan": "Akun tidak di temukan di temukan",
		})
	}

	if err := bcrypt.CompareHashAndPassword([]byte(users.Password), []byte(data["password"])); err != nil {
		c.Status(fiber.StatusUnauthorized)
		return c.JSON(fiber.Map{
			"pesan": "Incorrect password!",
		})
	} else {
		t := tools.GenerateToken(users)
		return c.JSON(fiber.Map{
			"t": t,
		})
	}
}

// Get a user by ID
func GetUserByID(c *fiber.Ctx) error {

	id_admin, _ := c.Locals("id_admin").(int)
	role, _ := c.Locals("role").(string)
	names, _ := c.Locals("name").(string)

	tools.ValidationJwt(c, role, id_admin, names)

	var user entity.Users
	if err := database.DB.Find(&user, id_admin).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"Message": "User not found",
		})
	}

	return c.JSON(user)
}

// Get all users
func GetAllUsers(c *fiber.Ctx) error {

	id := c.Query("id")

	if id != "" {
		var users entity.Users
		if err := database.DB.Where("id_users = ?", id).Find(&users).Error; err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"Message": "Failed to fetch users",
			})
		}

		return c.JSON(fiber.Map{
			"pesan": "Sukses",
			"data":  users,
		})

	}

	var users []entity.Users
	if err := database.DB.Find(&users).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"Pesan": "Failed to fetch users",
		})
	}

	return c.JSON(fiber.Map{
		"pesan": "Sukses",
		"data":  users,
	})
}

// Update user by ID
func UpdateUser(c *fiber.Ctx) error {
	id_admin, _ := c.Locals("id_admin").(int)
	role, _ := c.Locals("role").(string)
	names, _ := c.Locals("name").(string)

	tools.ValidationJwtUsers(c, role, id_admin, names)

	photoProfile, err := c.FormFile("Foto")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"Message": "Failed to retrieve file", "Error": err.Error()})
	}
	Ktp, err := c.FormFile("Ktp")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"Message": "Failed to retrieve file", "Error": err.Error()})
	}
	KK, err := c.FormFile("KK")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"Message": "Failed to retrieve file", "Error": err.Error()})
	}
	ijasah, err := c.FormFile("Ijazah")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"Message": "Failed to retrieve file", "Error": err.Error()})
	}

	suratSehat, err := c.FormFile("SuratKesehatan")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"Message": "Failed to retrieve file", "Error": err.Error()})
	}
	var request models.Users
	if err := c.BodyParser(&request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"Message": "Failed to parse request body", "Error": err.Error()})
	}

	var user entity.Users
	if err := database.DB.Find(&user, id_admin).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"Message": "User not found",
		})
	}
	//untuk delete file
	pathKK := "public/static/profile/kk/" + user.KK
	pathKTP := "public/static/profile/ktp/" + user.Ktp
	pathPoto := "public/static/profile/fotoProfile/" + user.Foto
	pathIzasah := "public/static/profile/ijazah/" + user.Ijazah
	pathSuratSehat := "public/static/profile/suratSehat/" + user.SuratKesehatan

	os.Remove(pathKK)
	os.Remove(pathKTP)
	os.Remove(pathPoto)
	os.Remove(pathIzasah)
	os.Remove(pathSuratSehat)
	// Update user fields
	update := entity.Users{
		Nama:                request.Nama,
		NoTelpon:            request.NoTelpon,
		Email:               request.Email,
		Password:            request.Password,
		Kota:                request.Kota,
		Provinsi:            request.Provinsi,
		Alamat:              request.Alamat,
		Nik:                 request.Nik,
		TempatLahir:         request.TempatLahir,
		TanggalLahir:        request.TanggalLahir,
		JenisKelamin:        request.JenisKelamin,
		Pekerjaan:           request.Pekerjaan,
		GolonganDarah:       request.GolonganDarah,
		StatusMenikah:       request.StatusMenikah,
		Kewarganegaraan:     request.Kewarganegaraan,
		IbuKandung:          request.IbuKandung,
		NegaraTujuanBekerja: request.NegaraTujuanBekerja,
		PendidikanTerakhir:  request.PendidikanTerakhir,
		Agama:               request.Agama,
		Foto:                photoProfile.Filename, //ganti format nama File Foto dengan Iser dan nama
		Status:              request.Status,
		CreateAt:            request.CreateAt,
		UpdateAt:            request.UpdateAt,
		KK:                  KK.Filename,
		Ktp:                 Ktp.Filename,
		Ijazah:              ijasah.Filename,
		SuratKesehatan:      suratSehat.Filename,
	}

	// Save the updated user
	if err := database.DB.Model(&user).Updates(&update).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"Message": "Failed to update user",
		})
	}

	if err := c.SaveFile(photoProfile, "public/static/profile/fotoProfile/"+strings.ReplaceAll(photoProfile.Filename, " ", "")); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"Message": "Failed to save fotoProfile", "Error": err.Error()})
	}
	if err := c.SaveFile(KK, "public/static/profile/kk/"+strings.ReplaceAll(KK.Filename, " ", "")); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"Message": "Failed to save kk", "Error": err.Error()})
	}
	if err := c.SaveFile(Ktp, "public/static/profile/ktp/"+strings.ReplaceAll(Ktp.Filename, " ", "")); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"Message": "Failed to save ktp", "Error": err.Error()})
	}
	if err := c.SaveFile(ijasah, "public/static/profile/ijazah/"+strings.ReplaceAll(ijasah.Filename, " ", "")); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"Message": "Failed to save ijazah", "Error": err.Error()})
	}
	if err := c.SaveFile(suratSehat, "public/static/profile/suratSehat/"+strings.ReplaceAll(suratSehat.Filename, " ", "")); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"Message": "Failed to save suratSehat", "Error": err.Error()})
	}

	return c.JSON(fiber.Map{
		"Message": "User updated successfully",
	})
}

// Delete user by ID
func DeleteUser(c *fiber.Ctx) error {
	userID := c.Params("id")

	var user entity.Users
	if err := database.DB.First(&user, userID).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"Message": "User not found",
		})
	}

	// Delete the user
	if err := database.DB.Delete(&user).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"Message": "Failed to delete user",
		})
	}

	return c.JSON(fiber.Map{
		"Message": "User deleted successfully",
	})
}

//Generat Referar Token
