package controllers

import (
	//"backend-elaut/app/entity"

	"template/app/entity"
	"template/app/models"
	"template/pkg/database"
	"template/pkg/tools"
	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
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
	email := data["email"]
	err := database.DB.Where("email = ?", email).Find(&existingEmail).Error
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"Message": err,
		})
	}

	//cek email
	if existingEmail.Email == email {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"Message": "This Email is Register",
		})
	}

	NewUser := entity.Users{
		Nama:          data["name"],
		Email:         email,
		Password:      tools.GeneratePassword(data["password"]),
		NoTelpon: tools.StringToInt(data["no_number"]),
		CreateAt:      tools.TimeNowJakarta(),
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
		"Message": "Sucsses To Register Accont, Contratulation You Have 100 Point Right Now",
	})
}

func Login(c *fiber.Ctx) error {

	var data map[string]string

	if err := c.BodyParser(&data); err != nil {
		c.JSON(fiber.Map{
			"pesan": err.Error(),
		})
	}

	var users entity.Users

	database.DB.Where("username = ? ", data["username"]).First(&users)
	if users.IdUsers == 0 {
		c.Status(fiber.StatusUnauthorized)
		return c.JSON(fiber.Map{
			"pesan": "Username tidak di temukan",
		})
	}

	if err := bcrypt.CompareHashAndPassword([]byte("Password database"), []byte(data["password"])); err != nil {
		c.Status(fiber.StatusUnauthorized)
		return c.JSON(fiber.Map{
			"pesan": "Incorrect password!",
		})
	} else {
		t:= tools.GenerateToken(models.User{})
		return c.JSON(fiber.Map{
			"t":t,
		})
	}
}

// Get a user by ID
func GetUserByID(c *fiber.Ctx) error {
	id_admin, _ := c.Locals("id_admin").(string)
	role, _ := c.Locals("role").(string)
	names, _ := c.Locals("name").(string)

	if role != "3" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"pesan": "Role Bukan User ",
		})
	}
	if id_admin == "" {
		return c.JSON(fiber.Map{
			"pesan": "Admin tidak terdaftar ",
		})
	}
	if names == "" {
		return c.JSON(fiber.Map{
			"pesan": "Tidak ada Nama di dalam Jwt",
		})
	}

	var user entity.Users
	if err := database.DB.First(&user, id_admin).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"Message": "User not found",
		})
	}

	return c.JSON(user)
}

// Get all users
func GetAllUsers(c *fiber.Ctx) error {
	var users []entity.Users
	if err := database.DB.Find(&users).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"Message": "Failed to fetch users",
		})
	}

	return c.JSON(users)
}

// Update user by ID
func UpdateUser(c *fiber.Ctx) error {
	id_admin, _ := c.Locals("id_admin").(string)
	role, _ := c.Locals("role").(string)
	names, _ := c.Locals("name").(string)

	if role != "3" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"pesan": "Role Bukan User ",
		})
	}
	if id_admin == "" {
		return c.JSON(fiber.Map{
			"pesan": "Admin tidak terdaftar ",
		})
	}
	if names == "" {
		return c.JSON(fiber.Map{
			"pesan": "Tidak ada Nama di dalam Jwt",
		})
	}

	var user entity.Users
	if err := database.DB.First(&user, id_admin).Error; err != nil {
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
	if name, ok := data["name"]; ok {
		user.Nama = name
	}
	if email, ok := data["email"]; ok {
		user.Email = email
	}

	// Update other fields similarly

	// Save the updated user
	if err := database.DB.Save(&user).Error; err != nil {
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

