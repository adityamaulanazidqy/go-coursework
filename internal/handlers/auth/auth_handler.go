package auth

import (
	"github.com/gofiber/fiber/v2"
	"go-coursework/internal/dto/auth"
	"go-coursework/internal/helpers"
	"go-coursework/internal/logger"
	"go-coursework/internal/repositories"
	"gorm.io/gorm"
)

type AuthenticationHandler struct {
	db        *gorm.DB
	logLogrus *logger.ErrorLogger
	authRepo  *repositories.AuthenticationRepo
}

func NewAuthenticationHandler(db *gorm.DB, logLogrus *logger.ErrorLogger) *AuthenticationHandler {
	return &AuthenticationHandler{db: db, logLogrus: logLogrus, authRepo: repositories.NewAuthenticationRepo(db)}
}

func (h *AuthenticationHandler) SignIn(ctx *fiber.Ctx) error {
	const op = "handler.AuthHandler.SignIn"

	var (
		msgBodyParse            = "Failed to parse body"
		msgBodyParseDetails     = []string{"One of the requests is not eligible"}
		msgValidateError        = "There was an error while validating Sign In"
		msgValidateErrorDetails = []string{"Use a valid email format", "Make sure you have filled in everything"}
	)

	var req auth.SignInRequest
	if err := ctx.BodyParser(&req); err != nil {
		h.logLogrus.LogUserError(req.Email, err, msgBodyParse)
		return h.logLogrus.LogRequestError(ctx, fiber.StatusBadRequest, op, err, msgBodyParse, msgBodyParseDetails)
	}

	if err := helpers.ValidateLoginRequest(req, "@gmail.com"); err != nil {
		h.logLogrus.LogUserError(req.Email, err, msgValidateError)
		return h.logLogrus.LogRequestError(ctx, fiber.StatusBadRequest, op, err, msgValidateError, msgValidateErrorDetails)
	}

	resp, code, opRepo, err, msg, details := h.authRepo.SignIn(&req)
	if err != nil {
		h.logLogrus.LogUserError(req.Email, err, msg)
		return h.logLogrus.LogRequestError(ctx, code, opRepo, err, msg, details)
	}

	return ctx.Status(code).JSON(fiber.Map{
		"message": msg,
		"data":    resp,
	})
}

func (h *AuthenticationHandler) SignUp(ctx *fiber.Ctx) error {
	const op = "handler.AuthHandler.SignUp"

	var (
		msgBodyParse            = "Failed to parse body"
		msgBodyParseDetails     = []string{"One of the requests is not eligible"}
		msgValidateError        = "There was an error while validating Sign In"
		msgValidateErrorDetails = []string{"Use a valid email format", "Make sure you have filled in everything"}
	)

	var req auth.SignUpRequest
	if err := ctx.BodyParser(&req); err != nil {
		h.logLogrus.LogUserError(req.Email, err, msgBodyParse)
		return h.logLogrus.LogRequestError(ctx, fiber.StatusBadRequest, op, err, msgBodyParse, msgBodyParseDetails)
	}

	if err := helpers.ValidateRegisterRequest(req, "@gmail.com"); err != nil {
		h.logLogrus.LogUserError(req.Email, err, msgValidateError)
		return h.logLogrus.LogRequestError(ctx, fiber.StatusBadRequest, op, err, msgValidateError, msgValidateErrorDetails)
	}

	resp, code, opRepo, err, msg, details := h.authRepo.SignUp(&req)
	if err != nil {
		h.logLogrus.LogUserError(req.Email, err, msg)
		return h.logLogrus.LogRequestError(ctx, code, opRepo, err, msg, details)
	}

	return ctx.Status(code).JSON(fiber.Map{
		"message": msg,
		"data":    resp,
	})
}
