package routes

import (
	"template/app/controllers"
	"template/pkg/middleware"
	"template/public/static"

	"github.com/gofiber/fiber/v2"
)


func SetupRoutesFiber(app *fiber.App){

	app.Get("/", func(c *fiber.Ctx) error {
        return c.Render("welcome.html", fiber.Map{})
    })

	lemdik:= app.Group("/lemdik")
	//adminPusat := app.Group("/adminpusat")
	app.Get("/test",controllers.CreateUser)

	//lemdik Area 
	//Pelatihan 
	lemdik.Post("/createPelatihan", controllers.CreatePelatihan)
	lemdik.Get("/getPelatihan", controllers.GetPelatihan)
	lemdik.Post("/login", controllers.LoginLemdik)



	//super admin
	//Create User area
	SuperAdmin:= app.Group("/superadmin")

	SuperAdmin.Post("/regiterLemdik",middleware.JwtProtect(), controllers.RegisterLemdik)


	SuperAdmin.Post("/login", controllers.SuperAdminLogin)


	//static file 

	app.Get("/public/static/pelatihan/:string", static.StaticPelatihan)
	app.Get("/public/static/prasarana/:string", static.StaticPrasarana)
	app.Get("/public/static/profile/:string", static.StaticProfile)
	app.Get("/public/static/sertifikasi/:string", static.StaticSertifikasi)
}