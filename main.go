package main

import (
	"github.com/gofiber/fiber/v2"
	"go-coursework/config"
	"go-coursework/internal/logger"
	"go-coursework/internal/models"
	"go-coursework/internal/routes"
	"go-coursework/pkg/jwt"
)

func main() {
	config.InitLogger()
	errorLogger := logger.NewErrorLogger()

	db, err := config.ConnectDB()
	if err != nil {
		return
	}

	err = config.RunMigrations()
	if err != nil {
		return
	}

	redisClient, err := config.InitRedis()
	if err != nil {
		return
	}

	app := fiber.New()

	routerCtx := models.RouterContext{
		DB:          db,
		App:         app,
		Logger:      errorLogger,
		RedisClient: redisClient,
	}

	jwt.SetRedisClientMiddleware(redisClient)

	routes.SetupRoutes(&routerCtx)

	if err := app.Listen(":8080"); err != nil {
		panic(err)
		return
	}
}
