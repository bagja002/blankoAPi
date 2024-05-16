package controllers

import (
	//"backend-elaut/app/entity"

	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
	"template/app/entity"
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
		Nama:     data["nama"],
		Nik:      nik,
		Password: tools.GeneratePassword(data["password"]),
		NoTelpon: tools.StringToInt(data["no_number"]),
		CreateAt: tools.TimeNowJakarta(),
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

	var user entity.Users
	if err := database.DB.Find(&user, id_admin).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"Message": "User not found",
		})
	}

	var data map[string]string
	if err := c.BodyParser(&data); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"Message": "Invalid request body",
		})
	}

	// Update user fields
	update := entity.Users{
		Nama:                 data["nama"],
		NoTelpon:             tools.StringToInt(data["no_telpon"]),
		Email:                data["email"],
		Password:             data["password"],
		Kota:                 data["kota"],
		Provinsi:             data["provinsi"],
		Alamat:               data["alamat"],
		Nik:                  tools.StringToInt(data["nik"]),
		TempatLahir:          data["tempat_lahir"],
		TanggalLahir:         data["tanggal_lahir"],
		JenisKelamin:         data["jenis_kelamin"],
		Pekerjaan:            data["pekerjaan"],
		GolonganDarah:        data["golongan_darah"],
		StatusMenikah:        data["status_menikah"],
		Kewarganegaraan:      data["kewarganegaraan"],
		IbuKandung:           data["ibu_kandung"],
		NegaraTujuanBekerja:  data["negara_tujuan_bekerja"],
		PendidikanTerakhir:   data["pendidikan_terakhir"],
		Agama:                data["agama"],
		Foto:                 data["foto"],
		Status:               data["status"],
		CreateAt:             data["create_at"],
		UpdateAt:             data["update_at"],
	}
	


	// Save the updated user
	if err := database.DB.Model(&user).Updates(&update).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"Message": "Failed to update user",
		})
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
