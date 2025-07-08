package settings

import (
	"github.com/gofiber/fiber/v2"
	"github.com/pkg/errors"
	"go-coursework/internal/dto/settings"
	"go-coursework/internal/handlers/otp"
	"go-coursework/internal/logger"
	"go-coursework/internal/models"
	"go-coursework/internal/repositories"
	pkgerr "go-coursework/pkg/errors"
	"go-coursework/pkg/jwt"
	"gorm.io/gorm"
	"mime/multipart"
)

type SettingHandler struct {
	db           *gorm.DB
	logLogrus    *logger.ErrorLogger
	SettingRepo  *repositories.SettingRepo
	waOtpHandler *otp.WaVerificationHandler
}

func NewSettingHandler(rctx *models.RouterContext) *SettingHandler {
	return &SettingHandler{
		db:           rctx.DB,
		logLogrus:    rctx.Logger,
		SettingRepo:  repositories.NewSettingRepo(rctx.DB),
		waOtpHandler: otp.NewWaVerificationHandler(rctx),
	}
}

func (h *SettingHandler) SetProfile(ctx *fiber.Ctx) error {
	const op = "handler.SettingHandler.SetProfile"

	var (
		msgInternalServerError        = "Internal Server Error"
		msgInternalServerErrorDetails = []string{"User Internal Server Error", "Make sure your internet is on", "There may be a problem with our server"}
		msgMissingClaims              = "Missing claims"
		msgMissingClaimsDetails       = []string{"There was an error with your token. Please re-login to get a new access token"}
	)

	claims, ok := ctx.Locals("user").(*jwt.Claims)
	if !ok || claims == nil {
		err := errors.New(msgMissingClaims)
		h.logLogrus.LogUserError("-", err, msgMissingClaims)
		return h.logLogrus.LogRequestError(ctx, fiber.StatusUnauthorized, op, err, msgMissingClaims, msgMissingClaimsDetails)
	}

	file, err := ctx.FormFile("profile")
	if err != nil {
		h.logLogrus.LogUserError(claims.Email, err, msgInternalServerError)
		return h.logLogrus.LogRequestError(ctx, fiber.StatusInternalServerError, op, err, msgInternalServerError, msgInternalServerErrorDetails)
	}

	multipart, err := file.Open()
	if err != nil {
		h.logLogrus.LogUserError(claims.Email, err, msgInternalServerError)
		return h.logLogrus.LogRequestError(ctx, fiber.StatusInternalServerError, op, err, msgInternalServerError, msgInternalServerErrorDetails)
	}
	defer multipart.Close()

	var req = settings.SetProfile{
		FileHeader:    file,
		MultipartFile: multipart,
	}

	resp, code, opRepo, err, msg, details := h.SettingRepo.SetProfile(&req, claims.UserID)
	if err != nil {
		h.logLogrus.LogUserError(claims.Email, err, msg)
		return h.logLogrus.LogRequestError(ctx, code, opRepo, err, msg, details)
	}

	return ctx.Status(code).JSON(fiber.Map{
		"message": msg,
		"data":    resp,
	})
}

func (h *SettingHandler) SetTelephone(ctx *fiber.Ctx) error {
	const op = "handler.SettingHandler.SetTelephone"

	var (
		msgBodyParse            = "Failed to parse body"
		msgBodyParseDetails     = []string{"One of the requests is not eligible"}
		msgMissingClaims        = "Missing claims"
		msgMissingClaimsDetails = []string{"There was an error with your token. Please re-login to get a new access token"}
	)

	claims, ok := ctx.Locals("user").(*jwt.Claims)
	if !ok || claims == nil {
		err := errors.New(msgMissingClaims)
		h.logLogrus.LogUserError("-", err, msgMissingClaims)
		return h.logLogrus.LogRequestError(ctx, fiber.StatusUnauthorized, op, err, msgMissingClaims, msgMissingClaimsDetails)
	}

	var req settings.SetTelephone

	if err := ctx.BodyParser(&req); err != nil {
		h.logLogrus.LogUserError(claims.Email, err, msgBodyParse)
		return h.logLogrus.LogRequestError(ctx, fiber.StatusBadRequest, op, err, msgBodyParse, msgBodyParseDetails)
	}

	resp, code, opRepo, err, msg, details := h.SettingRepo.SetTelephone(&req, claims.UserID)
	if err != nil {
		h.logLogrus.LogUserError(claims.Email, err, msg)
		return h.logLogrus.LogRequestError(ctx, code, opRepo, err, msg, details)
	}

	m, e := h.waOtpHandler.SendMessageWaOTP(claims.UserID, req.Telephone)
	if e != nil {
		h.logLogrus.LogUserError(claims.Email, e, m)
		return h.logLogrus.LogRequestError(ctx, fiber.StatusInternalServerError, op, e, pkgerr.ErrSendOTPWa.Message, pkgerr.ErrSendOTPWa.Details)
	}

	return ctx.Status(code).JSON(fiber.Map{
		"message": msg,
		"data":    resp,
	})
}

func (h *SettingHandler) UpdateUserInfo(ctx *fiber.Ctx) error {
	const op = "handler.SettingHandler.UpdateUsername"

	var (
		msgInternalServerError        = "Internal Server Error"
		msgInternalServerErrorDetails = []string{"User Internal Server Error", "Make sure your internet is on", "There may be a problem with our server"}
		msgMissingClaims              = "Missing claims"
		msgMissingClaimsDetails       = []string{"There was an error with your token. Please re-login to get a new access token"}
	)

	claims, ok := ctx.Locals("user").(*jwt.Claims)
	if !ok || claims == nil {
		err := errors.New(msgMissingClaims)
		h.logLogrus.LogUserError("-", err, msgMissingClaims)
		return h.logLogrus.LogRequestError(ctx, fiber.StatusUnauthorized, op, err, msgMissingClaims, msgMissingClaimsDetails)
	}

	var (
		fileHeader    *multipart.FileHeader
		multipartFile multipart.File
	)

	file, err := ctx.FormFile("profile")
	if err == nil && file != nil {
		multipartFileOpened, err := file.Open()
		if err != nil {
			h.logLogrus.LogUserError(claims.Email, err, msgInternalServerError)
			return h.logLogrus.LogRequestError(ctx, fiber.StatusInternalServerError, op, err, msgInternalServerError, msgInternalServerErrorDetails)
		}
		fileHeader = file
		multipartFile = multipartFileOpened
	}

	var req = settings.UpdateUserInfo{
		Username:  ctx.FormValue("username"),
		Email:     ctx.FormValue("email"),
		Telephone: ctx.FormValue("telephone"),
		Profile: settings.SetProfile{
			FileHeader:    fileHeader,
			MultipartFile: multipartFile,
		},
	}

	resp, code, opRepo, err, msg, details := h.SettingRepo.UpdateUserInfo(&req, claims.UserID)
	if err != nil {
		h.logLogrus.LogUserError(claims.Email, err, msg)
		return h.logLogrus.LogRequestError(ctx, code, opRepo, err, msg, details)
	}

	return ctx.Status(code).JSON(fiber.Map{
		"message": msg,
		"data":    resp,
	})
}
