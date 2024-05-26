package routes

import (
	"template/app/controllers"
	"template/pkg/middleware"
	"template/public/static"

	"github.com/gofiber/fiber/v2"
)

func SetupRoutesFiber(app *fiber.App) {

	app.Get("/", func(c *fiber.Ctx) error {
		return c.Render("welcome.html", fiber.Map{})
	})

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
	users.Put("/updateUsers", middleware.JwtProtect(), controllers.UpdateUser)

	users.Get("/test", controllers.TestPreloadPencapaian)

	//User Post Add Pelatihan
	users.Post("/addPelatihan", middleware.JwtProtect(), controllers.CreateUserPelatihan)

	//lemdik Area
	lemdik.Post("/login", controllers.LoginLemdik)
	lemdik.Get("/getLemdik", controllers.GetLemdik)
	lemdik.Put("/update", controllers.UpdateLemdik)
	//Pelatihan
	lemdik.Post("/createPelatihan", middleware.JwtProtect(), controllers.CreatePelatihan)
	lemdik.Get("/getPelatihan", controllers.GetPelatihan)

	//Sarpras
	lemdik.Post("/createSarpras", middleware.JwtProtect(), controllers.CreateSarpras)
	lemdik.Get("/getSarpras", middleware.JwtProtect(), controllers.GetSarpras)

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
	app.Get("/public/static/sertifikasi/:string", static.StaticSertifikasi)
}
