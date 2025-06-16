package limiter

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/gofiber/storage/redis/v3"
	"github.com/joho/godotenv"
	"go-coursework/internal/dto"
	"os"
	"strconv"
	"time"
)

func RateLimiter() fiber.Handler {
	if err := godotenv.Load(".env"); err != nil {
		log.Fatal("Error loading .env file")
	}

	port, err := strconv.Atoi(os.Getenv("REDIS_PORT"))
	if err != nil {
		log.Fatal("Error parsing REDIS_PORT to int")
	}

	storage := redis.New(redis.Config{
		Addrs:    []string{os.Getenv("REDIS_ADDR")},
		Host:     os.Getenv("REDIS_HOST"),
		Password: os.Getenv("REDIS_PASSWORD"),
		Port:     port,
		Database: 0,
		Reset:    false,
	})

	return limiter.New(limiter.Config{
		Max:        20,
		Expiration: 1 * time.Minute,
		Storage:    storage,
		KeyGenerator: func(ctx *fiber.Ctx) string {
			return ctx.IP()
		},
		LimitReached: func(ctx *fiber.Ctx) error {
			return ctx.Status(fiber.StatusTooManyRequests).JSON(dto.Response{
				Message: "many requests. Please try again later",
			})
		}},
	)
}
