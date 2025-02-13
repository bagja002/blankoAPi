package main

import (
	"template/app/routes"
	"template/pkg/config"
	"template/pkg/database"
	"template/pkg/tools"

	//"backend-elaut/pkg/config"
	"log"
	"os"
	//"gorm.io/gorm"
	//"github.com/gofiber/fiber/v2"
	//"github.com/gofiber/fiber/v2"
)

//Terbaru unuuu

func main() {
	tools.CreateFolder()
	viperConfig := config.NewViper()

	file, err := os.OpenFile("app.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal("Gagal membuka file log:", err)
	}
	defer file.Close()

	database.Connect()

	database.Connect2()
	// Set output log ke file yang telah dibuka
	log.SetOutput(file)

	app := config.NewFiber(viperConfig)

	routes.SetupRoutesFiber(app)

	log.Fatal(app.Listen(config.NewViper().GetString("web.port")))
}
