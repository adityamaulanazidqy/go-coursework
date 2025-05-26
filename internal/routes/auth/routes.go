package auth

import (
	"github.com/gofiber/fiber/v2"
	"go-coursework/internal/handlers/auth"
	"go-coursework/internal/models"
	"go-coursework/pkg/jwt"
	"go-coursework/pkg/limiter"
)

func Setup(router fiber.Router, ctx *models.RouterContext) {
	controller := auth.NewAuthenticationHandler(ctx.DB, ctx.Logger, ctx.RedisClient)
	authGroup := router.Group("/auth")
	{
		authGroup.Post("/signup", limiter.RateLimiter(), controller.SignUp)
		authGroup.Post("/signin", limiter.RateLimiter(), controller.SignIn)
		authGroup.Post("/logout", jwt.Middleware("Admin", "Lecturer", "Student"), controller.Logout)
	}
}
