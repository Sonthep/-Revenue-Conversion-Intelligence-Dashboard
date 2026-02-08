package main

import (
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"

	"revenue-dashboard-api/cache"
	"revenue-dashboard-api/db"
	"revenue-dashboard-api/handlers"
	"revenue-dashboard-api/middleware"
)

func main() {
	app := fiber.New()

	app.Use(logger.New())
	app.Use(cors.New(cors.Config{
		AllowOrigins: os.Getenv("ALLOWED_ORIGINS"),
	}))

	redisClient := cache.NewRedisClient()
	warehouse := db.NewWarehouseClient()

	api := app.Group("/api")
	api.Use(middleware.AuthMiddleware())

	api.Get("/metrics/revenue", handlers.GetRevenue(redisClient, warehouse))
	api.Get("/metrics/conversion-rate", handlers.GetConversionRate(redisClient, warehouse))
	api.Get("/health", handlers.Health())

	log.Fatal(app.Listen(":8080"))
}
