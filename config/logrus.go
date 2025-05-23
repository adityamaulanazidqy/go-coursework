package config

import (
	"github.com/getsentry/sentry-go"
	sentrylogrus "github.com/getsentry/sentry-go/logrus"
	"github.com/joho/godotenv"
	"os"
	"time"

	"github.com/sirupsen/logrus"
)

var Logger *logrus.Logger

func InitLogger() {
	_ = godotenv.Load()

	if sentryDSN := os.Getenv("SENTRY_DSN"); sentryDSN != "" {
		err := sentry.Init(sentry.ClientOptions{
			Dsn:              sentryDSN,
			Environment:      os.Getenv("ENVIRONMENT"),
			TracesSampleRate: 1.0,
			AttachStacktrace: true,
		})
		if err != nil {
			logrus.Errorf("Failed to initialize Sentry: %v", err)
		}
	}

	Logger = logrus.New()
	Logger.SetOutput(os.Stdout)
	Logger.SetFormatter(&logrus.JSONFormatter{
		TimestampFormat: "2006-01-02 15:04:05",
	})

	setLogLevel(Logger)

	if sentry.CurrentHub().Client() != nil {
		Logger.AddHook(sentrylogrus.NewFromClient([]logrus.Level{
			logrus.ErrorLevel,
			logrus.FatalLevel,
			logrus.PanicLevel,
		}, sentry.CurrentHub().Client()))
	}

	defer func() {
		if r := recover(); r != nil {
			Logger.Panicf("Recovered from panic: %v", r)
			sentry.Flush(2 * time.Second)
		}
	}()
}

func setLogLevel(logger *logrus.Logger) {
	level, err := logrus.ParseLevel(os.Getenv("LOG_LEVEL"))
	if err != nil {
		level = logrus.InfoLevel
	}
	logger.SetLevel(level)
}
