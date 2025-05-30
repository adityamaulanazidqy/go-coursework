package constant

import (
	"errors"
	"github.com/gofiber/fiber/v2"
	"go-coursework/constants"
	"go-coursework/internal/models"
	"go-coursework/internal/repositories"
	pkgerr "go-coursework/pkg/errors"
	"go-coursework/pkg/jwt"
)

type ConstHandler struct {
	rctx      *models.RouterContext
	constRepo *repositories.ConstRepo
}

func NewConstHandler(rctx *models.RouterContext) *ConstHandler {
	return &ConstHandler{rctx: rctx, constRepo: repositories.NewConstRepo(rctx)}
}

func (h *ConstHandler) GetSemesters(ctx *fiber.Ctx) error {
	const op = "handler.ConstHandler.GetSemester"

	user, ok := ctx.Locals("user").(*jwt.Claims)
	if !ok || user == nil {
		err := errors.New(pkgerr.ErrMissingClaims.Message)
		h.rctx.Logger.LogUserError("-", err, pkgerr.ErrMissingClaims.Message)
		return h.rctx.Logger.LogRequestError(ctx, fiber.StatusUnauthorized, op, err, pkgerr.ErrMissingClaims.Message, pkgerr.ErrMissingClaims.Details)
	}

	resp, code, opRepo, err, msg, details := h.constRepo.GetSemesters()
	if err != nil {
		return h.rctx.Logger.LogRequestError(ctx, code, opRepo, err, msg, details)
	}

	return ctx.Status(code).JSON(fiber.Map{
		"message": msg,
		"data":    resp,
	})
}

func (h *ConstHandler) PostSemesters(ctx *fiber.Ctx) error {
	const op = "handler.ConstHandler.PostSemesters"

	user, ok := ctx.Locals("user").(*jwt.Claims)
	if !ok || user == nil {
		err := errors.New(pkgerr.ErrMissingClaims.Message)
		h.rctx.Logger.LogUserError("-", err, pkgerr.ErrMissingClaims.Message)
		return h.rctx.Logger.LogRequestError(ctx, fiber.StatusUnauthorized, op, err, pkgerr.ErrMissingClaims.Message, pkgerr.ErrMissingClaims.Details)
	}

	var req constants.Semesters
	if err := ctx.BodyParser(&req); err != nil {
		return h.rctx.Logger.LogRequestError(ctx, fiber.StatusBadRequest, op, err, pkgerr.ErrBodyParse.Message, pkgerr.ErrBodyParse.Details)
	}

	if req.Name == "" {
		err := errors.New(pkgerr.ErrEmptyContent.Message)
		return h.rctx.Logger.LogRequestError(ctx, fiber.StatusBadRequest, op, err, pkgerr.ErrEmptyContent.Message, pkgerr.ErrEmptyContent.Details)
	}

	resp, code, opRepo, err, msg, details := h.constRepo.PostSemesters(&req)
	if err != nil {
		return h.rctx.Logger.LogRequestError(ctx, code, opRepo, err, msg, details)
	}

	return ctx.Status(code).JSON(fiber.Map{
		"message": msg,
		"data":    resp,
	})
}

func (h *ConstHandler) DeleteSemester(ctx *fiber.Ctx) error {
	const op = "handler.ConstHandler.DeleteSemester"

	user, ok := ctx.Locals("user").(*jwt.Claims)
	if !ok || user == nil {
		err := errors.New(pkgerr.ErrMissingClaims.Message)
		return h.rctx.Logger.LogRequestError(ctx, fiber.StatusUnauthorized, op, err, pkgerr.ErrMissingClaims.Message, pkgerr.ErrMissingClaims.Details)
	}

	var semester constants.Semesters
	if err := ctx.BodyParser(&semester); err != nil {
		return h.rctx.Logger.LogRequestError(ctx, fiber.StatusInternalServerError, op, err, pkgerr.ErrBodyParse.Message, pkgerr.ErrBodyParse.Details)
	}

	resp, code, opRepo, err, msg, details := h.constRepo.DeleteSemester(semester.ID)
	if err != nil {
		return h.rctx.Logger.LogRequestError(ctx, code, opRepo, err, msg, details)
	}

	return ctx.Status(code).JSON(fiber.Map{
		"message": msg,
		"data":    resp,
	})
}

