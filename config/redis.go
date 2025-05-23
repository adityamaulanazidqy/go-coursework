package config

import (
	"context"
	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
	"os"
)

var (
	redisClient *redis.Client
)

func InitRedis() (*redis.Client, error) {
	_ = godotenv.Load()

	Logger.SetFormatter(&logrus.JSONFormatter{
		TimestampFormat: "2006-01-02 15:04:05",
	})

	var err error
	redisClient, err = initRedis()
	if err != nil {
		Logger.WithError(err).Error("Failed to initialize Redis")
		return nil, err
	}

	return redisClient, nil
}

func initRedis() (*redis.Client, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_ADDR"),
		Password: os.Getenv("REDIS_PASSWORD"),
		DB:       0,
	})

	_, err := rdb.Ping(context.Background()).Result()
	return rdb, err
}
