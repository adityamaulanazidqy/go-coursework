package otp

import (
	"github.com/gofiber/fiber/v2"
	"go-coursework/internal/handlers/otp"
	"go-coursework/internal/models"
)

func Setup(router fiber.Router, ctx *models.RouterContext) {
	controller := otp.NewEmailVerification(ctx.DB, ctx.Logger, ctx.RedisClient)

	otpGroup := router.Group("/otp")
	{
		sendGroup := otpGroup.Group("/send")
		{
			sendGroup.Post("/email", controller.SendOtpEmail)
		}

		verifyGroup := otpGroup.Group("/verify")
		{
			verifyGroup.Post("/email", controller.VerifyOtpEmail)
		}
	}
}
