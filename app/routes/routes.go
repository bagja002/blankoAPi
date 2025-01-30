package routes

import (
	"template/app/controllers"
	"template/pkg/middleware"
	"template/pkg/static"

	"github.com/gofiber/fiber/v2"
)

func SetupRoutesFiber(app *fiber.App) {

	//superAdmin := app.Group("/superAdmin")

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Backend Blanko")
	})

	adminPusat := app.Group("/adminPusat")

	adminPusat.Post("/login", controllers.Login)
	adminPusat.Post("/register", controllers.Register)

	//CRUD Blanko
	adminPusat.Post("/addBlanko", middleware.JwtProtect(), controllers.CreteDataBlanko)
	adminPusat.Get("/getBlanko", controllers.GetBlanko)
	adminPusat.Put("/updateBlanko", middleware.JwtProtect(), controllers.UpdateBlanko)
	adminPusat.Delete("/deteleBlanko", middleware.JwtProtect(), controllers.DeleteBlanko)

	//CRUD BlankoKeluar
	// Routes for BlankoKeluar
	adminPusat.Post("/addBlankoKeluar", middleware.JwtProtect(), controllers.CreateBlankoKeluar)
	adminPusat.Get("/getBlankoKeluar", controllers.GetBlankoKeluar)
	adminPusat.Put("/updateBlankoKeluar", middleware.JwtProtect(), controllers.UpdateBlankoKeluar)
	adminPusat.Delete("/deleteBlankoKeluar", middleware.JwtProtect(), controllers.DeleteBlankoKeluar)

	// Routes for BlankoRusak
	adminPusat.Post("/addBlankoRusak", middleware.JwtProtect(), controllers.CreateBlankoRusak)
	adminPusat.Get("/getBlankoRusak", controllers.GetBlankoRusak)
	adminPusat.Put("/updateBlankoRusak", middleware.JwtProtect(), controllers.UpdateBlankoRusak)
	adminPusat.Delete("/deleteBlankoRusak", middleware.JwtProtect(), controllers.DeleteBlankoRusak)

	//Serah Terima
	adminPusat.Post("/createSerahTerima", controllers.CreateSerahterimaSertifikat)
	adminPusat.Get("/getSerahTerima", middleware.JwtProtect(), controllers.GetSerahterimaSertifikat)
	adminPusat.Put("/updateSerahTerima", middleware.JwtProtect(), controllers.UpdateSerahterimaSertifikat)
	adminPusat.Delete("/deleteSerahTerima/:id", middleware.JwtProtect(), controllers.DeleteSerahterimaSertifikat)

	//Pengiriman
	adminPusat.Post("/createPengiriman", middleware.JwtProtect(), controllers.CreatePengirimanSertifikat)
	adminPusat.Get("/getPengiriman", middleware.JwtProtect(), controllers.GetPengirimanSertifikat)
	adminPusat.Put("/updatePengiriman/:id", middleware.JwtProtect(), controllers.UpdatePengirimanSertifikat)
	adminPusat.Delete("/deletePengiriman/:id", middleware.JwtProtect(), controllers.DeletePengirimanSertifikat)

	app.Get("/public/static/foto-blanko-rusak/:string", static.StaticBelankoRusak)

	app.Get("/dataAkapi", controllers.GetDataSertifikat)
	app.Get("/lokasiDiklat", controllers.GetDataBalaiSertifikatLokasi)
	app.Get("/getSertifikatBalai", controllers.GetDataBalaiSertifikat)
	app.Get("/getByNameSertifkat", controllers.GetDataByNameUserSertifika)

	app.Get("/getBalaiSertifikat", controllers.GetDataBalaiSertifikats)
}
