package controllers

import (
	//"backend-elaut/app/entity"

	"os"

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

	tools.ValidationJwtUsers(c, role, id_admin, names)

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
	idAdmin, ok := c.Locals("id_admin").(int)
	if !ok {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid or missing id_admin",
		})
	}

	role, ok := c.Locals("role").(string)
	if !ok {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid or missing role",
		})
	}

	names, ok := c.Locals("name").(string)
	if !ok {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid or missing name",
		})
	}

	tools.ValidationJwtUsers(c, role, idAdmin, names)

	var request models.Users

	if err := c.BodyParser(&request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Failed to parse request body",
			"error":   err.Error(),
		})
	}

	photoProfile, _ := c.FormFile("Foto")

	Ktp, _ := c.FormFile("Ktp")

	KK, _ := c.FormFile("KK")

	ijasah, _ := c.FormFile("Ijazah")

	suratSehat, _ := c.FormFile("SuratKesehatan")

	var user entity.Users
	if err := database.DB.First(&user, idAdmin).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "User not found",
		})
	}

	// Paths for deleting old files
	pathKK := "public/static/profile/kk/" + user.KK
	pathKTP := "public/static/profile/ktp/" + user.Ktp
	pathPoto := "public/static/profile/fotoProfile/" + user.Foto
	pathIzasah := "public/static/profile/ijazah/" + user.Ijazah
	pathSuratSehat := "public/static/profile/suratSehat/" + user.SuratKesehatan

	tx := database.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Handle file uploads and updates
	if photoProfile != nil {
		user.Foto = photoProfile.Filename
		os.Remove(pathPoto)
		if err := c.SaveFile(photoProfile, "public/static/profile/fotoProfile/"+photoProfile.Filename); err != nil {
			tx.Rollback()
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"message": "Failed to save photo profile",
				"error":   err.Error(),
			})
		}
	}

	if Ktp != nil {
		user.Ktp = Ktp.Filename
		os.Remove(pathKTP)
		if err := c.SaveFile(Ktp, "public/static/profile/ktp/"+Ktp.Filename); err != nil {
			tx.Rollback()
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"message": "Failed to save KTP",
				"error":   err.Error(),
			})
		}
	}

	if KK != nil {
		user.KK = KK.Filename
		os.Remove(pathKK)
		if err := c.SaveFile(KK, "public/static/profile/kk/"+KK.Filename); err != nil {
			tx.Rollback()
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"message": "Failed to save KK",
				"error":   err.Error(),
			})
		}
	}

	if ijasah != nil {
		user.Ijazah = ijasah.Filename
		os.Remove(pathIzasah)
		if err := c.SaveFile(ijasah, "public/static/profile/ijazah/"+ijasah.Filename); err != nil {
			tx.Rollback()
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"message": "Failed to save Ijazah",
				"error":   err.Error(),
			})
		}
	}

	if suratSehat != nil {
		user.SuratKesehatan = suratSehat.Filename
		os.Remove(pathSuratSehat)
		if err := c.SaveFile(suratSehat, "public/static/profile/suratSehat/"+suratSehat.Filename); err != nil {
			tx.Rollback()
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"message": "Failed to save Surat Kesehatan",
				"error":   err.Error(),
			})
		}
	}

	// Update user fields
	user.Nama = request.Nama
	user.NoTelpon = tools.StringToInt(request.NoTelpon)
	user.Email = request.Email
	user.Password = request.Password
	user.Kota = request.Kota
	user.Provinsi = request.Provinsi
	user.Alamat = request.Alamat
	user.Nik = tools.StringToInt(request.Nik)
	user.TempatLahir = request.TempatLahir
	user.TanggalLahir = request.TanggalLahir
	user.JenisKelamin = request.JenisKelamin
	user.Pekerjaan = request.Pekerjaan
	user.GolonganDarah = request.GolonganDarah
	user.StatusMenikah = request.StatusMenikah
	user.Kewarganegaraan = request.Kewarganegaraan
	user.IbuKandung = request.IbuKandung
	user.NegaraTujuanBekerja = request.NegaraTujuanBekerja
	user.PendidikanTerakhir = request.PendidikanTerakhir
	user.Agama = request.Agama
	user.Status = request.Status
	user.CreateAt = request.CreateAt
	user.UpdateAt = request.UpdateAt

	// Save the updated user
	if err := tx.Model(&user).Where("id_users = ?", idAdmin).Updates(&user).Error; err != nil {
		tx.Rollback()
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to update user",
			"error":   err.Error(),
		})
	}

	if err := tx.Commit().Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to commit transaction",
			"error":   err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "User updated successfully",
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
