package handler

import (
	"strings"

	"github.com/dpalhz/microservice-exp-with-go/internal/services/auth/token"
	"github.com/gofiber/fiber/v2"
)

func AuthMiddleware(jwtGen *token.JWTGenerator) fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "missing or invalid authorization header",
			})
		}

		tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
		claims, err := jwtGen.ValidateToken(tokenStr)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "invalid or expired token",
			})
		}

		// Simpan claims ke context untuk handler bisa akses
		c.Locals("user_id", claims["sub"])
		return c.Next()
	}
}
