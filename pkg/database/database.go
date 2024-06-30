package database

import (
	"fmt"
	"log"
	"template/app/entity"
	"template/pkg/config"
	"template/pkg/tools"
	"time"

	//"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func GenerateLembagaDiklat() {

}

func Connect() {

	viper := config.NewViper()
	username := viper.GetString("database.username")
	password := viper.GetString("database.password")
	host := viper.GetString("database.host")
	port := viper.GetInt("database.port")
	database := viper.GetString("database.name")
	idleConnection := viper.GetInt("database.pool.idle")
	maxConnection := viper.GetInt("database.pool.max")
	maxLifeTimeConnection := viper.GetInt("database.pool.lifetime")

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local", username, password, host, port, database)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}

	connection, err := db.DB()
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}

	connection.SetMaxIdleConns(idleConnection)
	connection.SetMaxOpenConns(maxConnection)
	connection.SetConnMaxLifetime(time.Second * time.Duration(maxLifeTimeConnection))

	err = db.AutoMigrate(
		&entity.Admin{}, &entity.Blanko{}, &entity.BlankoKeluar{},
		&entity.BlankoRusak{},
	)
	if err != nil {
		log.Fatalf("failed to auto-migrate: %v", err)
	}

	// Cek apakah akun Super Admin sudah terbuat
	var existingSuperAdmin entity.SuperAdmin
	if err := db.Where("username = ? ", "super").Find(&existingSuperAdmin).Error; err != nil {
		// Penanganan kesalahan jika terjadi
		fmt.Println("Gagal membuat atau menemukan akun Super Admin:", err)
	} else {
		if existingSuperAdmin.IdSuperAdmin == 0 {
			// Akun Super Admin baru berhasil dibuat
			super := entity.SuperAdmin{
				Username: "super",
				Nama:     "superadmin",
				Email:    "superadmin@puslat.com",
				Password: tools.GeneratePassword("superadmin"),
			}
			db.Create(&super)
			fmt.Println("Akun Super Admin baru berhasil dibuat")
		} else {
			// Akun Super Admin sudah ada
			fmt.Println("Akun Super Admin sudah ada")
		}
	}
	DB = db

}
