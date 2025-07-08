package notification

import (
	"github.com/gofiber/fiber/v2"
	"go-coursework/internal/handlers/notification"
	"go-coursework/internal/models"
	"go-coursework/pkg/jwt"
	"go-coursework/pkg/limiter"
)

func Setup(router fiber.Router, rctx *models.RouterContext) {
	fcm, err := notification.NewFCMHandler("../../../go-coursework-fb7be95d2ee0.json", rctx)
	if err != nil {
		panic(err)
	}

	notificationGroup := router.Get("/notifications").Use(jwt.Middleware("Lecturer", "Student"), limiter.RateLimiter())
	notificationGroup.Post("/:id", fcm.SaveFCMToken)
}
