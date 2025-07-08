package otp

import (
	"github.com/gofiber/fiber/v2"
	"go-coursework/internal/handlers/otp"
	"go-coursework/internal/models"
	"go-coursework/pkg/jwt"
)

func Setup(router fiber.Router, rctx *models.RouterContext) {
	controller := otp.NewEmailVerification(rctx.DB, rctx.Logger, rctx.RedisClient, rctx)
	ct := otp.NewWaVerificationHandler(rctx)

	otpGroup := router.Group("/otp")
	{
		sendGroup := otpGroup.Group("/send")
		{
			sendGroup.Post("/email", controller.SendOtpEmail)
		}

		verifyGroup := otpGroup.Group("/verify")
		{
			verifyGroup.Post("/email", controller.VerifyOtpEmail)
			verifyGroup.Post("/telephone", jwt.Middleware("Lecturer", "Student"), ct.VerificationWaTokenService)
		}
	}
}
