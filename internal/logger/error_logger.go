package logger

import (
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
	"go-coursework/config"
	"go-coursework/internal/dto"
)

type ErrorLogger struct {
	Logger *logrus.Logger
}

func NewErrorLogger() *ErrorLogger {
	return &ErrorLogger{
		Logger: config.Logger,
	}
}

func (l *ErrorLogger) LogRequestError(ctx *fiber.Ctx, status int, op string, err error, message string, details []string) error {
	fields := logrus.Fields{
		"operation": op,
		"status":    status,
		"path":      ctx.Path(),
		"method":    ctx.Method(),
	}

	if err != nil {
		fields["error"] = err.Error()
	}

	l.Logger.WithFields(fields).Error(message)

	return ctx.Status(status).JSON(dto.ErrorResponse{
		Message: message,
		Details: details,
	})
}

func (l *ErrorLogger) LogUserError(email string, err error, message string) {
	fields := logrus.Fields{
		"email":   email,
		"message": message,
	}

	if err != nil {
		fields["error"] = err.Error()
	}

	l.Logger.WithFields(fields).Error(message)
}
