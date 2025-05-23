package models

import (
	"github.com/gofiber/fiber/v2"
	"github.com/redis/go-redis/v9"
	"go-coursework/internal/logger"
	"gorm.io/gorm"
)

type RouterContext struct {
	App         *fiber.App
	DB          *gorm.DB
	Logger      *logger.ErrorLogger
	RedisClient *redis.Client
}
