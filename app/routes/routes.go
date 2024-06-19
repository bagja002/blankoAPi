package routes

import (
	"template/app/controllers"
	"template/pkg/middleware"
	"template/pkg/tools"
	"template/public/static"

	"github.com/gofiber/fiber/v2"
)

func SetupRoutesFiber(app *fiber.App) {

	app.Get("/", func(c *fiber.Ctx) error {
		return c.Render("welcome.html", fiber.Map{})
	})

	app.Get("/getDataKusuka", tools.GetDataKusuka)

	lemdik := app.Group("/lemdik")
	users := app.Group("/users")
	adminPusat := app.Group("/adminPusat")

	//adminPusat := app.Group("/adminpusat")
	app.Get("/test", controllers.CreateUser)

	app.Get("/getUserPelatihan", controllers.GetPelatihanByUser)
	//
	app.Get("/getPelatihanUser", controllers.GetUsersByPelatihan)

	//AdminPusat Area
	adminPusat.Post("/login", controllers.LoginAdminPusat)
	adminPusat.Get("/getAdminPusat", middleware.JwtProtect(), controllers.GetAdminPusat)

	//Users Area
	users.Post("/registerUser", controllers.CreateUser)
	users.Post("/login", controllers.LoginUsers)
	users.Get("/getUsersById", middleware.JwtProtect(), controllers.GetUserByID)
	lemdik.Get("/getAllUsers", middleware.JwtProtect(), controllers.GetAllUsers)
	users.Put("/updateUsers", middleware.JwtProtect(), controllers.UpdateUser)

	users.Get("/test", controllers.TestPreloadPencapaian)

	//User Post Add Pelatihan
	users.Post("/addPelatihan", middleware.JwtProtect(), controllers.CreateUserPelatihan)
	lemdik.Put("/updatePelatihanUsers", middleware.JwtProtect(), controllers.UpdateUsersPelatihan)
	//lemdik Area
	lemdik.Post("/login", controllers.LoginLemdik)
	lemdik.Get("/getLemdik", middleware.JwtProtect(), controllers.GetLemdik)
	lemdik.Put("/update", middleware.JwtProtect(), controllers.UpdateLemdik)
	lemdik.Get("/getAllUsers", middleware.JwtProtect(), controllers.GetAllUsers)
	//Pelatihan
	lemdik.Post("/createPelatihan", middleware.JwtProtect(), controllers.CreatePelatihan)
	lemdik.Post("/createMateriPelatihan", middleware.JwtProtect(), controllers.CreateMateriPelatihan)
	lemdik.Put("/updatePelatihan", middleware.JwtProtect(), controllers.UpdatePelatihan)
	lemdik.Get("/getPelatihan", controllers.GetPelatihan)

	lemdik.Post("/PublishSertifikat", middleware.JwtProtect(), controllers.PublishSertifikat)
	lemdik.Post("/LastNomorSertifBalai", middleware.JwtProtect(), controllers.LastNomorSertifBalai)

	//lemdik.Put("/updateLastSertif", middleware.JwtProtect(), controllers.AddLastSertifLowBalai)

	//Sarpras
	lemdik.Post("/createSarpras", middleware.JwtProtect(), controllers.CreateSarpras)
	lemdik.Get("/getSarpras", middleware.JwtProtect(), controllers.GetSarpras)
	lemdik.Put("/updateSarpras", middleware.JwtProtect(), controllers.UpdateSarpras)
	lemdik.Delete("/deleteSarpras", middleware.JwtProtect(), controllers.DeleteSarpras)

	//Pelatihan Users Area

	//super admin
	//Create User area
	SuperAdmin := app.Group("/superadmin")

	SuperAdmin.Post("/registerAdminPusat", middleware.JwtProtect(), controllers.CreateAdminPusat)
	SuperAdmin.Post("/regiterLemdik", middleware.JwtProtect(), controllers.RegisterLemdik)
	SuperAdmin.Post("/login", controllers.SuperAdminLogin)

	//static file

	app.Get("/public/static/pelatihan/:string", static.StaticPelatihan)
	app.Get("/public/static/prasarana/:string", static.StaticPrasarana)
	app.Get("/public/static/profile/:string", static.StaticProfile)
	app.Get("/public/silabus/sertifikasi/:string", static.StaticSilabusSertifikasi)
	app.Get("/public/static/sarpras/:string", static.StaticSarpras)

	//Get Pelatihan FIle
	app.Get("/public/silabus/pelatihan/:string", static.StaticSilabusPelatihan)
	app.Get("/public/module/pelatihan/:string", static.StaticModulePelatihan)
	app.Get("/public/static/BeritaAcara/:string", static.StaticBeritaAcara)
	app.Get("/public/static/suratPemberitahuan/:string", static.StaticSuratPemberitahuan)
	app.Get("/public/static/memo/:string", static.StaticMemo)
	//Get Users File File

	app.Get("/public/static/profile/fotoProfile/:string", static.StaticFotoUsers)
	app.Get("/public/static/profile/ijazah/:string", static.StaticIjazah)
	app.Get("/public/static/profile/kk/:string", static.StaticKK)
	app.Get("/public/static/profile/ktp/:string", static.StaticKtp)
	app.Get("/public/static/profile/suratSehat/:string", static.StaticSuratSehat)
	//Users Static area

	app.Get("/public/static/sertifikasi/:string", static.StaticSertifikasi)

	app.Get("/public/static/sertifikasi/:string", static.StaticSertifikasi)

	///Cek Sertifikat
	app.Post("/cekSertifikat", controllers.CekSertifikat)
}
