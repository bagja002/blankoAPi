package routes

import (
	"template/app/controllers"

	"github.com/gofiber/fiber/v2"
)


func SetupRoutesFiber(app *fiber.App){

	app.Get("/", func(c *fiber.Ctx) error {
        return c.Render("welcome.html", fiber.Map{})
    })

	app.Get("/test",controllers.CreateUser)

	SuperAdmin:= app.Group("/superadmin")

	SuperAdmin.Post("/login", controllers.SuperAdminLogin)
	
}