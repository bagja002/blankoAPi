package controllers

import (
	"template/app/entity"
	"template/pkg/database"
	"template/pkg/generator"
	"template/pkg/tools"

	"github.com/gofiber/fiber/v2"
)

// From Register Akun User
func CreateUserPelatihan(c *fiber.Ctx) error {

	id_admin, _ := c.Locals("id_admin").(int)
	role, _ := c.Locals("role").(string)
	names, _ := c.Locals("name").(string)

	tools.ValidationJwtUsers(c, role, id_admin, names)

	var data map[string]string

	if err := c.BodyParser(&data); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"Message": "Failed to parse request body", "Error": err.Error()})
	}
	var exixtingPelatihanUsers entity.UsersPelatihan
	idPelatihan := uint(tools.StringToInt(data["id_pelatihan"]))

	database.DB.Where("id_users = ? AND id_pelatihan = ?", id_admin, idPelatihan).Find(&exixtingPelatihanUsers)

	if exixtingPelatihanUsers.IdUserPelatihan != 0 {
		return c.Status(400).JSON(fiber.Map{
			"Pesan": "Anda Sudah Mendaftar Pelatihan",
		})
	}
	//
	var Pelatihan entity.Pelatihan
	database.DB.Where("id_pelatihan =?", idPelatihan).Find(&Pelatihan)

	//ambil id lemdiknya trius ambil Id namanya
	var lemdik entity.Lemdiklat
	database.DB.Where("id_lemdik = ? ", Pelatihan.IdLemdik).Find(&lemdik)

	NoRegistrasi := generator.GeneratorNoRegister(lemdik.NamaLemdik, Pelatihan.BidangPelatihan, Pelatihan.IdPelatihan, uint(id_admin), int(lemdik.IdLemdik))

	newUserPelatihan := entity.UsersPelatihan{
		IdUsers:          uint(id_admin),
		Nama:             names,
		IdPelatihan:      idPelatihan,
		NoRegistrasi:     NoRegistrasi,
		TotalBayar:       data["totalBayar"],
		StatusPembayaran: "pending",
		CreteAt:          tools.TimeNowJakarta(),
	}

	database.DB.Create(&newUserPelatihan)

	var testing entity.UsersPelatihan

	database.DB.Preload("Users").Find(&testing)

	return c.JSON(fiber.Map{
		"pesan": "berhasil membuat data",
		"data":  newUserPelatihan,
	})
}

func UpdateUsersPelatihan(c *fiber.Ctx) error {

	id_admin, _ := c.Locals("id_admin").(int)
	role, _ := c.Locals("role").(string)
	names, _ := c.Locals("name").(string)

	tools.ValidationJwtLemdik(c, role, id_admin, names)

	//INI CARANYA BAGAIMANA
	id := c.Query("id")
	FileSertifikat, _ := c.FormFile("FileSertifikat")

	var usersPelatihan entity.UsersPelatihan

	database.DB.Where("id_user_pelatihan = ?", id).Find(&usersPelatihan)

	// Menginisialisasi koneksi database
	tx := database.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if FileSertifikat != nil {
		if FileSertifikat != nil {
			usersPelatihan.FileSertifikat = FileSertifikat.Filename
			if err := c.SaveFile(FileSertifikat, "public/static/fileSertifikat/"+FileSertifikat.Filename); err != nil {
				tx.Rollback()
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
					"pesan": "Gagal menyimpan file EvaluasiRenaksi",
					"error": err.Error(),
				})
			}
		}

		if err := tx.Model(&usersPelatihan).Where("id_user_pelatihan = ?", id).Updates(&usersPelatihan).Error; err != nil {
			tx.Rollback()
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"pesan": "Gagal memperbarui MonitoringEvaluasi",
				"error": err.Error(),
			})
		}

	}

	//data biasa
	var request entity.UsersPelatihan

	if err := c.BodyParser(&request); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"pesan": "gagal reques",
		})
	}

	updates := entity.UsersPelatihan{
		NoSertifikat: request.NoSertifikat,
		NoRegistrasi: request.NoRegistrasi,
		PreTest:      request.PreTest,
		PostTest:     request.PostTest,
		NilaiTeory:   request.NilaiTeory,
		NilaiPraktek: request.NilaiPraktek,

		//Nilai Materi
		StatusPembayaran: request.StatusPembayaran, //Pending dan Void
		MetodoPembayaran: request.MetodoPembayaran,
		WaktuPembayaran:  request.WaktuPembayaran,
		Keterangan:       request.Keterangan,

		UpdateAt: tools.TimeNowJakarta(),
	}

	if err := tx.Model(&usersPelatihan).Where("id_pelatihan = ?", id).Updates(&updates).Error; err != nil {
		tx.Rollback()
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"pesan": "Gagal memperbarui MonitoringEvaluasi",
			"error": err.Error(),
		})
	}

	if err := tx.Commit().Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"pesan": "Gagal melakukan commit transaksi",
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"Pesan": "Berhasil Update Pelatihan ",
		"data":  usersPelatihan,
	})
}

func GetUserPelatihan(c *fiber.Ctx) error {

	/*
		id_admin, _ := c.Locals("id_admin").(int)
		role, _ := c.Locals("role").(string)
		names, _ := c.Locals("name").(string)

		tools.ValidationJwtUsers(c, role, id_admin, names)

	*/

	id_users := c.Query("idUsers")
	id_pelatihan := c.Query("idPelatihan")

	var usersPelatihan []entity.UsersPelatihan
	baseQuey := database.DB

	if id_users != "" {
		baseQuey = baseQuey.Where("id_users = ?", id_users)
	}
	if id_pelatihan != "" {
		baseQuey = baseQuey.Where("id_pelatihan = ?", id_pelatihan)
	}

	baseQuey.Find(&usersPelatihan)

	return c.JSON(fiber.Map{
		"pesan": "Sukses Mengambil data",
		"data":  usersPelatihan,
	})
}

// Test dengan Gorm relasi cuy
func GetPelatihanByUser(c *fiber.Ctx) error {
	userID := c.Query("userID")
	var user entity.Users

	if err := database.DB.Preload("Pelatihan").First(&user, userID).Error; err != nil {
		return c.Status(404).SendString(err.Error())
	}

	return c.JSON(user)
}

func GetUsersByPelatihan(c *fiber.Ctx) error {
	/*
		//JWT nya siapa ntar ?
		id_admin, _ := c.Locals("id_admin").(int)
		role, _ := c.Locals("role").(string)
		names, _ := c.Locals("name").(string)

		tools.ValidationJwt(c, role, id_admin, names)

	*/

	idPelatihan := c.Query("idPelatihan")

	var pelatihan entity.Pelatihan

	if err := database.DB.Preload("UserPelatihan").Find(&pelatihan, idPelatihan).Error; err != nil {
		return c.Status(404).SendString(err.Error())
	}

	return c.JSON(pelatihan)
}
