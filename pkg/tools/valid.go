package tools

import "github.com/gofiber/fiber/v2"





func ValidationJwt(c *fiber.Ctx, role string, id_admin int, names string) *fiber.Map {
    if role != "1" {
        return &fiber.Map{
            "pesan": "Role Bukan Admin Pusat",
        }
    }
    if id_admin == 0 {
        return &fiber.Map{
            "pesan": "Admin tidak terdaftar",
        }
    }
    if names == "" {
        return &fiber.Map{
            "pesan": "Tidak ada Nama di dalam Jwt",
        }
    }
    // Jika semuanya valid, kembalikan nilai null (tidak ada pesan kesalahan)
    return nil
}

func ValidationJwtLemdik(c *fiber.Ctx, role string, id_admin int, names string) *fiber.Map {
    if role != "2" {
        return &fiber.Map{
            "pesan": "Role Bukan Admin Pusat",
        }
    }
    if id_admin == 0 {
        return &fiber.Map{
            "pesan": "Admin tidak terdaftar",
        }
    }
    if names == "" {
        return &fiber.Map{
            "pesan": "Tidak ada Nama di dalam Jwt",
        }
    }
    // Jika semuanya valid, kembalikan nilai null (tidak ada pesan kesalahan)
    return nil
}

func ValidationJwtMitra(c *fiber.Ctx, role string, id_admin int, names string) *fiber.Map {
    if role != "3" {
        return &fiber.Map{
            "pesan": "Role Bukan Admin Pusat",
        }
    }
    if id_admin == 0 {
        return &fiber.Map{
            "pesan": "Admin tidak terdaftar",
        }
    }
    if names == "" {
        return &fiber.Map{
            "pesan": "Tidak ada Nama di dalam Jwt",
        }
    }
    // Jika semuanya valid, kembalikan nilai null (tidak ada pesan kesalahan)
    return nil
}

func ValidationJwtBPPSDM(c *fiber.Ctx, role string, id_admin int, names string) *fiber.Map {
    if role != "4" {
        return &fiber.Map{
            "pesan": "Role Bukan Admin Pusat",
        }
    }
    if id_admin == 0 {
        return &fiber.Map{
            "pesan": "Admin tidak terdaftar",
        }
    }
    if names == "" {
        return &fiber.Map{
            "pesan": "Tidak ada Nama di dalam Jwt",
        }
    }
    // Jika semuanya valid, kembalikan nilai null (tidak ada pesan kesalahan)
    return nil
}

 
func ValidationJwtUsers(c *fiber.Ctx, role string, id_admin int, names string) *fiber.Map {
    if role != "5" {
        return &fiber.Map{
            "pesan": "Role Bukan Admin Pusat",
        }
    }
    if id_admin == 0 {
        return &fiber.Map{
            "pesan": "Admin tidak terdaftar",
        }
    }
    if names == "" {
        return &fiber.Map{
            "pesan": "Tidak ada Nama di dalam Jwt",
        }
    }
    // Jika semuanya valid, kembalikan nilai null (tidak ada pesan kesalahan)
    return nil
}



