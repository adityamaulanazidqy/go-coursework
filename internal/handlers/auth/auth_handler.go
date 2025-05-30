package auth

import (
	"context"
	"github.com/gofiber/fiber/v2"
	"github.com/redis/go-redis/v9"
	"go-coursework/internal/dto/auth"
	"go-coursework/internal/helpers"
	"go-coursework/internal/logger"
	"go-coursework/internal/repositories"
	"go-coursework/pkg/jwt"
	"gorm.io/gorm"
	"time"
)

type AuthenticationHandler struct {
	db          *gorm.DB
	logLogrus   *logger.ErrorLogger
	redisClient *redis.Client
	authRepo    *repositories.AuthenticationRepo
}

func NewAuthenticationHandler(db *gorm.DB, logLogrus *logger.ErrorLogger, redisClient *redis.Client) *AuthenticationHandler {
	return &AuthenticationHandler{db: db, logLogrus: logLogrus, redisClient: redisClient, authRepo: repositories.NewAuthenticationRepo(db)}
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
		msgValidateErrorDetails = []string{"Use a valid email format", "Make sure you have filled in everything", "make sure the password is more than 6 long"}
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
	req.SemesterID = 1

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

func (h *AuthenticationHandler) Logout(ctx *fiber.Ctx) error {
	const op = "handler.AuthHandler.Logout"

	var (
		msgTokenError                 = "Failed to parse token"
		msgTokenErrorDetails          = []string{"Maybe you haven't logged in yet"}
		msgTokenFailedVerify          = "Token verification failed"
		msgTokenFailedVerifyDetails   = []string{"There was an error with your token. Please re-login to get a new access token"}
		msgInternalServerError        = "Internal Server Error"
		msgInternalServerErrorDetails = []string{"User Internal Server Error", "Make sure your internet is on", "There may be a problem with our server"}
	)

	token, err := jwt.ExtractTokenFromHeader(ctx)
	if err != nil {
		h.logLogrus.LogUserError("-", err, msgTokenError)
		return h.logLogrus.LogRequestError(ctx, fiber.StatusUnauthorized, op, err, msgTokenError, msgTokenErrorDetails)
	}

	claims, err := jwt.VerifyToken(token)
	if err != nil {
		h.logLogrus.LogUserError("-", err, msgTokenFailedVerify)
		return h.logLogrus.LogRequestError(ctx, fiber.StatusUnauthorized, op, err, msgTokenFailedVerify, msgTokenFailedVerifyDetails)
	}

	expDuration := time.Until(claims.ExpiresAt.Time)
	if expDuration <= 0 {
		expDuration = time.Minute * 1
	}

	ctxBackground := context.Background()
	err = h.redisClient.Set(ctxBackground, "blacklist:"+token, "true", expDuration).Err()
	if err != nil {
		h.logLogrus.LogUserError("-", err, msgInternalServerError)
		return h.logLogrus.LogRequestError(ctx, fiber.StatusInternalServerError, op, err, msgInternalServerError, msgInternalServerErrorDetails)
	}

	h.logLogrus.Logger.Info(ctx.IP() + " => User successfully logged out.")

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Successfully logged out.",
	})
}