func (h *ConstHandler) GetStudyPrograms(ctx *fiber.Ctx) error {
	const op = "handler.ConstHandler.GetStudyPrograms"

	user, ok := ctx.Locals("user").(*jwt.Claims)
	if !ok || user == nil {
		err := errors.New(pkgerr.ErrMissingClaims.Message)
		h.rctx.Logger.LogUserError("-", err, pkgerr.ErrMissingClaims.Message)
		return h.rctx.Logger.LogRequestError(ctx, fiber.StatusUnauthorized, op, err, pkgerr.ErrMissingClaims.Message, pkgerr.ErrMissingClaims.Details)
	}

	resp, code, opRepo, err, msg, details := h.constRepo.GetStudyPrograms()
	if err != nil {
		return h.rctx.Logger.LogRequestError(ctx, code, opRepo, err, msg, details)
	}

	return ctx.Status(code).JSON(fiber.Map{
		"message": msg,
		"data":    resp,
	})
}

func (h *ConstHandler) PostStudyPrograms(ctx *fiber.Ctx) error {
	const op = "handler.ConstHandler.PostStudyPrograms"

	user, ok := ctx.Locals("user").(*jwt.Claims)
	if !ok || user == nil {
		err := errors.New(pkgerr.ErrMissingClaims.Message)
		return h.rctx.Logger.LogRequestError(ctx, fiber.StatusUnauthorized, op, err, pkgerr.ErrMissingClaims.Message, pkgerr.ErrMissingClaims.Details)
	}

	var req constants.StudyPrograms
	if err := ctx.BodyParser(&req); err != nil {
		return h.rctx.Logger.LogRequestError(ctx, fiber.StatusBadRequest, op, err, pkgerr.ErrBodyParse.Message, pkgerr.ErrBodyParse.Details)
	}

	if req.Name == "" {
		err := errors.New(pkgerr.ErrEmptyContent.Message)
		return h.rctx.Logger.LogRequestError(ctx, fiber.StatusBadRequest, op, err, pkgerr.ErrEmptyContent.Message, pkgerr.ErrEmptyContent.Details)
	}

	resp, code, opRepo, err, msg, details := h.constRepo.PostStudyPrograms(&req)
	if err != nil {
		return h.rctx.Logger.LogRequestError(ctx, code, opRepo, err, msg, details)
	}

	return ctx.Status(code).JSON(fiber.Map{
		"message": msg,
		"data":    resp,
	})
}

func (h *ConstHandler) DeleteProgram(ctx *fiber.Ctx) error {
	const op = "handler.ConstHandler.DeleteProgram"

	user, ok := ctx.Locals("user").(*jwt.Claims)
	if !ok || user == nil {
		err := errors.New(pkgerr.ErrMissingClaims.Message)
		return h.rctx.Logger.LogRequestError(ctx, fiber.StatusUnauthorized, op, err, pkgerr.ErrMissingClaims.Message, pkgerr.ErrMissingClaims.Details)
	}

	var req constants.StudyPrograms
	if err := ctx.BodyParser(&req); err != nil {
		return h.rctx.Logger.LogRequestError(ctx, fiber.StatusInternalServerError, op, err, pkgerr.ErrBodyParse.Message, pkgerr.ErrBodyParse.Details)
	}

	resp, code, opRepo, err, msg, details := h.constRepo.DeleteProgram(req.ID)
	if err != nil {
		return h.rctx.Logger.LogRequestError(ctx, code, opRepo, err, msg, details)
	}

	return ctx.Status(code).JSON(fiber.Map{
		"message": msg,
		"data":    resp,
	})
}
