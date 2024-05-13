package controllers

import (
	"template/app/entity"
	"template/pkg/database"
	"template/pkg/tools"

	"github.com/gofiber/fiber/v2"
)


func CreateSarpras(c *fiber.Ctx)error{

	id_admin, _ := c.Locals("id_admin").(int)
	role, _ := c.Locals("role").(string)
	names, _ := c.Locals("name").(string)

	tools.ValidationJwtLemdik(c ,role,id_admin,names)


	var data map[string]string

	if err := c.BodyParser(&data); err != nil {
		c.JSON(fiber.Map{
			"pesan": err.Error(),
		})
	}

	//masukan ke sarpras untuk membuat sarpras baru 

	newSarpras:= entity.Sarpras{
		IdLemdik: uint(id_admin),
		NamaSarpras: data["nama_sarpras"],
		Harga: tools.StringToInt(data["harga"]),
		Deskripsi: data["deskripsi"],
		Jenis: data["jenis"],
		CreateAt: tools.TimeNowJakarta(),
	}

	//save to database

	database.DB.Create(&newSarpras)


	return c.JSON(fiber.Map{
		"pesan":"Sukses Menambahkan Database Sarpras",
	})
}

func GetSarpras(c *fiber.Ctx) error {
    id_admin, _ := c.Locals("id_admin").(int)
    role, _ := c.Locals("role").(string)
    names, _ := c.Locals("name").(string)

    tools.ValidationJwtLemdik(c, role, id_admin, names)

    id := c.Query("id")

    baseQuery := database.DB.Model(&entity.Sarpras{})

    // Tambahkan kondisi filter berdasarkan ID jika ada
    if id != "" {
        baseQuery = baseQuery.Where("id_sarpras = ?", id)
    }

    jenisSarpras := c.Query("jenis_sarpras")
    // Tambahkan kondisi filter berdasarkan jenis sarpras jika ada
    if jenisSarpras != "" {
        baseQuery = baseQuery.Where("jenis= ?", jenisSarpras)
    }

    // Eksekusi query dan simpan hasilnya ke dalam slice Sarpras
    var Sarpras []entity.Sarpras
    baseQuery.Find(&Sarpras)

    return c.JSON(fiber.Map{
        "Pesan": "Sukses mendapatkan data",
        "data":  Sarpras,
    })
}
