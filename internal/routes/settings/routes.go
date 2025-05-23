package settings

import (
	"github.com/gofiber/fiber/v2"
	"go-coursework/internal/handlers/settings"
	"go-coursework/internal/models"
	"go-coursework/pkg/jwt"
)

func Setup(router fiber.Router, ctx *models.RouterContext) {
	controller := settings.NewSettingHandler(ctx.DB, ctx.Logger)
	SettingGroup := router.Group("/setting")
	{
		SetGroup := SettingGroup.Group("/set")
		{
			SetGroup.Post("/profile", jwt.Middleware("Admin", "Lecturer", "Student"), controller.SetProfile)
			SetGroup.Post("/telephone", jwt.Middleware("Admin", "Lecturer", "Student"), controller.SetTelephone)
		}
	}
}
