package middleware

import (
	"os"
	"strings"

	"github.com/gofiber/fiber/v2"
)

func AuthMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		requiredKey := os.Getenv("API_KEY")
		if requiredKey == "" {
			return c.Next()
		}

		provided := c.Get("X-API-Key")
		if provided == "" {
			provided = parseBearer(c.Get("Authorization"))
		}

		if provided == "" || provided != requiredKey {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "unauthorized",
			})
		}

		return c.Next()
	}
}

func parseBearer(value string) string {
	if value == "" {
		return ""
	}
	parts := strings.SplitN(value, " ", 2)
	if len(parts) != 2 {
		return ""
	}
	if strings.ToLower(parts[0]) != "bearer" {
		return ""
	}
	return strings.TrimSpace(parts[1])
}
