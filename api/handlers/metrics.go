package handlers

import (
	"context"
	"net/http"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/redis/go-redis/v9"

	"revenue-dashboard-api/db"
)

type MetricResponse struct {
	Metric     string      `json:"metric"`
	Value      interface{} `json:"value"`
	UpdatedAt  string      `json:"updated_at"`
	Cached     bool        `json:"cached"`
	TimeWindow string      `json:"time_window"`
}

func GetRevenue(cache *redis.Client, warehouse *db.WarehouseClient) fiber.Handler {
	return func(c *fiber.Ctx) error {
		startDate, endDate := resolveDateRange(c.Query("start_date"), c.Query("end_date"))
		accountID := c.Query("account_id")

		cacheKey := "revenue:" + startDate + ":" + endDate + ":" + accountID
		if cached, ok := getCache(c.Context(), cache, cacheKey); ok {
			return c.Status(http.StatusOK).JSON(MetricResponse{
				Metric:     "revenue",
				Value:      cached,
				UpdatedAt:  time.Now().UTC().Format(time.RFC3339),
				Cached:     true,
				TimeWindow: startDate + " to " + endDate,
			})
		}

		value := warehouse.GetRevenue(context.Background(), startDate, endDate, accountID)
		setCache(c.Context(), cache, cacheKey, value, 5*time.Minute)

		return c.Status(http.StatusOK).JSON(MetricResponse{
			Metric:     "revenue",
			Value:      value,
			UpdatedAt:  time.Now().UTC().Format(time.RFC3339),
			Cached:     false,
			TimeWindow: startDate + " to " + endDate,
		})
	}
}

func GetConversionRate(cache *redis.Client, warehouse *db.WarehouseClient) fiber.Handler {
	return func(c *fiber.Ctx) error {
		startDate, endDate := resolveDateRange(c.Query("start_date"), c.Query("end_date"))
		accountID := c.Query("account_id")

		cacheKey := "conversion_rate:" + startDate + ":" + endDate + ":" + accountID
		if cached, ok := getCache(c.Context(), cache, cacheKey); ok {
			return c.Status(http.StatusOK).JSON(MetricResponse{
				Metric:     "conversion_rate",
				Value:      cached,
				UpdatedAt:  time.Now().UTC().Format(time.RFC3339),
				Cached:     true,
				TimeWindow: startDate + " to " + endDate,
			})
		}

		value := warehouse.GetConversionRate(context.Background(), startDate, endDate, accountID)
		setCache(c.Context(), cache, cacheKey, value, 10*time.Minute)

		return c.Status(http.StatusOK).JSON(MetricResponse{
			Metric:     "conversion_rate",
			Value:      value,
			UpdatedAt:  time.Now().UTC().Format(time.RFC3339),
			Cached:     false,
			TimeWindow: startDate + " to " + endDate,
		})
	}
}

func Health() fiber.Handler {
	return func(c *fiber.Ctx) error {
		return c.Status(http.StatusOK).JSON(fiber.Map{"status": "ok"})
	}
}

func getCache(ctx context.Context, client *redis.Client, key string) (string, bool) {
	val, err := client.Get(ctx, key).Result()
	if err != nil {
		return "", false
	}
	return val, true
}

func setCache(ctx context.Context, client *redis.Client, key string, value string, ttl time.Duration) {
	_ = client.Set(ctx, key, value, ttl).Err()
}

func resolveDateRange(startDate, endDate string) (string, string) {
	if startDate == "" || endDate == "" {
		now := time.Now().UTC()
		endDate = now.Format("2006-01-02")
		startDate = now.AddDate(0, 0, -30).Format("2006-01-02")
		return startDate, endDate
	}
	return startDate, endDate
}
