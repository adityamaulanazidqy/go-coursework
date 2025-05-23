package auth

import (
	"github.com/gofiber/fiber/v2"
	"go-coursework/internal/handlers/auth"
	"go-coursework/internal/models"
)

func Setup(router fiber.Router, ctx *models.RouterContext) {
	controller := auth.NewAuthenticationHandler(ctx.DB, ctx.Logger)
	authGroup := router.Group("/auth")
	{
		authGroup.Post("/signup", controller.SignUp)
		authGroup.Post("/signin", controller.SignIn)
	}
}
