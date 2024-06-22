package middleware

import (
	"github.com/gofiber/fiber/v2"

	"fmt"
	"github.com/golang-jwt/jwt/v4"
)

/*
func JWTProtected() func(*fiber.Ctx) error {
	// Create config for JWT authentication middleware.
	config := jwtMiddleware.Config{
		SigningKey:   jwtMiddleware.SigningKey{Key: []byte(os.Getenv("JWT_SECRET_KEY"))},
		ContextKey:   "jwt", // used in private routes
		ErrorHandler: jwtError,
	}

	return jwtMiddleware.New(config)
}

*/

func JwtProtect() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Get Authorization header value
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"Pesan": "PLEASE USE JWT TOKEN !!!!!",
			})
		}
		// Check if Authorization header value is in the correct format
		if len(authHeader) < 8 || authHeader[:7] != "Bearer " {
			return c.JSON(fiber.Map{
				"status": "Invalid",
			})
		}

		// Extract token string
		tokenString := authHeader[7:]

		// Parse token
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			// Verify signing method and secret key
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
			}
			return []byte("secret"), nil
		})

		if err != nil {
			// Token not valid, handle error
			return c.Status(fiber.StatusUnauthorized).
				JSON(fiber.Map{
					"Pesan": "Maaf Gagal Mengambil Token",
				})
		}

		// Get claims from token
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok || !token.Valid {
			// Claims not valid, handle error
			return c.SendStatus(fiber.StatusUnauthorized)
		}

		// Extract data from claims and add it to the context
		c.Locals("name", claims["name"].(string))
		c.Locals("id_admin", int(claims["id_admin"].(float64)))
		c.Locals("role", claims["role"].(string))
		c.Locals("type", claims["type"].(string))

		// Continue with next middleware/handler
		return c.Next()
	}
}

func JwtExamProtect() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Get Authorization header value
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"Pesan": "PLEASE USE JWT TOKEN !!!!!",
			})
		}
		// Check if Authorization header value is in the correct format
		if len(authHeader) < 8 || authHeader[:7] != "Bearer " {
			return c.JSON(fiber.Map{
				"status": "Invalid",
			})
		}

		// Extract token string
		tokenString := authHeader[7:]

		// Parse token
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			// Verify signing method and secret key
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
			}
			return []byte("secret"), nil
		})

		if err != nil {
			// Token not valid, handle error
			return c.Status(fiber.StatusUnauthorized).
				JSON(fiber.Map{
					"Pesan": "Maaf Gagal Mengambil Token",
				})
		}

		// Get claims from token
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok || !token.Valid {
			// Claims not valid, handle error
			return c.SendStatus(fiber.StatusUnauthorized)
		}

		// Extract data from claims and add it to the context
		c.Locals("types", claims["type"].(string))
		c.Locals("id_users", int(claims["id_users"].(float64)))
		c.Locals("id_users_pelatihan", int(claims["id_users_pelatihan"].(float64)))

		// Continue with next middleware/handler
		return c.Next()
	}

}

// FiberMiddleware provide Fiber's built-in middlewares.
// See: https://docs.gofiber.io/api/middleware
