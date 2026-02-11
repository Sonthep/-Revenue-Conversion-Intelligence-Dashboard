package middleware

import (
	"os"
	"strings"

	"github.com/gofiber/fiber/v2"
)

func AuthMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		keyMap := parseKeyMap(os.Getenv("API_KEYS"))
		requiredKey := os.Getenv("API_KEY")
		if requiredKey == "" && len(keyMap) == 0 {
			return c.Next()
		}

		provided := c.Get("X-API-Key")
		if provided == "" {
			provided = parseBearer(c.Get("Authorization"))
		}
		if provided == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "unauthorized",
			})
		}

		if len(keyMap) > 0 {
			accountScope, ok := keyMap[provided]
			if !ok {
				return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
					"error": "unauthorized",
				})
			}
			if accountScope != "*" {
				requested := c.Query("account_id")
				if requested != "" && requested != accountScope {
					return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
						"error": "forbidden",
					})
				}
				c.Locals("account_id", accountScope)
			}
			return c.Next()
		}

		if provided != requiredKey {
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

func parseKeyMap(raw string) map[string]string {
	keyMap := map[string]string{}
	if raw == "" {
		return keyMap
	}
	entries := strings.Split(raw, ",")
	for _, entry := range entries {
		entry = strings.TrimSpace(entry)
		if entry == "" {
			continue
		}
		parts := strings.SplitN(entry, ":", 2)
		if len(parts) != 2 {
			continue
		}
		key := strings.TrimSpace(parts[0])
		account := strings.TrimSpace(parts[1])
		if key == "" || account == "" {
			continue
		}
		keyMap[key] = account
	}
	return keyMap
}
