package config

import (
	"errors"
	"github.com/golang-migrate/migrate/v4"
	postgres2 "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"os"
)

var db *gorm.DB

func Init() error {
	err := godotenv.Load()
	if err != nil {
		Logger.WithError(err).Error("Error loading .env file")
		return err
	}
	return nil
}

func ConnectDB() (*gorm.DB, error) {
	Logger.SetFormatter(&logrus.JSONFormatter{
		TimestampFormat: "2006-01-02 15:04:05",
	})

	err := Init()
	if err != nil {
		Logger.WithError(err).Error("Failed to Initialize DB")
		return nil, err
	}

	dsn := os.Getenv("DATABASE_URL")
	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		Logger.WithError(err).Error("Failed to connect to database")
		return nil, err
	}

	return db, nil
}

func RunMigrations() error {
	sqlDB, err := db.DB()
	if err != nil {
		Logger.WithError(err).Error("Failed to run migrations 1")
		return err
	}

	driver, err := postgres2.WithInstance(sqlDB, &postgres2.Config{})
	if err != nil {
		Logger.WithError(err).Error("Failed to run migrations 2")
		return err
	}

	m, err := migrate.NewWithDatabaseInstance("file://migrations", "postgres", driver)
	if err != nil {
		Logger.WithError(err).Error("Failed to run migrations 3")
		return err
	}

	if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		Logger.WithError(err).Error("Failed to run migrations 4")
		return err
	}

	log.Println("migration succeeded")
	return nil
}
