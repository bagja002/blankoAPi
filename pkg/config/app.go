package config

import (
	//"log"

	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/spf13/viper"
)

// Kumpulan Konfigurasi Go Fiber
func NewFiber(config *viper.Viper) *fiber.App {
	var app = fiber.New(fiber.Config{
		AppName:      config.GetString("app.name"),
		ErrorHandler: NewErrorHandler(),
		Prefork:      config.GetBool("web.prefork"),
		BodyLimit:    20 * 1024 * 1024,
		Concurrency:  1024 * 1024,
	})

	app.Use(
		cors.New(),
		logger.New(),
		limiter.New(limiter.Config{
			Max:        1000,
			Expiration: 30 * time.Second,
			KeyGenerator: func(c *fiber.Ctx) string {
				return c.IP()
			},
			LimitReached: func(c *fiber.Ctx) error {
				return c.SendStatus(fiber.StatusTooManyRequests)
			},
			SkipFailedRequests:     false,
			SkipSuccessfulRequests: false,
			LimiterMiddleware:      limiter.FixedWindow{},
		}),
	)

	return app
}

func NewErrorHandler() fiber.ErrorHandler {
	return func(ctx *fiber.Ctx, err error) error {
		code := fiber.StatusInternalServerError
		if e, ok := err.(*fiber.Error); ok {
			code = e.Code
		}

		return ctx.Status(code).JSON(fiber.Map{
			"errors": err.Error(),
		})
	}
}
