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

func GetARPU(cache *redis.Client, warehouse *db.WarehouseClient) fiber.Handler {
	return func(c *fiber.Ctx) error {
		startDate, endDate := resolveDateRange(c.Query("start_date"), c.Query("end_date"))
		accountID := c.Query("account_id")

		cacheKey := "arpu:" + startDate + ":" + endDate + ":" + accountID
		if cached, ok := getCache(c.Context(), cache, cacheKey); ok {
			return c.Status(http.StatusOK).JSON(MetricResponse{
				Metric:     "arpu",
				Value:      cached,
				UpdatedAt:  time.Now().UTC().Format(time.RFC3339),
				Cached:     true,
				TimeWindow: startDate + " to " + endDate,
			})
		}

		value := warehouse.GetARPU(context.Background(), startDate, endDate, accountID)
		setCache(c.Context(), cache, cacheKey, value, 10*time.Minute)

		return c.Status(http.StatusOK).JSON(MetricResponse{
			Metric:     "arpu",
			Value:      value,
			UpdatedAt:  time.Now().UTC().Format(time.RFC3339),
			Cached:     false,
			TimeWindow: startDate + " to " + endDate,
		})
	}
}

func GetMRR(cache *redis.Client, warehouse *db.WarehouseClient) fiber.Handler {
	return func(c *fiber.Ctx) error {
		accountID := c.Query("account_id")

		cacheKey := "mrr:" + accountID
		if cached, ok := getCache(c.Context(), cache, cacheKey); ok {
			return c.Status(http.StatusOK).JSON(MetricResponse{
				Metric:     "mrr",
				Value:      cached,
				UpdatedAt:  time.Now().UTC().Format(time.RFC3339),
				Cached:     true,
				TimeWindow: "current",
			})
		}

		value := warehouse.GetMRR(context.Background(), accountID)
		setCache(c.Context(), cache, cacheKey, value, 15*time.Minute)

		return c.Status(http.StatusOK).JSON(MetricResponse{
			Metric:     "mrr",
			Value:      value,
			UpdatedAt:  time.Now().UTC().Format(time.RFC3339),
			Cached:     false,
			TimeWindow: "current",
		})
	}
}

func GetNRR(cache *redis.Client, warehouse *db.WarehouseClient) fiber.Handler {
	return func(c *fiber.Ctx) error {
		startDate, endDate := resolveDateRange(c.Query("start_date"), c.Query("end_date"))
		accountID := c.Query("account_id")

		cacheKey := "nrr:" + startDate + ":" + endDate + ":" + accountID
		if cached, ok := getCache(c.Context(), cache, cacheKey); ok {
			return c.Status(http.StatusOK).JSON(MetricResponse{
				Metric:     "nrr",
				Value:      cached,
				UpdatedAt:  time.Now().UTC().Format(time.RFC3339),
				Cached:     true,
				TimeWindow: startDate + " to " + endDate,
			})
		}

		value := warehouse.GetNRR(context.Background(), startDate, endDate, accountID)
		setCache(c.Context(), cache, cacheKey, value, 30*time.Minute)

		return c.Status(http.StatusOK).JSON(MetricResponse{
			Metric:     "nrr",
			Value:      value,
			UpdatedAt:  time.Now().UTC().Format(time.RFC3339),
			Cached:     false,
			TimeWindow: startDate + " to " + endDate,
		})
	}
}

func GetChurnRate(cache *redis.Client, warehouse *db.WarehouseClient) fiber.Handler {
	return func(c *fiber.Ctx) error {
		startDate, endDate := resolveDateRange(c.Query("start_date"), c.Query("end_date"))
		accountID := c.Query("account_id")

		cacheKey := "churn_rate:" + startDate + ":" + endDate + ":" + accountID
		if cached, ok := getCache(c.Context(), cache, cacheKey); ok {
			return c.Status(http.StatusOK).JSON(MetricResponse{
				Metric:     "churn_rate",
				Value:      cached,
				UpdatedAt:  time.Now().UTC().Format(time.RFC3339),
				Cached:     true,
				TimeWindow: startDate + " to " + endDate,
			})
		}

		value := warehouse.GetChurnRate(context.Background(), startDate, endDate, accountID)
		setCache(c.Context(), cache, cacheKey, value, 30*time.Minute)

		return c.Status(http.StatusOK).JSON(MetricResponse{
			Metric:     "churn_rate",
			Value:      value,
			UpdatedAt:  time.Now().UTC().Format(time.RFC3339),
			Cached:     false,
			TimeWindow: startDate + " to " + endDate,
		})
	}
}

func GetLTV(cache *redis.Client, warehouse *db.WarehouseClient) fiber.Handler {
	return func(c *fiber.Ctx) error {
		startDate, endDate := resolveDateRange(c.Query("start_date"), c.Query("end_date"))
		accountID := c.Query("account_id")

		cacheKey := "ltv:" + startDate + ":" + endDate + ":" + accountID
		if cached, ok := getCache(c.Context(), cache, cacheKey); ok {
			return c.Status(http.StatusOK).JSON(MetricResponse{
				Metric:     "ltv",
				Value:      cached,
				UpdatedAt:  time.Now().UTC().Format(time.RFC3339),
				Cached:     true,
				TimeWindow: startDate + " to " + endDate,
			})
		}

		value := warehouse.GetLTV(context.Background(), startDate, endDate, accountID)
		setCache(c.Context(), cache, cacheKey, value, 30*time.Minute)

		return c.Status(http.StatusOK).JSON(MetricResponse{
			Metric:     "ltv",
			Value:      value,
			UpdatedAt:  time.Now().UTC().Format(time.RFC3339),
			Cached:     false,
			TimeWindow: startDate + " to " + endDate,
		})
	}
}

func GetCAC(cache *redis.Client, warehouse *db.WarehouseClient) fiber.Handler {
	return func(c *fiber.Ctx) error {
		startDate, endDate := resolveDateRange(c.Query("start_date"), c.Query("end_date"))
		accountID := c.Query("account_id")

		cacheKey := "cac:" + startDate + ":" + endDate + ":" + accountID
		if cached, ok := getCache(c.Context(), cache, cacheKey); ok {
			return c.Status(http.StatusOK).JSON(MetricResponse{
				Metric:     "cac",
				Value:      cached,
				UpdatedAt:  time.Now().UTC().Format(time.RFC3339),
				Cached:     true,
				TimeWindow: startDate + " to " + endDate,
			})
		}

		value := warehouse.GetCAC(context.Background(), startDate, endDate, accountID)
		setCache(c.Context(), cache, cacheKey, value, 30*time.Minute)

		return c.Status(http.StatusOK).JSON(MetricResponse{
			Metric:     "cac",
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
